-- +goose Up
CREATE TABLE IF NOT EXISTS kitchen_recipe_steps(
  step_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  recipe_id uuid NOT NULL REFERENCES kitchen_recipes(id) ON DELETE CASCADE,
  step_order integer NOT NULL CHECK (step_order >= 1),
  title text NOT NULL DEFAULT '',
  description text NOT NULL DEFAULT '',
  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT chk_kitchen_recipe_steps_title_not_null CHECK (title IS NOT NULL),
  CONSTRAINT chk_kitchen_recipe_steps_description_not_null CHECK (description IS NOT NULL)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_kitchen_recipe_steps_recipe_order_unique ON kitchen_recipe_steps(recipe_id, step_order);
CREATE INDEX IF NOT EXISTS idx_kitchen_recipe_steps_recipe_order ON kitchen_recipe_steps(recipe_id, step_order);
CREATE INDEX IF NOT EXISTS idx_kitchen_recipe_steps_created_at ON kitchen_recipe_steps(created_at DESC);
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_kitchen_recipe_steps_updated_at()
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
DROP TRIGGER IF EXISTS trg_kitchen_recipe_steps_updated_at ON kitchen_recipe_steps;
CREATE TRIGGER trg_kitchen_recipe_steps_updated_at
  BEFORE UPDATE ON kitchen_recipe_steps
  FOR EACH ROW
  EXECUTE FUNCTION set_kitchen_recipe_steps_updated_at();
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_kitchen_recipe_steps_updated_at ON kitchen_recipe_steps;
DROP FUNCTION IF EXISTS set_kitchen_recipe_steps_updated_at();
-- +goose StatementEnd
DROP INDEX IF EXISTS idx_kitchen_recipe_steps_created_at;
DROP INDEX IF EXISTS idx_kitchen_recipe_steps_recipe_order;
DROP INDEX IF EXISTS idx_kitchen_recipe_steps_recipe_order_unique;
DROP TABLE IF EXISTS kitchen_recipe_steps;
