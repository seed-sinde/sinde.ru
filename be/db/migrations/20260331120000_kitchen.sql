-- +goose Up
-- Kitchen domain init.
-- Dependency: users(user_id) must already exist.
-- Seed data is imported separately from CSV files in be/db/seeds/kitchen/.
CREATE TABLE kitchen_recipes(
  id uuid PRIMARY KEY,
  title text NOT NULL,
  description text NOT NULL DEFAULT '',
  prep_minutes int NOT NULL DEFAULT 0 CHECK (prep_minutes >= 0),
  cook_minutes int NOT NULL DEFAULT 0 CHECK (cook_minutes >= 0),
  servings int NOT NULL DEFAULT 1 CHECK (servings > 0),
  difficulty text NOT NULL DEFAULT 'easy',
  meal_type text NOT NULL DEFAULT 'other',
  cuisine text NOT NULL DEFAULT '',
  diet_type text NOT NULL DEFAULT '',
  cooking_method text NOT NULL DEFAULT '',
  ingredients jsonb NOT NULL,
  ingredients_search text[] NOT NULL DEFAULT '{}',
  steps jsonb NOT NULL,
  tags text[] NOT NULL DEFAULT '{}',
  is_public boolean NOT NULL DEFAULT TRUE,
  cover_image_key text NOT NULL DEFAULT '',
  owner_user_id uuid NULL REFERENCES users(user_id) ON DELETE SET NULL,
  moderation_status text NOT NULL DEFAULT 'approved',
  moderated_by_user_id uuid NULL REFERENCES users(user_id) ON DELETE SET NULL,
  moderated_at timestamptz NULL,
  moderation_note text NOT NULL DEFAULT '',
  kcal int NOT NULL DEFAULT 0 CHECK (kcal >= 0),
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT chk_kitchen_ingredients_array CHECK (jsonb_typeof(ingredients) = 'array'),
  CONSTRAINT chk_kitchen_steps_array CHECK (jsonb_typeof(steps) = 'array'),
  CONSTRAINT chk_kitchen_recipes_moderation_status CHECK (moderation_status IN ('draft', 'pending', 'approved', 'rejected'))
);
CREATE INDEX idx_kitchen_recipes_created_at
  ON kitchen_recipes(created_at DESC);
CREATE INDEX idx_kitchen_recipes_title_lower
  ON kitchen_recipes(lower(title));
CREATE INDEX idx_kitchen_recipes_ingredients_search_gin
  ON kitchen_recipes USING GIN(ingredients_search);
CREATE INDEX idx_kitchen_recipes_diet_type
  ON kitchen_recipes(diet_type);
CREATE INDEX idx_kitchen_recipes_cooking_method
  ON kitchen_recipes(cooking_method);
CREATE INDEX idx_kitchen_recipes_owner_user_id
  ON kitchen_recipes(owner_user_id, created_at DESC);
CREATE INDEX idx_kitchen_recipes_moderation_status_created
  ON kitchen_recipes(moderation_status, created_at DESC);
-- +goose StatementBegin
CREATE FUNCTION set_kitchen_recipes_updated_at()
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
CREATE TRIGGER trg_kitchen_recipes_updated_at
  BEFORE UPDATE ON kitchen_recipes
  FOR EACH ROW
  EXECUTE FUNCTION set_kitchen_recipes_updated_at();
-- +goose StatementEnd
CREATE TABLE kitchen_user_ingredients(
  ingredient_id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  name text NOT NULL,
  category text NOT NULL DEFAULT 'другое',
  normalized_name text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT uq_kitchen_user_ingredients_user_normalized UNIQUE (user_id, normalized_name)
);
CREATE INDEX idx_kitchen_user_ingredients_user_id
  ON kitchen_user_ingredients(user_id, created_at DESC);
-- +goose StatementBegin
CREATE FUNCTION set_kitchen_user_ingredients_updated_at()
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
CREATE TRIGGER trg_kitchen_user_ingredients_updated_at
  BEFORE UPDATE ON kitchen_user_ingredients
  FOR EACH ROW
  EXECUTE FUNCTION set_kitchen_user_ingredients_updated_at();
