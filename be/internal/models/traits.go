package models

import "github.com/google/uuid"

type Trait struct {
	TUUID  uuid.UUID `json:"t_uuid"`  // UUID
	TKey   int64     `json:"t_key"`   // foreign key (ID from traits_k)
	TValue string    `json:"t_value"` // value
}
type Key struct {
	ID      int64  `json:"id"`     // PK
	SynID   int64  `json:"syn_id"` // FK to key_syns
	Syn     string `json:"syn"`    // human-friendly key name
	Meta    string `json:"meta"`   // JSON description
	RawJSON []byte `json:"-"`
}
type Set struct {
	SUUID   uuid.UUID   `json:"s_uuid"`
	SChilds []uuid.UUID `json:"s_childs"` // array of t_uuid or s_uuid (order matters)
}
type TraitKV struct {
	TKey   string `json:"t_key"`
	TValue any    `json:"t_value"`
}
type TraitResponse struct {
	TUUID  uuid.UUID `json:"t_uuid"`
	TKey   string    `json:"t_key"` // human-readable syn
	TKeyID int64     `json:"t_key_id"`
	TValue string    `json:"t_value"`
}
type SetResponse struct {
	SUUID  uuid.UUID        `json:"s_uuid"`
	Traits []*TraitResponse `json:"traits"`
}
