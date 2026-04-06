package handlers

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"log"
	"sinde.ru/db/services"
	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/models"
	"sinde.ru/internal/store"
	"sinde.ru/utils"
	"strings"
)

// CreateSetHandler создаёт или находит набор из двух потомков.
//
// Маршрут:
//
//	POST /api/v1/sets
//
// Тело запроса:
//
//	{"s_childs":["<t|s uuid>","<t uuid>"]}
//	{"s_childs":["<t|s uuid>",{"t_key":"syn","t_value":"val"}]}
//
// Возвращает:
//
//	201 + {"s_uuid":"..."} (и "t_uuid", если особенность создана из объекта).
func CreateSetHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		// верхний уровень: строгое декодирование
		var in struct {
			SChilds []json.RawMessage `json:"s_childs"`
		}
		dec := json.NewDecoder(bytes.NewReader(c.Body()))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&in); err != nil {
			return responses.Error(c, 400, "Некорректное тело запроса", err.Error())
		}
		if len(in.SChilds) != 2 {
			return responses.Error(c, 400, `Тело должно быть формата {"s_childs":[left,right]}`)
		}
		// --- LEFT: UUID-строка (t_uuid или s_uuid), уже должна существовать в Postgres
		var leftStr string
		if err := json.Unmarshal(in.SChilds[0], &leftStr); err != nil || leftStr == "" {
			return responses.Error(c, 400, "левый потомок должен быть UUID-строкой")
		}
		left, err := uuid.Parse(leftStr)
		if err != nil {
			return responses.Error(c, 400, "левый потомок должен быть UUID-строкой")
		}
		leftExists, err := services.PdbTraitTargetExists(c, left)
		if err != nil {
			return responses.Error(c, 500, "Не удалось проверить левый потомок", err.Error())
		}
		if !leftExists {
			return responses.Error(c, 400, "левый потомок должен существовать как t_uuid или s_uuid")
		}
		// --- RIGHT: либо UUID-строка (существующая особенность), либо объект {t_key,t_value} без лишних полей
		var rightTUUID uuid.UUID
		var echoTraitUUID string
		// сначала пробуем как UUID-строку
		var tryUUID string
		if err := json.Unmarshal(in.SChilds[1], &tryUUID); err == nil && tryUUID != "" {
			id, err := uuid.Parse(tryUUID)
			if err != nil {
				return responses.Error(c, 400, "правый потомок должен быть существующим t_uuid (особенность)")
			}
			rightExists, err := services.PdbTraitExists(c, id)
			if err != nil {
				return responses.Error(c, 500, "Не удалось проверить правый потомок", err.Error())
			}
			if !rightExists {
				return responses.Error(c, 400, "правый потомок должен быть существующим t_uuid (особенность)")
			}
			rightTUUID = id
		} else {
			var kv models.TraitKV
			decKV := json.NewDecoder(bytes.NewReader(in.SChilds[1]))
			decKV.DisallowUnknownFields()
			if err := decKV.Decode(&kv); err != nil {
				return responses.Error(c, 400, "Некорректный правый потомок", err.Error())
			}
			kv.TKey = strings.TrimSpace(kv.TKey)
			if kv.TKey == "" || kv.TValue == nil {
				return responses.Error(c, 400, "правый потомок должен содержать и t_key, и t_value")
			}
			// нормализуем value к строке
			val := fmt.Sprintf("%v", kv.TValue)
			key, err := services.PdbGetOrCreateCanonicalKeyBySyn(c, kv.TKey)
			if err != nil {
				return responses.Error(c, 500, "Не удалось сохранить ключ правого потомка", err.Error())
			}
			if key == nil {
				return responses.Error(c, 400, "Некорректный ключ правого потомка")
			}
			store.SyncKey(key)
			trait, err := services.PdbGetOrCreateTraitByKeyValue(c, key.ID, val)
			if err != nil {
				return responses.Error(c, 500, "Не удалось сохранить правую особенность", err.Error())
			}
			if trait == nil {
				return responses.Error(c, 500, "Не удалось сохранить правую особенность")
			}
			store.SyncTrait(trait)
			rightTUUID = trait.TUUID
			echoTraitUUID = trait.TUUID.String()
		}
		setItem, err := services.PdbGetOrCreateSet(c, []uuid.UUID{left, rightTUUID})
		if err != nil {
			return responses.Error(c, 500, "Не удалось создать или найти набор", err.Error())
		}
		if setItem == nil {
			return responses.Error(c, 500, "Не удалось создать или найти набор")
		}
		store.SyncSet(setItem)
		out := fiber.Map{"s_uuid": setItem.SUUID.String()}
		if echoTraitUUID != "" {
			out["t_uuid"] = echoTraitUUID
		}
		return responses.Success(c, 201, out)
	}
}

