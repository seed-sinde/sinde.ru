-- +goose Up
CREATE TABLE IF NOT EXISTS storage_object_usages (
  usage_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  object_id UUID NOT NULL REFERENCES storage_objects(object_id) ON DELETE CASCADE,
  entity_type TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  usage_type TEXT NOT NULL,
  field_name TEXT NOT NULL DEFAULT '',
  sort_order INTEGER NOT NULL DEFAULT 0,
  is_primary BOOLEAN NOT NULL DEFAULT false,
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_storage_object_usages_entity_type_not_blank
    CHECK (btrim(entity_type) <> ''),
  CONSTRAINT chk_storage_object_usages_entity_id_not_blank
    CHECK (btrim(entity_id) <> ''),
  CONSTRAINT chk_storage_object_usages_usage_type_not_blank
    CHECK (btrim(usage_type) <> '')
);
CREATE INDEX IF NOT EXISTS idx_storage_object_usages_entity_lookup
  ON storage_object_usages (
    entity_type,
    entity_id,
    usage_type,
    field_name,
    sort_order,
    created_at DESC
  );
CREATE INDEX IF NOT EXISTS idx_storage_object_usages_object_id
  ON storage_object_usages (object_id);
CREATE INDEX IF NOT EXISTS idx_storage_object_usages_primary_lookup
  ON storage_object_usages (entity_type, entity_id, is_primary, sort_order);
CREATE UNIQUE INDEX IF NOT EXISTS idx_storage_object_usages_unique_binding
  ON storage_object_usages (
    object_id,
    entity_type,
    entity_id,
    usage_type,
    field_name
  );
-- +goose Down
DROP INDEX IF EXISTS idx_storage_object_usages_unique_binding;
DROP INDEX IF EXISTS idx_storage_object_usages_primary_lookup;
DROP INDEX IF EXISTS idx_storage_object_usages_object_id;
DROP INDEX IF EXISTS idx_storage_object_usages_entity_lookup;
DROP TABLE IF EXISTS storage_object_usages;