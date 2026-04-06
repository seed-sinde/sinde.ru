package store

import (
	"fmt"
	"github.com/google/uuid"
	"sinde.ru/internal/models"
	"strings"
)

func traitIndexKey(keyID int64, value string) string {
	return fmt.Sprintf("%d\x1f%s", keyID, value)
}
func normalizeSyn(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// LoadInitialData атомарно сбрасывает и заполняет in-memory store.
//
// Параметры:
//
//	traits — список особенностей.
//	keys — список ключей.
//	sets — список наборов.
//
// Возвращает:
//
//	Ничего; обновляет глобальное хранилище в памяти.
func LoadInitialData(traits []*models.Trait, keys []*models.Key, sets []*models.Set) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// сброс
	s.maxKeyID = 0
	s.Keys = make(map[int64]*models.Key)
	s.Traits = make(map[uuid.UUID]*models.Trait)
	s.Sets = make(map[uuid.UUID]*models.Set)
	s.SynToKeyIDs = make(map[string][]int64)
	s.SetChildsToSUUID = make(map[string]uuid.UUID)
	s.TraitResponses = make(map[uuid.UUID]*models.TraitResponse)
	s.TraitKVToUUID = make(map[string]uuid.UUID)
	// ключи
	for _, k := range keys {
		if k == nil {
			continue
		}
		s.Keys[k.ID] = k
		if k.ID > s.maxKeyID {
			s.maxKeyID = k.ID
		}
		if syn := strings.ToLower(k.Syn); syn != "" {
			s.SynToKeyIDs[syn] = append(s.SynToKeyIDs[syn], k.ID)
		}
	}
	// особенности
	for _, t := range traits {
		if t == nil || t.TUUID == uuid.Nil {
			continue
		}
		s.Traits[t.TUUID] = t
		s.TraitKVToUUID[traitIndexKey(t.TKey, t.TValue)] = t.TUUID
		var syn string
		if k := s.Keys[t.TKey]; k != nil {
			syn = k.Syn
		}
		s.TraitResponses[t.TUUID] = &models.TraitResponse{
			TUUID:  t.TUUID,
			TKey:   syn, // syn stored here as well
			TKeyID: t.TKey,
			TValue: t.TValue,
		}
	}
	// наборы
	for _, set := range sets {
		if set == nil || set.SUUID == uuid.Nil {
			continue
		}
		s.Sets[set.SUUID] = set
		key := joinChilds(set.SChilds)
		if key != "" {
			s.SetChildsToSUUID[key] = set.SUUID
		}
	}
}