// FindOrBuildSetHandler строит или находит цепочку наборов из списка t_uuid.
//
// Маршрут:
//
//	POST /api/v1/sets/find-or-build
//
// Тело запроса:
//
//	["<t_uuid>", "<t_uuid>", ...]
//
// Возвращает:
//
//	JSON-объект с `s_uuid` верхнего набора цепочки.
func FindOrBuildSetHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		var childsRaw []string
		if err := c.Bind().Body(&childsRaw); err != nil {
			return responses.Error(c, 400, "Некорректное тело запроса", err.Error())
		}
		if len(childsRaw) < 2 {
			return responses.Error(c, 400, "Нужно минимум два t_uuid")
		}
		childs := make([]uuid.UUID, 0, len(childsRaw))
		// В цепочке разрешены только особенности
		for i, u := range childsRaw {
			id, err := uuid.Parse(u)
			if err != nil {
				return responses.Error(c, 400, fmt.Sprintf("Потомок #%d не является существующим t_uuid", i))
			}
			exists, err := services.PdbTraitExists(c, id)
			if err != nil {
				return responses.Error(c, 500, "Не удалось проверить потомка", err.Error())
			}
			if !exists {
				return responses.Error(c, 400, fmt.Sprintf("Потомок #%d не является существующим t_uuid", i))
			}
			childs = append(childs, id)
		}
		current := childs[0]
		for i := 1; i < len(childs); i++ {
			setItem, err := services.PdbGetOrCreateSet(c, []uuid.UUID{current, childs[i]})
			if err != nil {
				return responses.Error(c, 500, "Не удалось построить или найти цепочку наборов", err.Error())
			}
			if setItem == nil {
				return responses.Error(c, 500, "Не удалось построить или найти цепочку наборов")
			}
			store.SyncSet(setItem)
			current = setItem.SUUID
		}
		return responses.Success(c, 200, fiber.Map{"s_uuid": current.String()})
	}
}

// StoreSetHandler возвращает набор с раскрытыми особенностями.
//
// Маршрут:
//
//	GET /store/sets/:uuid
//
// Возвращает:
//
//	JSON-объект SetResponse со списком traits.
func StoreSetHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer utils.Benchmark("storeSetHandler")()
		uuidStr := c.Params("uuid")
		if uuidStr == "" {
			return responses.Error(c, fiber.StatusBadRequest, "UUID обязателен")
		}
		id, err := uuid.Parse(uuidStr)
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "UUID некорректен")
		}
		if _, ok := store.GetSet(id); !ok {
			return responses.Error(c, fiber.StatusNotFound, "Набор не найден")
		}
		// Динамически собираем все дочерние особенности
		allTraits := make([]*models.TraitResponse, 0)
		err = CollectAllTraits(id, func(t *models.TraitResponse) error {
			allTraits = append(allTraits, t)
			return nil
		})
		if err != nil {
			log.Printf("Ошибка при сборе особенностей для набора %s: %v", uuidStr, err)
			return responses.Error(c, 500, "Не удалось собрать особенности")
		}
		resp := models.SetResponse{
			SUUID:  id,
			Traits: allTraits,
		}
		return c.JSON(resp)
	}
}

// StoreSetStreamHandler стримит особенности набора в формате NDJSON.
//
// Маршрут:
//
//	GET /store/sets/:uuid/stream
//
// Возвращает:
//
//	NDJSON-поток, по одной особенности на строку.
func StoreSetStreamHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer utils.Benchmark("storeSetStreamHandler")()
		uuidStr := c.Params("uuid")
		if uuidStr == "" {
			return responses.Error(c, fiber.StatusBadRequest, "UUID обязателен")
		}
		id, err := uuid.Parse(uuidStr)
		if err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "UUID некорректен")
		}
		if _, ok := store.GetSet(id); !ok {
			return responses.Error(c, fiber.StatusNotFound, "Набор не найден")
		}
		// Заголовок Content-Type для NDJSON-стриминга. Остальные заголовки (кэш, соединение) обрабатывает Caddy.
		c.Set("Content-Type", "application/x-ndjson")
		// Используем CollectAllTraits для стриминга
		err = CollectAllTraits(id, func(trait *models.TraitResponse) error {
			jsonLine, err := json.Marshal(trait)
			if err != nil {
				return fmt.Errorf("не удалось преобразовать особенность в JSON: %w", err)
			}
			if _, err := c.Write(append(jsonLine, '\n')); err != nil {
				return err // Останавливаем стриминг при отключении клиента.
			}
			return nil
		})
		if err != nil {
			log.Printf("Ошибка во время стрима набора для UUID %s: %v", uuidStr, err)
			// Не отправляем ошибку клиенту, если стриминг уже начался
		}
		return nil
	}
}
