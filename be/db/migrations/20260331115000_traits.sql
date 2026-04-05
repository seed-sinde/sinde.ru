-- +goose Up

-- Traits domain: immutable append-only trait graph.
-- Rules:
-- - s_childs length = 2
-- - allowed ordered forms: [t_uuid, t_uuid] or [s_uuid, t_uuid]
-- - all referenced UUIDs must exist
-- - s_childs array itself must not be NULL
-- - s_childs elements must not be NULL
-- - s_childs elements must be distinct
-- - UPDATE sets is forbidden
-- - DELETE sets is forbidden

CREATE TABLE IF NOT EXISTS key_syns (
  id   BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS traits_k (
  id     BIGSERIAL PRIMARY KEY,
  syn_id BIGINT NOT NULL REFERENCES key_syns(id) ON DELETE RESTRICT,
  meta   JSONB  NOT NULL DEFAULT '{}'::jsonb,
  CONSTRAINT uq_traits_k_syn_id_meta UNIQUE (syn_id, meta)
);

CREATE TABLE IF NOT EXISTS traits_v (
  t_uuid  UUID PRIMARY KEY,
  t_key   BIGINT NOT NULL REFERENCES traits_k(id) ON DELETE RESTRICT,
  t_value TEXT   NOT NULL,
  CONSTRAINT uq_traits_v_t_key_t_value UNIQUE (t_key, t_value)
);

CREATE TABLE IF NOT EXISTS sets (
  s_uuid   UUID PRIMARY KEY,
  s_childs UUID[] NOT NULL,
  CONSTRAINT chk_sets_s_childs_length CHECK (cardinality(s_childs) = 2),
  CONSTRAINT chk_sets_s_childs_no_nulls CHECK (s_childs[1] IS NOT NULL AND s_childs[2] IS NOT NULL),
  CONSTRAINT chk_sets_s_childs_distinct CHECK (s_childs[1] <> s_childs[2]),
  CONSTRAINT uq_sets_s_childs UNIQUE (s_childs)
);

COMMENT ON TABLE sets IS 'Immutable append-only trait graph. INSERT allowed; UPDATE and DELETE forbidden.';

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION validate_sets_s_childs() RETURNS TRIGGER AS $fn$
BEGIN
  -- Allowed ordered forms only:
  --   [t_uuid, t_uuid]
  --   [s_uuid, t_uuid]

  IF EXISTS (
    SELECT 1
    FROM traits_v t1
    JOIN traits_v t2 ON TRUE
    WHERE t1.t_uuid = NEW.s_childs[1]
      AND t2.t_uuid = NEW.s_childs[2]
  ) THEN
    RETURN NEW;
  END IF;

  IF EXISTS (
    SELECT 1
    FROM sets s
    JOIN traits_v t ON TRUE
    WHERE s.s_uuid = NEW.s_childs[1]
      AND t.t_uuid = NEW.s_childs[2]
  ) THEN
    RETURN NEW;
  END IF;

  RAISE EXCEPTION 'invalid s_childs: allowed only [t_uuid, t_uuid] or [s_uuid, t_uuid] with existing referenced rows';
END;
$fn$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_validate_sets_s_childs ON sets;
CREATE TRIGGER trg_validate_sets_s_childs
  BEFORE INSERT ON sets
  FOR EACH ROW
  EXECUTE FUNCTION validate_sets_s_childs();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION forbid_sets_update() RETURNS TRIGGER AS $fn$
BEGIN
  RAISE EXCEPTION 'sets rows are immutable: UPDATE is not allowed';
END;
$fn$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_forbid_sets_update ON sets;
CREATE TRIGGER trg_forbid_sets_update
  BEFORE UPDATE ON sets
  FOR EACH ROW
  EXECUTE FUNCTION forbid_sets_update();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION forbid_sets_delete() RETURNS TRIGGER AS $fn$
BEGIN
  RAISE EXCEPTION 'sets rows are immutable: DELETE is not allowed';
END;
$fn$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_forbid_sets_delete ON sets;
CREATE TRIGGER trg_forbid_sets_delete
  BEFORE DELETE ON sets
  FOR EACH ROW
  EXECUTE FUNCTION forbid_sets_delete();
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_forbid_sets_delete ON sets;
DROP TRIGGER IF EXISTS trg_forbid_sets_update ON sets;
DROP TRIGGER IF EXISTS trg_validate_sets_s_childs ON sets;
-- +goose StatementEnd

DROP FUNCTION IF EXISTS forbid_sets_delete();
DROP FUNCTION IF EXISTS forbid_sets_update();
DROP FUNCTION IF EXISTS validate_sets_s_childs();

DROP TABLE IF EXISTS sets;
DROP TABLE IF EXISTS traits_v;
DROP TABLE IF EXISTS traits_k;
DROP TABLE IF EXISTS key_syns;
