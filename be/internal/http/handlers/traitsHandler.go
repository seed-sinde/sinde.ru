package handlers
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"sinde.ru/db/services"
	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/models"
	"sinde.ru/internal/store"
	"sinde.ru/utils"
)
// MemoryTraitHandler возвращает особенность по UUID.
//
// Маршрут:
//
//	GET /api/v1/traits/:uuid
//
// Возвращает:
//
//	JSON-массив вида [t_key_syn, t_value].
func MemoryTraitHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer utils.Benchmark("MemoryTraitHandler")()
		uuidStr := c.Params("uuid")
		if uuidStr == "" {
			return responses.Error(c, fiber.StatusBadRequest, "UUID обязателен")
		}
		id, err := uuid.Parse(uuidStr)
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "UUID некорректен")
		}
		trait, ok := store.GetTrait(id)
		if !ok {
			return responses.Error(c, fiber.StatusNotFound, "Особенность не найдена")
		}
		keys, ok := store.GetKeys(trait.TKey)
		if !ok || len(keys) == 0 {
			return responses.Error(c, fiber.StatusInternalServerError, "Ключ не найден")
		}
		key := keys[0]
		resp := []string{key.Syn, trait.TValue}
		c.Type("json")
		return c.JSON(resp)
	}
}
// MemoryAddTraitHandler добавляет особенность или возвращает существующую.
//
// Маршрут:
//
//	POST /api/v1/traits
//
// Тело запроса:
//
//	{"t_key":"<syn>","t_value":"<value>"}
//
// Возвращает:
//
//	Строку с t_uuid (канонический UUID из in-memory store).
func MemoryAddTraitHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer utils.Benchmark("MemoryAddTraitHandler")()
		var input struct {
			TKey   string `json:"t_key"`
			TValue string `json:"t_value"`
		}
		if err := c.Bind().JSON(&input); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "Некорректное тело запроса", err.Error())
		}
		input.TKey = strings.TrimSpace(input.TKey)
		input.TValue = strings.TrimSpace(input.TValue)
		if input.TKey == "" || input.TValue == "" {
			return responses.Error(c, fiber.StatusBadRequest, "Отсутствуют обязательные поля")
		}
		key, err := services.PdbGetOrCreateCanonicalKeyBySyn(c, input.TKey)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "Ошибка сохранения ключа", err.Error())
		}
		if key == nil {
			return responses.Error(c, fiber.StatusBadRequest, "Некорректный ключ")
		}
		store.SyncKey(key)
		trait, err := services.PdbGetOrCreateTraitByKeyValue(c, key.ID, input.TValue)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "Ошибка сохранения особенности", err.Error())
		}
		if trait == nil {
			return responses.Error(c, fiber.StatusInternalServerError, "Ошибка сохранения особенности")
		}
		resp := store.SyncTrait(trait)
		return c.SendString(resp.TUUID.String())
	}
}
// MemoryResolveUUIDHandler разрешает UUID в особенность или набор.
//
// Маршрут:
//
//	GET /api/v1/traits/resolve/:uuid
//
// Возвращает:
//
//	Для trait — компактный JSON-массив.
//	Для set — NDJSON-стрим всех особенностей.
func MemoryResolveUUIDHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer utils.Benchmark("MemoryResolveUUIDHandler")()
		uuidStr := c.Params("uuid")
		if uuidStr == "" {
			return responses.Error(c, 400, "UUID обязателен")
		}
		id, err := uuid.Parse(uuidStr)
		if err != nil {
			return responses.Error(c, 400, "UUID некорректен")
		}
		// ---------- 1) ОСОБЕННОСТЬ: компактный JSON + ручной ETag ----------
		if trait, ok := store.GetTrait(id); ok {
			keys, ok := store.GetKeys(trait.TKey)
			if !ok || len(keys) == 0 {
				return responses.Error(c, 500, "Ключ для особенности не найден")
			}
			key := keys[0]
			// Формируем минимальный JSON без лишних аллокаций.
			// Пример: ["height","180"]
			var buf bytes.Buffer
			buf.Grow(len(key.Syn) + len(trait.TValue) + 8)
			buf.WriteByte('[')
			buf.WriteString(strconv.Quote(key.Syn))
			buf.WriteByte(',')
			buf.WriteString(strconv.Quote(trait.TValue))
			buf.WriteByte(']')
			body := buf.Bytes()
			// Считаем и выставляем ETag (weak допустим). Если совпал — 304 без тела.
			sum := sha256.Sum256(body)
			etag := `W/"` + hex.EncodeToString(sum[:]) + `"`
			if inm := c.Get("If-None-Match"); inm == etag {
				c.Set("ETag", etag)
				return c.SendStatus(fiber.StatusNotModified)
			}
			c.Set("Content-Type", "application/json")
			c.Set("ETag", etag)
			return c.Send(body)
		}
		if _, ok := store.GetSet(id); !ok {
			return responses.Error(c, 404, "UUID не найден")
		}
		c.Set("Content-Type", "application/x-ndjson")
		c.Status(fiber.StatusOK)
		// Используем функцию обратного вызова для стриминга, чтобы не аллоцировать память под все результаты.
		err = CollectAllTraits(id, func(t *models.TraitResponse) error {
			jsonLine, err := json.Marshal(t)
			if err != nil {
				// Эта ошибка критична для стрима, поэтому нужно остановиться.
				return fmt.Errorf("не удалось преобразовать особенность в JSON: %w", err)
			}
			// Добавляем перевод строки для NDJSON и отправляем клиенту.
			// Если клиент отключится, c.Write вернёт ошибку и мы остановим обход.
			if _, err := c.Write(append(jsonLine, '\n')); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			// Ошибка, скорее всего, из-за отключения клиента. Логируем.
			log.Printf("Ошибка во время стрима набора для UUID %s: %v", uuidStr, err)
		}
		return nil
	}
}
// CollectAllTraits обходит все особенности в наборе, включая вложенные s_uuid.
//
// Параметры:
//
//	root — корневой UUID набора.
//	callback — вызывается для каждой найденной особенности.
//
// Возвращает:
//
//	Ошибку callback или ошибку обхода.
func CollectAllTraits(root uuid.UUID, callback func(*models.TraitResponse) error) error {
	seen := make(map[uuid.UUID]bool)
	queue := []uuid.UUID{root}
	seen[root] = true // сразу помечаем корневой узел как посещённый
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		// Проверяем, является ли текущий UUID набором
		if set, ok := store.GetSet(curr); ok {
			for _, child := range set.SChilds {
				if !seen[child] {
					seen[child] = true
					queue = append(queue, child)
				}
			}
			continue // переходим к следующему элементу очереди
		}
		// Если это не набор, проверяем, является ли это особенностью
		if trait, ok := store.GetTrait(curr); ok {
			if keys, keyOk := store.GetKeys(trait.TKey); keyOk && len(keys) > 0 {
				key := keys[0]
				resp := &models.TraitResponse{
					TUUID:  trait.TUUID,
					TKey:   key.Syn,
					TKeyID: trait.TKey,
					TValue: trait.TValue,
				}
				// Вызываем функцию обратного вызова. Если она вернёт ошибку (например, клиент отключился),
				// останавливаем обход и возвращаем ошибку.
				if err := callback(resp); err != nil {
					return err
				}
			} else {
				log.Printf("CollectAllTraits: несогласованность, ключ не найден для TKey %d особенности %s", trait.TKey, trait.TUUID)
			}
		} else {
			// UUID не является ни набором, ни особенностью (возможна битая ссылка в данных)
			log.Printf("CollectAllTraits: UUID %s не найден ни как набор, ни как особенность при обходе", curr)
		}
	}
	return nil
}
