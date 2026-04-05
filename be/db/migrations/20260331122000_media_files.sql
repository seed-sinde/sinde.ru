-- +goose Up
CREATE TABLE IF NOT EXISTS media_files (
  image_key TEXT PRIMARY KEY,
  file_hash TEXT NOT NULL,
  media_kind TEXT NOT NULL,
  owner_user_id UUID NULL REFERENCES users(user_id) ON DELETE SET NULL,
  recipe_id UUID NULL REFERENCES kitchen_recipes(id) ON DELETE SET NULL,
  content_ext TEXT NOT NULL,
  byte_size BIGINT NOT NULL CHECK (byte_size > 0),
  storage_section TEXT NOT NULL,
  storage_collection TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_media_files_kind CHECK (media_kind IN ('recipe', 'avatar'))
);

CREATE INDEX IF NOT EXISTS idx_media_files_owner_user_id ON media_files (owner_user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_media_files_recipe_id ON media_files (recipe_id, created_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS idx_media_files_hash_scope_unique
  ON media_files (file_hash, storage_section, storage_collection);
CREATE INDEX IF NOT EXISTS idx_media_files_section_collection
  ON media_files (storage_section, storage_collection, created_at DESC);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_media_files_updated_at() RETURNS TRIGGER AS $fn$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$fn$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_media_files_updated_at ON media_files;
CREATE TRIGGER trg_media_files_updated_at
  BEFORE UPDATE ON media_files
  FOR EACH ROW
  EXECUTE FUNCTION set_media_files_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_media_files_updated_at ON media_files;
DROP FUNCTION IF EXISTS set_media_files_updated_at();
-- +goose StatementEnd

DROP TABLE IF EXISTS media_files;
