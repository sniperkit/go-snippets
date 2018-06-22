-- +migrate Up

CREATE TABLE languages (
  nutrient_id text,
  factor_code text
);
CREATE INDEX languages_nutrient_idx ON languages (nutrient_id);
CREATE INDEX languages_factor_code_idx ON languages (factor_code);

-- +migrate Down
DROP TABLE languages;
