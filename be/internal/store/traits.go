package store
import (
	"github.com/google/uuid"
	"sinde.ru/internal/models"
)
func GetTrait(id uuid.UUID) (*models.Trait, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.Traits[id]
	return t, ok
}
// GetOrAddTraitByKeyValue ищет особенность по паре (keyID, value) или создаёт новую.
//
// Параметры:
//
//	keyID — идентификатор ключа.
//	value — значение особенности.
//
// Возвращает:
//
//	TraitResponse для найденной или созданной особенности.
func GetOrAddTraitByKeyValue(keyID int64, value string) *models.TraitResponse {
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := traitIndexKey(keyID, value)
	if existingUUID, ok := s.TraitKVToUUID[idx]; ok {
		if resp, ok := s.TraitResponses[existingUUID]; ok {
			return resp
		}
		if ex, ok := s.Traits[existingUUID]; ok {
			var syn string
			if k := s.Keys[ex.TKey]; k != nil {
				syn = k.Syn
			}
			resp := &models.TraitResponse{
				TUUID:  existingUUID,
				TKey:   syn,
				TKeyID: ex.TKey,
				TValue: ex.TValue,
			}
			s.TraitResponses[existingUUID] = resp
			return resp
		}
	}
	tuuid := s.generateFreeUUIDLocked()
	t := &models.Trait{TUUID: tuuid, TKey: keyID, TValue: value}
	s.Traits[tuuid] = t
	s.TraitKVToUUID[idx] = tuuid
	var syn string
	if k := s.Keys[keyID]; k != nil {
		syn = k.Syn
	}
	resp := &models.TraitResponse{
		TUUID:  tuuid,
		TKey:   syn,
		TKeyID: keyID,
		TValue: value,
	}
	s.TraitResponses[tuuid] = resp
	return resp
}
// SyncTrait синхронизирует особенность из Postgres в in-memory store.
func AddTrait(t *models.Trait) *models.TraitResponse {
	return SyncTrait(t)
}
func SyncTrait(t *models.Trait) *models.TraitResponse {
	if t == nil || t.TUUID == uuid.Nil {
		return nil
	}
	copyTrait := *t
	s.mu.Lock()
	defer s.mu.Unlock()
	idx := traitIndexKey(copyTrait.TKey, copyTrait.TValue)
	s.Traits[copyTrait.TUUID] = &copyTrait
	s.TraitKVToUUID[idx] = copyTrait.TUUID
	var syn string
	if k := s.Keys[copyTrait.TKey]; k != nil {
		syn = k.Syn
	}
	resp := &models.TraitResponse{
		TUUID:  copyTrait.TUUID,
		TKey:   syn,
		TKeyID: copyTrait.TKey,
		TValue: copyTrait.TValue,
	}
	s.TraitResponses[copyTrait.TUUID] = resp
	return resp
}
