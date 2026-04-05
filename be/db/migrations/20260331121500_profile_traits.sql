-- +goose Up
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS primary_trait_uuid UUID NULL;

CREATE TABLE IF NOT EXISTS auth_saved_trait_sets (
  saved_set_id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  trait_uuid UUID NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_auth_saved_trait_sets_user_trait_uuid UNIQUE (user_id, trait_uuid),
  CONSTRAINT chk_auth_saved_trait_sets_name CHECK (char_length(btrim(name)) BETWEEN 1 AND 120),
  CONSTRAINT chk_auth_saved_trait_sets_description CHECK (char_length(description) <= 280)
);

CREATE INDEX IF NOT EXISTS idx_auth_saved_trait_sets_user_id ON auth_saved_trait_sets (user_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_users_primary_trait_uuid ON users (primary_trait_uuid);

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_auth_saved_trait_sets_updated_at ON auth_saved_trait_sets;
CREATE TRIGGER trg_auth_saved_trait_sets_updated_at
  BEFORE UPDATE ON auth_saved_trait_sets
  FOR EACH ROW
  EXECUTE FUNCTION set_auth_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_auth_saved_trait_sets_updated_at ON auth_saved_trait_sets;
-- +goose StatementEnd

DROP TABLE IF EXISTS auth_saved_trait_sets;
DROP INDEX IF EXISTS idx_users_primary_trait_uuid;
ALTER TABLE users
  DROP COLUMN IF EXISTS primary_trait_uuid;
