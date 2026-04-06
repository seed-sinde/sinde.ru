package handlers

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"sinde.ru/db/services"
	"sinde.ru/internal/http/responses"
	"sinde.ru/internal/store"
	"sort"
	"strings"
)

type KeyMetaRequest struct {
	Syn string `json:"syn"`
}

// KeyMetaHandler возвращает метаданные ключа по `syn`.
//
// Маршрут:
//
//	POST /api/v1/keys/meta
//
// Тело запроса:
//
//	{"syn":"..."}
//
// Возвращает:
//
//	JSON-объект с id, syn и meta.
func KeyMetaHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req KeyMetaRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if req.Syn == "" {
			return responses.Error(c, fiber.StatusBadRequest, "ключ не указан", nil)
		}
		key, ok := store.GetKeyBySyn(req.Syn)
		if !ok {
			return responses.Error(c, fiber.StatusNotFound, "ключ не найден", nil)
		}
		var meta any = map[string]any{}
		if key.Meta != "" {
			_ = json.Unmarshal([]byte(key.Meta), &meta) // return {} if JSON is invalid
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{
			"id":   key.ID,
			"syn":  key.Syn,
			"meta": meta,
		})
	}
}

// KeyMetaAllHandler возвращает все уникальные метаданные для указанного `syn`.
//
// Маршрут:
//
//	POST /api/v1/keys/meta/all
//
// Тело запроса:
//
//	{"syn":"..."}
//
// Возвращает:
//
//	JSON-объект с syn и списком уникальных items.
func KeyMetaAllHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req KeyMetaRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if req.Syn == "" {
			return responses.Error(c, fiber.StatusBadRequest, "ключ не указан", nil)
		}
		keys, ok := store.GetKeysBySyn(req.Syn)
		if !ok {
			return responses.Error(c, fiber.StatusNotFound, "ключ не найден", nil)
		}
		unique := make(map[string]struct{})
		items := make([]fiber.Map, 0, len(keys))
		for _, key := range keys {
			if key == nil {
				continue
			}
			var meta any = map[string]any{}
			if key.Meta != "" {
				_ = json.Unmarshal([]byte(key.Meta), &meta)
			}
			dedupKey := key.Meta
			if dedupKey == "" {
				if raw, err := json.Marshal(meta); err == nil {
					dedupKey = string(raw)
				}
			}
			if dedupKey == "" {
				continue
			}
			if _, seen := unique[dedupKey]; seen {
				continue
			}
			unique[dedupKey] = struct{}{}
			items = append(items, fiber.Map{
				"id":   key.ID,
				"syn":  key.Syn,
				"meta": meta,
			})
		}
		if len(items) == 0 {
			return responses.Error(c, fiber.StatusNotFound, "ключ не найден", nil)
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{
			"syn":   req.Syn,
			"items": items,
		})
	}
}

type KeyMetaBulkRequest struct {
	IDs []int64 `json:"ids"`
}

// KeyMetaBulkHandler возвращает метаданные по списку ID с дедупликацией на стороне сервера.
//
// Маршрут:
//
//	POST /api/v1/keys/meta/bulk
//
// Тело запроса:
//
//	{"ids":[1,2,3]}
//
// Возвращает:
//
//	JSON-объект со списком items.
func KeyMetaBulkHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req KeyMetaBulkRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if len(req.IDs) == 0 {
			return responses.Success(c, fiber.StatusOK, fiber.Map{"items": []fiber.Map{}})
		}
		seen := make(map[int64]struct{})
		uniq := make([]int64, 0, len(req.IDs))
		for _, id := range req.IDs {
			if id == 0 {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			uniq = append(uniq, id)
		}
		keys, ok := store.GetKeys(uniq...)
		if !ok || len(keys) == 0 {
			return responses.Success(c, fiber.StatusOK, fiber.Map{"items": []fiber.Map{}})
		}
		items := make([]fiber.Map, 0, len(keys))
		for _, key := range keys {
			if key == nil {
				continue
			}
			var meta any = map[string]any{}
			if key.Meta != "" {
				_ = json.Unmarshal([]byte(key.Meta), &meta)
			}
			items = append(items, fiber.Map{
				"id":   key.ID,
				"syn":  key.Syn,
				"meta": meta,
			})
		}
		return responses.Success(c, fiber.StatusOK, fiber.Map{"items": items})
	}
}

type KeyMetaUpdateRequest struct {
	ID   int64          `json:"id"`
	Meta map[string]any `json:"meta"`
}
type KeyEnumOptionsRequest struct {
	Syn string `json:"syn"`
}

// KeyEnumOptionsHandler возвращает уникальные enum-значения для `syn`,
// объединяя значения по всем ключам с этим синонимом.
//
// Маршрут:
//
//	POST /api/v1/keys/enum-options
//
// Тело запроса:
//
//	{"syn":"..."}
//
// Возвращает:
//
//	JSON-объект со списком options.
func KeyEnumOptionsHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req KeyEnumOptionsRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if strings.TrimSpace(req.Syn) == "" {
			return responses.Error(c, fiber.StatusBadRequest, "ключ не указан", nil)
		}
		keys, ok := store.GetKeysBySyn(req.Syn)
		if !ok {
			return responses.Success(c, fiber.StatusOK, fiber.Map{"options": []string{}})
		}
		uniq := make(map[string]struct{})
		for _, k := range keys {
			if k == nil {
				continue
			}
			for _, v := range store.GetTraitValuesByKey(k.ID) {
				uniq[v] = struct{}{}
			}
		}
		out := make([]string, 0, len(uniq))
		for v := range uniq {
			out = append(out, v)
		}
		sort.Strings(out)
		return responses.Success(c, fiber.StatusOK, fiber.Map{"options": out})
	}
}

// KeyMetaUpdateHandler обновляет метаданные ключа по его ID.
//
// Маршрут:
//
//	PATCH /api/v1/keys/meta
//
// Тело запроса:
//
//	{"id":123,"meta":{...}}
//
// Возвращает:
//
//	JSON-объект обновлённого ключа.
func KeyMetaUpdateHandler() fiber.Handler {
	return func(c fiber.Ctx) error {
		var req KeyMetaUpdateRequest
		if err := c.Bind().Body(&req); err != nil {
			return responses.Error(c, fiber.StatusBadRequest, "некорректный запрос", err.Error())
		}
		if req.ID == 0 {
			return responses.Error(c, fiber.StatusBadRequest, "ID ключа обязателен", nil)
		}
		if req.Meta == nil {
			req.Meta = map[string]any{}
		}
		metaBytes, err := json.Marshal(req.Meta)
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "ошибка сериализации метаданных", err.Error())
		}
		key, err := services.PDBUpdateKeyMetaByID(c, req.ID, string(metaBytes))
		if err != nil {
			return responses.Error(c, fiber.StatusInternalServerError, "ошибка обновления базы данных", err.Error())
		}
		if key == nil {
			return responses.Error(c, fiber.StatusNotFound, "ключ не найден", nil)
		}
		store.SyncKey(key)
		return responses.Success(c, fiber.StatusOK, fiber.Map{
			"id":   key.ID,
			"syn":  key.Syn,
			"meta": req.Meta,
		})
	}
}
