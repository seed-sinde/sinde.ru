package store

import (
	"encoding/json"
	"sinde.ru/internal/models"
	"strings"
)

func appendUniqueInt64(items []int64, target int64) []int64 {
	for _, item := range items {
		if item == target {
			return items
		}
	}
	return append(items, target)
}
func removeInt64(items []int64, target int64) []int64 {
	if len(items) == 0 {
		return items
	}
	out := items[:0]
	for _, item := range items {
		if item != target {
			out = append(out, item)
		}
	}
	return out
}

// GetTraitValuesByKey возвращает все уникальные значения t_value для указанного keyID.
//
// Параметры:
//
//	keyID — идентификатор ключа.
//
// Возвращает:
//
//	Срез уникальных значений t_value.
func GetTraitValuesByKey(keyID int64) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	uniq := make(map[string]struct{})
	for _, t := range s.Traits {
		if t != nil && t.TKey == keyID {
			uniq[t.TValue] = struct{}{}
		}
	}
	res := make([]string, 0, len(uniq))
	for v := range uniq {
		res = append(res, v)
	}
	return res
}

// GetKeys возвращает ключи по одному или нескольким ID.
//
// Параметры:
//
//	ids — идентификаторы ключей.
//
// Возвращает:
//
//	Список найденных ключей и признак успешного поиска.
func GetKeys(ids ...int64) ([]*models.Key, bool) {
	if len(ids) == 0 {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]*models.Key, 0, len(ids))
	for _, id := range ids {
		if k := s.Keys[id]; k != nil {
			res = append(res, k)
		}
	}
	if len(res) == 0 {
		return nil, false
	}
	return res, true
}

// SyncKey синхронизирует ключ из Postgres в in-memory store.
func SyncKey(key *models.Key) *models.Key {
	if key == nil || key.ID == 0 {
		return nil
	}
	copyKey := *key
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing := s.Keys[key.ID]; existing != nil {
		oldSyn := normalizeSyn(existing.Syn)
		newSyn := normalizeSyn(copyKey.Syn)
		if oldSyn != "" && oldSyn != newSyn {
			s.SynToKeyIDs[oldSyn] = removeInt64(s.SynToKeyIDs[oldSyn], key.ID)
			if len(s.SynToKeyIDs[oldSyn]) == 0 {
				delete(s.SynToKeyIDs, oldSyn)
			}
		}
	}
	s.Keys[key.ID] = &copyKey
	if key.ID > s.maxKeyID {
		s.maxKeyID = key.ID
	}
	if syn := normalizeSyn(copyKey.Syn); syn != "" {
		s.SynToKeyIDs[syn] = appendUniqueInt64(s.SynToKeyIDs[syn], copyKey.ID)
	}
	return s.Keys[key.ID]
}

// / GetOrAddKeyBySyn возвращает keyID по синониму (без учёта регистра) или создаёт его в памяти.
// (Запись в БД выполняется асинхронно на уровне service layer.)
// раньше: func GetOrAddKeyBySyn(syn string) (int, bool)
func GetOrAddKeyBySyn(syn string) *models.Key {
	trimmed := strings.TrimSpace(syn)
	lc := strings.ToLower(trimmed)
	if lc == "" {
		return nil
	}
	s.mu.RLock()
	if ids, ok := s.SynToKeyIDs[lc]; ok && len(ids) > 0 {
		minID := ids[0]
		for _, id := range ids[1:] {
			if id < minID {
				minID = id
			}
		}
		k := s.Keys[minID]
		s.mu.RUnlock()
		return k
	}
	s.mu.RUnlock()
	s.mu.Lock()
	defer s.mu.Unlock()
	if ids, ok := s.SynToKeyIDs[lc]; ok && len(ids) > 0 {
		minID := ids[0]
		for _, id := range ids[1:] {
			if id < minID {
				minID = id
			}
		}
		return s.Keys[minID]
	}
	// вычисляем новый int64 ID без полного сканирования
	newID := s.maxKeyID + 1
	s.maxKeyID = newID
	k := &models.Key{ID: newID, Syn: trimmed, Meta: "{}"}
	s.Keys[newID] = k
	s.SynToKeyIDs[lc] = append(s.SynToKeyIDs[lc], newID)
	return k
}

// GetKeyBySyn ищет первый ключ по syn.
//
// Параметры:
//
//	syn — синоним ключа.
//
// Возвращает:
//
//	Ключ с наименьшим ID для syn и признак успешного поиска.
func GetKeyBySyn(syn string) (*models.Key, bool) {
	lc := normalizeSyn(syn)
	if lc == "" {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids, ok := s.SynToKeyIDs[lc]
	if !ok || len(ids) == 0 {
		return nil, false
	}
	minKey := (*models.Key)(nil)
	for _, id := range ids {
		if k := s.Keys[id]; k != nil {
			if minKey == nil || k.ID < minKey.ID {
				minKey = k
			}
		}
	}
	if minKey == nil {
		return nil, false
	}
	return minKey, true
}

// GetKeysBySyn возвращает все ключи с указанным syn.
//
// Параметры:
//
//	syn — синоним ключа.
//
// Возвращает:
//
//	Список ключей и признак успешного поиска.
func GetKeysBySyn(syn string) ([]*models.Key, bool) {
	lc := normalizeSyn(syn)
	if lc == "" {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := s.SynToKeyIDs[lc]
	if len(ids) == 0 {
		return nil, false
	}
	res := make([]*models.Key, 0, len(ids))
	for _, id := range ids {
		if k := s.Keys[id]; k != nil {
			res = append(res, k)
		}
	}
	if len(res) == 0 {
		return nil, false
	}
	return res, true
}

// UpdateKeyMetaByID полностью заменяет meta у ключа по ID.
//
// Параметры:
//
//	id — идентификатор ключа.
//	meta — новые метаданные.
//
// Возвращает:
//
//	Обновлённый ключ и признак успешного обновления.
func UpdateKeyMetaByID(id int64, meta map[string]any) (*models.Key, bool) {
	if id == 0 {
		return nil, false
	}
	b, err := json.Marshal(meta)
	if err != nil {
		return nil, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	k, ok := s.Keys[id]
	if !ok || k == nil {
		return nil, false
	}
	k.Meta = string(b)
	k.RawJSON = b
	return k, true
}
