-- +goose Up
CREATE TABLE IF NOT EXISTS storage_objects(
  object_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  storage_key text NOT NULL UNIQUE,
  bucket_name text NOT NULL,
  media_family text NOT NULL,
  content_type text NOT NULL,
  byte_size bigint NOT NULL CHECK (byte_size > 0),
  file_hash text NOT NULL DEFAULT '',
  source_kind text NOT NULL DEFAULT 'runtime',
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT chk_storage_objects_storage_key_not_blank CHECK (btrim(storage_key) <> ''),
  CONSTRAINT chk_storage_objects_bucket_name_not_blank CHECK (btrim(bucket_name) <> ''),
  CONSTRAINT chk_storage_objects_media_family_not_blank CHECK (btrim(media_family) <> ''),
  CONSTRAINT chk_storage_objects_content_type_not_blank CHECK (btrim(content_type) <> ''),
  CONSTRAINT chk_storage_objects_source_kind CHECK (source_kind IN ('seed', 'runtime', 'import'))
);
ALTER TABLE storage_objects
  ADD COLUMN IF NOT EXISTS file_hash text,
  ADD COLUMN IF NOT EXISTS source_kind text,
  ADD COLUMN IF NOT EXISTS metadata jsonb,
  ADD COLUMN IF NOT EXISTS updated_at timestamptz;

UPDATE storage_objects
SET
  file_hash = COALESCE(file_hash, ''),
  source_kind = COALESCE(source_kind, 'runtime'),
  metadata = COALESCE(metadata, '{}'::jsonb),
  updated_at = COALESCE(updated_at, created_at, now())
WHERE
  file_hash IS NULL
  OR source_kind IS NULL
  OR metadata IS NULL
  OR updated_at IS NULL;

ALTER TABLE storage_objects
  ALTER COLUMN file_hash SET DEFAULT '',
  ALTER COLUMN file_hash SET NOT NULL,
  ALTER COLUMN source_kind SET DEFAULT 'runtime',
  ALTER COLUMN source_kind SET NOT NULL,
  ALTER COLUMN metadata SET DEFAULT '{}'::jsonb,
  ALTER COLUMN metadata SET NOT NULL,
  ALTER COLUMN updated_at SET DEFAULT now(),
  ALTER COLUMN updated_at SET NOT NULL;

-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'chk_storage_objects_source_kind'
      AND conrelid = 'storage_objects'::regclass
  ) THEN
    ALTER TABLE storage_objects
      ADD CONSTRAINT chk_storage_objects_source_kind
      CHECK (source_kind IN ('seed', 'runtime', 'import'));
  END IF;
END
$$;
-- +goose StatementEnd
CREATE INDEX IF NOT EXISTS idx_storage_objects_bucket_created ON storage_objects(bucket_name, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_storage_objects_media_family_created ON storage_objects(media_family, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_storage_objects_source_kind_created ON storage_objects(source_kind, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_storage_objects_file_hash ON storage_objects(file_hash);
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_storage_objects_updated_at()
  RETURNS TRIGGER
  AS $fn$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$fn$
LANGUAGE plpgsql;
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
DROP INDEX IF EXISTS idx_storage_objects_file_hash;
DROP INDEX IF EXISTS idx_storage_objects_source_kind_created;
DROP INDEX IF EXISTS idx_storage_objects_media_family_created;
DROP INDEX IF EXISTS idx_storage_objects_bucket_created;
DROP TABLE IF EXISTS storage_objects;
