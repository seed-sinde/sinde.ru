-- +goose Up
CREATE TABLE IF NOT EXISTS storage_objects (
  storage_key TEXT PRIMARY KEY,
  storage_driver TEXT NOT NULL DEFAULT 'minio',
  bucket_name TEXT NOT NULL,
  file_hash TEXT NOT NULL,
  byte_size BIGINT NOT NULL CHECK (byte_size > 0),
  content_ext TEXT NOT NULL,
  content_type TEXT NOT NULL,
  media_family TEXT NOT NULL,
  source_path TEXT NOT NULL DEFAULT '',
  metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_storage_objects_driver CHECK (storage_driver IN ('minio'))
);

CREATE INDEX IF NOT EXISTS idx_storage_objects_bucket_created
  ON storage_objects (bucket_name, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_storage_objects_media_family_created
  ON storage_objects (media_family, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_storage_objects_file_hash
  ON storage_objects (file_hash);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_storage_objects_updated_at() RETURNS TRIGGER AS $fn$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$fn$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_storage_objects_updated_at ON storage_objects;
CREATE TRIGGER trg_storage_objects_updated_at
  BEFORE UPDATE ON storage_objects
  FOR EACH ROW
  EXECUTE FUNCTION set_storage_objects_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_storage_objects_updated_at ON storage_objects;
DROP FUNCTION IF EXISTS set_storage_objects_updated_at();
-- +goose StatementEnd

DROP TABLE IF EXISTS storage_objects;
