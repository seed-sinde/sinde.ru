-- +goose Up
CREATE TABLE chemistry_elements(
  number integer PRIMARY KEY,
  symbol text NOT NULL,
  name text NOT NULL,
  russian_name text NOT NULL,
  appearance text,
  atomic_mass numeric,
  boil numeric,
  category text NOT NULL,
  density numeric,
  discovered_by text,
  melt numeric,
  molar_heat numeric,
  named_by text,
  period integer,
  group_number integer,
  phase text,
  source text,
  bohr_model_image text,
  bohr_model_3d text,
  spectral_img text,
  samples jsonb NOT NULL DEFAULT '[]'::jsonb,
  sample_fallback jsonb,
  summary text NOT NULL DEFAULT '',
  xpos integer NOT NULL,
  ypos integer NOT NULL,
  wxpos integer,
  wypos integer,
  shells integer[] NOT NULL DEFAULT '{}',
  electron_configuration text,
  electron_configuration_semantic text,
  electron_affinity numeric,
  electronegativity_pauling numeric,
  ionization_energies numeric[] NOT NULL DEFAULT '{}',
  cpk_hex text,
  block text,
  CONSTRAINT chemistry_elements_number_positive CHECK (number > 0),
  CONSTRAINT chemistry_elements_symbol_not_blank CHECK (btrim(symbol) <> ''),
  CONSTRAINT chemistry_elements_name_not_blank CHECK (btrim(name) <> ''),
  CONSTRAINT chemistry_elements_russian_name_not_blank CHECK (btrim(russian_name) <> ''),
  CONSTRAINT chemistry_elements_category_not_blank CHECK (btrim(category) <> ''),
  CONSTRAINT chemistry_elements_summary_not_blank CHECK (btrim(summary) <> ''),
  CONSTRAINT chemistry_elements_samples_array CHECK (jsonb_typeof(samples) = 'array'),
  CONSTRAINT chemistry_elements_sample_fallback_object CHECK (sample_fallback IS NULL OR jsonb_typeof(sample_fallback) = 'object')
);
CREATE UNIQUE INDEX idx_chemistry_elements_symbol_unique ON chemistry_elements(symbol);
CREATE INDEX idx_chemistry_elements_category ON chemistry_elements(category);
CREATE INDEX idx_chemistry_elements_period ON chemistry_elements(period);
CREATE INDEX idx_chemistry_elements_group_number ON chemistry_elements(group_number);
CREATE INDEX idx_chemistry_elements_block ON chemistry_elements(block);
-- +goose Down
DROP INDEX idx_chemistry_elements_block;
DROP INDEX idx_chemistry_elements_group_number;
DROP INDEX idx_chemistry_elements_period;
DROP INDEX idx_chemistry_elements_category;
DROP INDEX idx_chemistry_elements_symbol_unique;
DROP TABLE chemistry_elements;
