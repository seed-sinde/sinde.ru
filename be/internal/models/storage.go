package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type StorageObject struct {
	ObjectID    uuid.UUID       `json:"object_id"`
	StorageKey  string          `json:"storage_key"`
	BucketName  string          `json:"bucket_name"`
	MediaFamily string          `json:"media_family"`
	ContentType string          `json:"content_type"`
	ByteSize    int64           `json:"byte_size"`
	FileHash    string          `json:"file_hash,omitempty"`
	SourceKind  string          `json:"source_kind"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type StorageObjectUsage struct {
	UsageID    uuid.UUID       `json:"usage_id"`
	ObjectID   uuid.UUID       `json:"object_id"`
	EntityType string          `json:"entity_type"`
	EntityID   string          `json:"entity_id"`
	UsageType  string          `json:"usage_type"`
	FieldName  string          `json:"field_name"`
	SortOrder  int             `json:"sort_order"`
	IsPrimary  bool            `json:"is_primary"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	Object     *StorageObject  `json:"object,omitempty"`
}
