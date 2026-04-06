-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE TYPE mineral_crystal_system AS ENUM(
  'cubic',
  'hexagonal',
  'monoclinic',
  'orthorhombic',
  'tetragonal',
  'triclinic',
  'unknown'
);
CREATE TABLE IF NOT EXISTS chemistry_minerals(
  id bigserial PRIMARY KEY,
  database_id bigint NOT NULL,
  mineral_name text NOT NULL,
  mineral_name_plain text NOT NULL,
  mineral_name_search text NOT NULL,
  valence_chemistry text,
  ima_chemistry text,
  chemistry_elements text[] NOT NULL DEFAULT '{}',
  ima_number text,
  rruff_ids text[] NOT NULL DEFAULT '{}',
  country_of_type_locality text,
  year_first_published integer,
  ima_status text,
  structural_groupname text,
  status_notes text,
  crystal_systems mineral_crystal_system[] NOT NULL DEFAULT '{}',
  space_groups text[] NOT NULL DEFAULT '{}',
  oldest_known_age_ma numeric(12, 4),
  valence_elements text[] NOT NULL DEFAULT '{}',
  ima_mineral_symbol text,
  paragenetic_modes text[] NOT NULL DEFAULT '{}',
  images jsonb NOT NULL DEFAULT '[]'::jsonb,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT chemistry_minerals_database_id_not_blank CHECK (database_id > 0),
  CONSTRAINT chemistry_minerals_images_array CHECK (jsonb_typeof(images) = 'array')
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_chemistry_minerals_database_id_unique ON chemistry_minerals(database_id);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_name_plain ON chemistry_minerals(mineral_name_plain);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_name_search_trgm ON chemistry_minerals USING GIN(mineral_name_search gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_ima_status ON chemistry_minerals(ima_status);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_structural_groupname ON chemistry_minerals(structural_groupname);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_year_first_published ON chemistry_minerals(year_first_published);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_chemistry_elements_gin ON chemistry_minerals USING GIN(chemistry_elements);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_valence_elements_gin ON chemistry_minerals USING GIN(valence_elements);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_crystal_systems_gin ON chemistry_minerals USING GIN(crystal_systems);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_space_groups_gin ON chemistry_minerals USING GIN(space_groups);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_rruff_ids_gin ON chemistry_minerals USING GIN(rruff_ids);
CREATE INDEX IF NOT EXISTS idx_chemistry_minerals_paragenetic_modes_gin ON chemistry_minerals USING GIN(paragenetic_modes);
-- +goose Down
DROP TABLE IF EXISTS chemistry_minerals;