-- +goose StatementEnd
CREATE TABLE kitchen_categories(
  category_id bigserial PRIMARY KEY,
  slug text NOT NULL UNIQUE,
  label text NOT NULL UNIQUE,
  description text NOT NULL DEFAULT '',
  sort_order int NOT NULL DEFAULT 0
);
CREATE TABLE kitchen_filter_options(
  option_id bigserial PRIMARY KEY,
  kind text NOT NULL,
  code text NOT NULL,
  label text NOT NULL,
  description text NOT NULL DEFAULT '',
  sort_order int NOT NULL DEFAULT 0,
  CONSTRAINT uq_kitchen_filter_options_kind_code UNIQUE (kind, code)
);
CREATE TABLE kitchen_ingredients_catalog(
  ingredient_id bigserial PRIMARY KEY,
  name text NOT NULL UNIQUE,
  normalized_name text NOT NULL UNIQUE,
  category_id bigint NOT NULL REFERENCES kitchen_categories(category_id) ON DELETE RESTRICT,
  description text NOT NULL DEFAULT '',
  protein_g numeric(8, 2) NULL,
  fat_g numeric(8, 2) NULL,
  carbs_g numeric(8, 2) NULL,
  kcal numeric(8, 2) NULL,
  sort_order int NOT NULL DEFAULT 0
);
CREATE INDEX idx_kitchen_ingredients_catalog_category_id
  ON kitchen_ingredients_catalog(category_id, sort_order, ingredient_id);
CREATE INDEX idx_kitchen_filter_options_kind
  ON kitchen_filter_options(kind, sort_order, option_id);
CREATE TABLE kitchen_favorite_ingredients(
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  ingredient_id bigint NOT NULL REFERENCES kitchen_ingredients_catalog(ingredient_id) ON DELETE CASCADE,
  name text NOT NULL,
  category text NOT NULL DEFAULT 'другое',
  normalized_name text NOT NULL,
  list_type text NOT NULL DEFAULT 'include',
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT pk_kitchen_favorite_ingredients PRIMARY KEY (user_id, ingredient_id, list_type),
  CONSTRAINT ck_kitchen_favorite_ingredients_list_type CHECK (list_type IN ('include', 'exclude'))
);
CREATE INDEX idx_kitchen_favorite_ingredients_user_id
  ON kitchen_favorite_ingredients(user_id, created_at DESC);
CREATE UNIQUE INDEX uq_kitchen_favorite_ingredients_user_normalized
  ON kitchen_favorite_ingredients(user_id, normalized_name, list_type);
CREATE TABLE kitchen_favorite_recipes(
  user_id uuid NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  recipe_id uuid NOT NULL REFERENCES kitchen_recipes(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT pk_kitchen_favorite_recipes PRIMARY KEY (user_id, recipe_id)
);
CREATE INDEX idx_kitchen_favorite_recipes_user_id
  ON kitchen_favorite_recipes(user_id, created_at DESC);
CREATE INDEX idx_kitchen_favorite_recipes_recipe_id
  ON kitchen_favorite_recipes(recipe_id, created_at DESC);
-- +goose Down
DROP INDEX idx_kitchen_favorite_recipes_recipe_id;
DROP INDEX idx_kitchen_favorite_recipes_user_id;
DROP TABLE kitchen_favorite_recipes;
DROP INDEX uq_kitchen_favorite_ingredients_user_normalized;
DROP INDEX idx_kitchen_favorite_ingredients_user_id;
DROP TABLE kitchen_favorite_ingredients;
DROP INDEX idx_kitchen_ingredients_catalog_category_id;
DROP TABLE kitchen_ingredients_catalog;
DROP INDEX idx_kitchen_filter_options_kind;
DROP TABLE kitchen_filter_options;
DROP TABLE kitchen_categories;
DROP TRIGGER trg_kitchen_user_ingredients_updated_at ON kitchen_user_ingredients;
DROP FUNCTION set_kitchen_user_ingredients_updated_at();
DROP INDEX idx_kitchen_user_ingredients_user_id;
DROP TABLE kitchen_user_ingredients;
DROP TRIGGER trg_kitchen_recipes_updated_at ON kitchen_recipes;
DROP FUNCTION set_kitchen_recipes_updated_at();
DROP INDEX idx_kitchen_recipes_moderation_status_created;
DROP INDEX idx_kitchen_recipes_owner_user_id;
DROP INDEX idx_kitchen_recipes_cooking_method;
DROP INDEX idx_kitchen_recipes_diet_type;
DROP INDEX idx_kitchen_recipes_ingredients_search_gin;
DROP INDEX idx_kitchen_recipes_title_lower;
DROP INDEX idx_kitchen_recipes_created_at;
DROP TABLE kitchen_recipes;