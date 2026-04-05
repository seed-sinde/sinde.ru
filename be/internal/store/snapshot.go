package store
import (
	"sort"
	"github.com/google/uuid"
	"sinde.ru/internal/models"
)
type Snapshot struct {
	Keys   []*models.Key
	Traits []*models.Trait
	Sets   []*models.Set
}
func SnapshotData() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keyIDs := make([]int64, 0, len(s.Keys))
	for id := range s.Keys {
		keyIDs = append(keyIDs, id)
	}
	sort.Slice(keyIDs, func(i, j int) bool { return keyIDs[i] < keyIDs[j] })
	keys := make([]*models.Key, 0, len(keyIDs))
	for _, id := range keyIDs {
		key := s.Keys[id]
		if key == nil {
			continue
		}
		copyKey := *key
		if key.RawJSON != nil {
			copyKey.RawJSON = append([]byte(nil), key.RawJSON...)
		}
		keys = append(keys, &copyKey)
	}
	traitIDs := make([]uuid.UUID, 0, len(s.Traits))
	for id := range s.Traits {
		traitIDs = append(traitIDs, id)
	}
	sort.Slice(traitIDs, func(i, j int) bool { return traitIDs[i].String() < traitIDs[j].String() })
	traits := make([]*models.Trait, 0, len(traitIDs))
	for _, id := range traitIDs {
		trait := s.Traits[id]
		if trait == nil {
			continue
		}
		copyTrait := *trait
		traits = append(traits, &copyTrait)
	}
	setIDs := make([]uuid.UUID, 0, len(s.Sets))
	for id := range s.Sets {
		setIDs = append(setIDs, id)
	}
	sort.Slice(setIDs, func(i, j int) bool { return setIDs[i].String() < setIDs[j].String() })
	sets := make([]*models.Set, 0, len(setIDs))
	for _, id := range setIDs {
		setItem := s.Sets[id]
		if setItem == nil {
			continue
		}
		copySet := *setItem
		copySet.SChilds = append([]uuid.UUID(nil), setItem.SChilds...)
		sets = append(sets, &copySet)
	}
	return Snapshot{
		Keys:   keys,
		Traits: traits,
		Sets:   sets,
	}
}
