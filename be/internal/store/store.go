package store

import (
	"github.com/google/uuid"
	"sinde.ru/internal/models"
	"sync"
)

type store struct {
	mu sync.RWMutex
	// основные данные
	Keys   map[int64]*models.Key
	Traits map[uuid.UUID]*models.Trait
	Sets   map[uuid.UUID]*models.Set
	// счётчики
	maxKeyID int64
	// индексы / кэши
	SynToKeyIDs      map[string][]int64
	SetChildsToSUUID map[string]uuid.UUID
	TraitResponses   map[uuid.UUID]*models.TraitResponse
	// быстрый поиск по паре (t_key, t_value)
	TraitKVToUUID map[string]uuid.UUID
}

var s = &store{
	Keys:             make(map[int64]*models.Key),
	Traits:           make(map[uuid.UUID]*models.Trait),
	Sets:             make(map[uuid.UUID]*models.Set),
	SynToKeyIDs:      make(map[string][]int64),
	SetChildsToSUUID: make(map[string]uuid.UUID),
	TraitResponses:   make(map[uuid.UUID]*models.TraitResponse),
	TraitKVToUUID:    make(map[string]uuid.UUID),
}

// Функция GenerateFreeUUID возвращает уникальный UUID, который не конфликтует с особенностями или наборами.
func (st *store) GenerateFreeUUID() uuid.UUID {
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.generateFreeUUIDLocked()
}
func (st *store) generateFreeUUIDLocked() uuid.UUID {
	for {
		u := uuid.New()
		if _, ok := st.Traits[u]; ok {
			continue
		}
		if _, ok := st.Sets[u]; ok {
			continue
		}
		return u
	}
}
