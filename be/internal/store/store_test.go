package store
import (
	"encoding/json"
	"reflect"
	"testing"
	"github.com/google/uuid"
	"sinde.ru/internal/models"
)
func resetStoreForTest(t *testing.T) {
	t.Helper()
	old := s
	s = &store{
		Keys:             make(map[int64]*models.Key),
		Traits:           make(map[uuid.UUID]*models.Trait),
		Sets:             make(map[uuid.UUID]*models.Set),
		SynToKeyIDs:      make(map[string][]int64),
		SetChildsToSUUID: make(map[string]uuid.UUID),
		TraitResponses:   make(map[uuid.UUID]*models.TraitResponse),
		TraitKVToUUID:    make(map[string]uuid.UUID),
	}
	t.Cleanup(func() {
		s = old
	})
}
func TestGetOrAddKeyBySynAndUpdateMeta(t *testing.T) {
	resetStoreForTest(t)
	key := GetOrAddKeyBySyn("  Color ")
	if key == nil {
		t.Fatal("expected key to be created")
	}
	if key.ID != 1 {
		t.Fatalf("expected first key id to be 1, got %d", key.ID)
	}
	if key.Syn != "Color" {
		t.Fatalf("unexpected key syn: %q", key.Syn)
	}
	if key.Meta != "{}" {
		t.Fatalf("unexpected default meta: %q", key.Meta)
	}
	sameKey := GetOrAddKeyBySyn("color")
	if sameKey == nil || sameKey.ID != key.ID {
		t.Fatalf("expected existing key to be reused, got %#v", sameKey)
	}
	updated, ok := UpdateKeyMetaByID(key.ID, map[string]any{"type": "string"})
	if !ok {
		t.Fatal("expected key meta update to succeed")
	}
	var decoded map[string]any
	if err := json.Unmarshal(updated.RawJSON, &decoded); err != nil {
		t.Fatalf("failed to decode raw json: %v", err)
	}
	if decoded["type"] != "string" {
		t.Fatalf("unexpected decoded meta: %#v", decoded)
	}
}
func TestLoadInitialDataSnapshotAndTraitLookup(t *testing.T) {
	resetStoreForTest(t)
	key := &models.Key{ID: 10, Syn: "Color", Meta: `{"kind":"enum"}`, RawJSON: []byte(`{"kind":"enum"}`)}
	traitID1 := uuid.New()
	traitID2 := uuid.New()
	setID := uuid.New()
	LoadInitialData(
		[]*models.Trait{
			{TUUID: traitID1, TKey: 10, TValue: "red"},
			{TUUID: traitID2, TKey: 10, TValue: "red"},
		},
		[]*models.Key{key},
		[]*models.Set{{SUUID: setID, SChilds: []uuid.UUID{traitID1, traitID2}}},
	)
	foundKey, ok := GetKeyBySyn(" color ")
	if !ok || foundKey == nil || foundKey.ID != 10 {
		t.Fatalf("expected key lookup by syn to succeed, got %#v", foundKey)
	}
	values := GetTraitValuesByKey(10)
	if len(values) != 1 || values[0] != "red" {
		t.Fatalf("unexpected unique trait values: %#v", values)
	}
	snapshot := SnapshotData()
	if len(snapshot.Keys) != 1 || len(snapshot.Traits) != 2 || len(snapshot.Sets) != 1 {
		t.Fatalf("unexpected snapshot sizes: %#v", snapshot)
	}
	snapshot.Keys[0].RawJSON[0] = '['
	snapshot.Sets[0].SChilds[0] = uuid.Nil
	if s.Keys[10].RawJSON[0] != '{' {
		t.Fatalf("expected key raw json to remain unchanged, got %q", string(s.Keys[10].RawJSON))
	}
	if s.Sets[setID].SChilds[0] == uuid.Nil {
		t.Fatal("expected set childs to be copied in snapshot")
	}
}
func TestCreateOrFindSetChainAndTraitReuse(t *testing.T) {
	resetStoreForTest(t)
	key := &models.Key{ID: 1, Syn: "Color", Meta: "{}"}
	s.Keys[key.ID] = key
	s.SynToKeyIDs[normalizeSyn(key.Syn)] = []int64{key.ID}
	s.maxKeyID = key.ID
	left := uuid.New()
	right := uuid.New()
	last := uuid.New()
	setID, _, createdSet, err := CreateOrFindSet([]any{left, right})
	if err != nil {
		t.Fatalf("CreateOrFindSet returned error: %v", err)
	}
	if createdSet == nil || createdSet.SUUID != setID {
		t.Fatalf("expected created set, got %#v", createdSet)
	}
	sameSetID, _, noSet, err := CreateOrFindSet([]any{left, right})
	if err != nil {
		t.Fatalf("CreateOrFindSet returned error: %v", err)
	}
	if sameSetID != setID {
		t.Fatalf("expected existing set id %s, got %s", setID, sameSetID)
	}
	if noSet != nil {
		t.Fatalf("expected no new set on reuse, got %#v", noSet)
	}
	topID, created, err := FindOrBuildSetChain([]uuid.UUID{left, right, last})
	if err != nil {
		t.Fatalf("FindOrBuildSetChain returned error: %v", err)
	}
	if topID == uuid.Nil {
		t.Fatal("expected non-nil top set id")
	}
	if len(created) != 1 {
		t.Fatalf("expected one new set in chain, got %d", len(created))
	}
	if !reflect.DeepEqual(created[0].SChilds, []uuid.UUID{setID, last}) {
		t.Fatalf("unexpected created chain set childs: %#v", created[0].SChilds)
	}
	resp1 := GetOrAddTraitByKeyValue(1, "red")
	resp2 := GetOrAddTraitByKeyValue(1, "red")
	if resp1.TUUID != resp2.TUUID {
		t.Fatalf("expected trait reuse for identical key/value, got %s and %s", resp1.TUUID, resp2.TUUID)
	}
	if resp1.TKey != "Color" {
		t.Fatalf("unexpected trait response key: %q", resp1.TKey)
	}
}
