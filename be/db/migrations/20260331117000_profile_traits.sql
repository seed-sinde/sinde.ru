-- +goose Up
ALTER TABLE users
  ADD COLUMN primary_trait_uuid uuid NULL;
CREATE TABLE auth_saved_trait_sets(
  saved_set_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  trait_uuid uuid NOT NULL,
  name text NOT NULL,
  description text NOT NULL DEFAULT '',
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT uq_auth_saved_trait_sets_user_trait_uuid UNIQUE (user_id, trait_uuid),
  CONSTRAINT chk_auth_saved_trait_sets_name CHECK (char_length(btrim(name)) BETWEEN 1 AND 120),
  CONSTRAINT chk_auth_saved_trait_sets_description CHECK (char_length(description) <= 280)
);
CREATE INDEX idx_auth_saved_trait_sets_user_id
  ON auth_saved_trait_sets(user_id, updated_at DESC);
CREATE INDEX idx_users_primary_trait_uuid
  ON users(primary_trait_uuid);
-- +goose StatementBegin
CREATE TRIGGER trg_auth_saved_trait_sets_updated_at
  BEFORE UPDATE ON auth_saved_trait_sets
  FOR EACH ROW
  EXECUTE FUNCTION set_auth_updated_at();
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER trg_auth_saved_trait_sets_updated_at ON auth_saved_trait_sets;
-- +goose StatementEnd
DROP INDEX idx_auth_saved_trait_sets_user_id;
DROP INDEX idx_users_primary_trait_uuid;
DROP TABLE auth_saved_trait_sets;
ALTER TABLE users
  DROP COLUMN primary_trait_uuid;