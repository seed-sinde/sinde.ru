-- +goose Up

-- Kitchen domain init.
-- Dependency: users(user_id) must already exist.
-- Seed data is imported separately from CSV files in be/db/seeds/kitchen/.

CREATE TABLE IF NOT EXISTS kitchen_recipes (
  id UUID PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  prep_minutes INT NOT NULL DEFAULT 0 CHECK (prep_minutes >= 0),
  cook_minutes INT NOT NULL DEFAULT 0 CHECK (cook_minutes >= 0),
  servings INT NOT NULL DEFAULT 1 CHECK (servings > 0),
  difficulty TEXT NOT NULL DEFAULT 'easy',
  meal_type TEXT NOT NULL DEFAULT 'other',
  cuisine TEXT NOT NULL DEFAULT '',
  diet_type TEXT NOT NULL DEFAULT '',
  cooking_method TEXT NOT NULL DEFAULT '',
  ingredients JSONB NOT NULL,
  ingredients_search TEXT[] NOT NULL DEFAULT '{}',
  steps JSONB NOT NULL,
  tags TEXT[] NOT NULL DEFAULT '{}',
  is_public BOOLEAN NOT NULL DEFAULT TRUE,
  cover_image_key TEXT NOT NULL DEFAULT '',
  owner_user_id UUID NULL REFERENCES users(user_id) ON DELETE SET NULL,
  moderation_status TEXT NOT NULL DEFAULT 'approved',
  moderated_by_user_id UUID NULL REFERENCES users(user_id) ON DELETE SET NULL,
  moderated_at TIMESTAMPTZ NULL,
  moderation_note TEXT NOT NULL DEFAULT '',
  kcal INT NOT NULL DEFAULT 0 CHECK (kcal >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_kitchen_ingredients_array CHECK (jsonb_typeof(ingredients) = 'array'),
  CONSTRAINT chk_kitchen_steps_array CHECK (jsonb_typeof(steps) = 'array'),
  CONSTRAINT chk_kitchen_recipes_moderation_status CHECK (moderation_status IN ('draft', 'pending', 'approved', 'rejected'))
);

CREATE INDEX IF NOT EXISTS idx_kitchen_recipes_created_at
  ON kitchen_recipes (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_kitchen_recipes_title_lower
  ON kitchen_recipes (lower(title));

CREATE INDEX IF NOT EXISTS idx_kitchen_recipes_ingredients_search_gin
  ON kitchen_recipes USING GIN (ingredients_search);

CREATE INDEX IF NOT EXISTS idx_kitchen_recipes_diet_type
  ON kitchen_recipes (diet_type);

CREATE INDEX IF NOT EXISTS idx_kitchen_recipes_cooking_method
  ON kitchen_recipes (cooking_method);

CREATE INDEX IF NOT EXISTS idx_kitchen_recipes_owner_user_id
  ON kitchen_recipes (owner_user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_kitchen_recipes_moderation_status_created
  ON kitchen_recipes (moderation_status, created_at DESC);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_kitchen_recipes_updated_at() RETURNS TRIGGER AS $fn$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$fn$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_kitchen_recipes_updated_at ON kitchen_recipes;
CREATE TRIGGER trg_kitchen_recipes_updated_at
  BEFORE UPDATE ON kitchen_recipes
  FOR EACH ROW
  EXECUTE FUNCTION set_kitchen_recipes_updated_at();
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS kitchen_user_ingredients (
  ingredient_id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT 'другое',
  normalized_name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_kitchen_user_ingredients_user_normalized UNIQUE (user_id, normalized_name)
);

CREATE INDEX IF NOT EXISTS idx_kitchen_user_ingredients_user_id
  ON kitchen_user_ingredients (user_id, created_at DESC);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_kitchen_user_ingredients_updated_at() RETURNS TRIGGER AS $fn$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$fn$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_kitchen_user_ingredients_updated_at ON kitchen_user_ingredients;
CREATE TRIGGER trg_kitchen_user_ingredients_updated_at
  BEFORE UPDATE ON kitchen_user_ingredients
  FOR EACH ROW
  EXECUTE FUNCTION set_kitchen_user_ingredients_updated_at();
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS kitchen_categories (
  category_id BIGSERIAL PRIMARY KEY,
  slug TEXT NOT NULL UNIQUE,
  label TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS kitchen_filter_options (
  option_id BIGSERIAL PRIMARY KEY,
  kind TEXT NOT NULL,
  code TEXT NOT NULL,
  label TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  sort_order INT NOT NULL DEFAULT 0,
  CONSTRAINT uq_kitchen_filter_options_kind_code UNIQUE (kind, code)
);

CREATE TABLE IF NOT EXISTS kitchen_ingredients_catalog (
  ingredient_id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  normalized_name TEXT NOT NULL UNIQUE,
  category_id BIGINT NOT NULL REFERENCES kitchen_categories(category_id) ON DELETE RESTRICT,
  description TEXT NOT NULL DEFAULT '',
  protein_g NUMERIC(8,2) NULL,
  fat_g NUMERIC(8,2) NULL,
  carbs_g NUMERIC(8,2) NULL,
  kcal NUMERIC(8,2) NULL,
  sort_order INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_kitchen_ingredients_catalog_category_id
  ON kitchen_ingredients_catalog (category_id, sort_order, ingredient_id);

CREATE INDEX IF NOT EXISTS idx_kitchen_filter_options_kind
  ON kitchen_filter_options (kind, sort_order, option_id);

CREATE TABLE IF NOT EXISTS kitchen_favorite_ingredients (
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  ingredient_id BIGINT NOT NULL REFERENCES kitchen_ingredients_catalog(ingredient_id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT 'другое',
  normalized_name TEXT NOT NULL,
  list_type TEXT NOT NULL DEFAULT 'include',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT pk_kitchen_favorite_ingredients PRIMARY KEY (user_id, ingredient_id, list_type),
  CONSTRAINT ck_kitchen_favorite_ingredients_list_type CHECK (list_type IN ('include', 'exclude'))
);

CREATE INDEX IF NOT EXISTS idx_kitchen_favorite_ingredients_user_id
  ON kitchen_favorite_ingredients (user_id, created_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS uq_kitchen_favorite_ingredients_user_normalized
  ON kitchen_favorite_ingredients (user_id, normalized_name, list_type);

CREATE TABLE IF NOT EXISTS kitchen_favorite_recipes (
  user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  recipe_id UUID NOT NULL REFERENCES kitchen_recipes(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT pk_kitchen_favorite_recipes PRIMARY KEY (user_id, recipe_id)
);

CREATE INDEX IF NOT EXISTS idx_kitchen_favorite_recipes_user_id
  ON kitchen_favorite_recipes (user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_kitchen_favorite_recipes_recipe_id
  ON kitchen_favorite_recipes (recipe_id, created_at DESC);

-- +goose Down

DROP TABLE IF EXISTS kitchen_favorite_recipes;
DROP TABLE IF EXISTS kitchen_favorite_ingredients;
DROP TABLE IF EXISTS kitchen_ingredients_catalog;
DROP TABLE IF EXISTS kitchen_filter_options;
DROP TABLE IF EXISTS kitchen_categories;
DROP TABLE IF EXISTS kitchen_user_ingredients;
DROP TABLE IF EXISTS kitchen_recipes;

DROP FUNCTION IF EXISTS set_kitchen_user_ingredients_updated_at();
DROP FUNCTION IF EXISTS set_kitchen_recipes_updated_at();
