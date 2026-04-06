package store

import (
	"fmt"
	"github.com/google/uuid"
	"sinde.ru/internal/models"
	"strings"
)

func joinChilds(childs []uuid.UUID) string {
	if len(childs) == 0 {
		return ""
	}
	// важно: порядок имеет значение; используем '\x1f' как редкий разделитель
	parts := make([]string, 0, len(childs))
	for _, c := range childs {
		parts = append(parts, c.String())
	}
	return strings.Join(parts, "\x1f")
}
func GetSet(id uuid.UUID) (*models.Set, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.Sets[id]
	return st, ok
}

// SyncSet синхронизирует набор из Postgres в in-memory store.
func SyncSet(set *models.Set) *models.Set {
	if set == nil || set.SUUID == uuid.Nil {
		return nil
	}
	copySet := *set
	copySet.SChilds = append([]uuid.UUID(nil), set.SChilds...)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Sets[copySet.SUUID] = &copySet
	if key := joinChilds(copySet.SChilds); key != "" {
		s.SetChildsToSUUID[key] = copySet.SUUID
	}
	return s.Sets[copySet.SUUID]
}

// Создаёт или находит набор РОВНО из двух потомков.
// Возвращает: s_uuid, (nil, так как здесь особенность не создаётся), *models.Set при создании нового набора, ошибка.
func CreateOrFindSet(children []any) (uuid.UUID, *models.Trait, *models.Set, error) {
	if len(children) != 2 {
		return uuid.Nil, nil, nil, fmt.Errorf("a set must be created from exactly two children")
	}
	// оба элемента должны быть UUID-строками
	a, ok1 := children[0].(uuid.UUID)
	b, ok2 := children[1].(uuid.UUID)
	if !ok1 || !ok2 {
		return uuid.Nil, nil, nil, fmt.Errorf("children must be uuid strings")
	}
	key := joinChilds([]uuid.UUID{a, b})
	s.mu.Lock()
	defer s.mu.Unlock()
	if suuid, ok := s.SetChildsToSUUID[key]; ok {
		return suuid, nil, nil, nil
	}
	suuid := s.generateFreeUUIDLocked()
	st := &models.Set{SUUID: suuid, SChilds: []uuid.UUID{a, b}}
	s.Sets[suuid] = st
	s.SetChildsToSUUID[key] = suuid
	return suuid, nil, st, nil
}

// Лево-ассоциированная цепочка из произвольного списка UUID.
// Возвращает верхний s_uuid и список НОВЫХ наборов для асинхронной персистентности.
func FindOrBuildSetChain(childs []uuid.UUID) (uuid.UUID, []*models.Set, error) {
	if len(childs) == 0 {
		return uuid.Nil, nil, nil
	}
	current := childs[0]
	var created []*models.Set
	for i := 1; i < len(childs); i++ {
		next := childs[i]
		joined := joinChilds([]uuid.UUID{current, next})
		s.mu.Lock()
		if suuid, ok := s.SetChildsToSUUID[joined]; ok {
			current = suuid
			s.mu.Unlock()
			continue
		}
		suuid := s.generateFreeUUIDLocked()
		st := &models.Set{SUUID: suuid, SChilds: []uuid.UUID{current, next}}
		s.Sets[suuid] = st
		s.SetChildsToSUUID[joined] = suuid
		s.mu.Unlock()
		created = append(created, st)
		current = suuid
	}
	return current, created, nil
}
