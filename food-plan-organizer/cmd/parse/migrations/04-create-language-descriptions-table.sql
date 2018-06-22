-- +migrate Up

CREATE TABLE language_descriptions (
  factor_code text,
  description text
);
CREATE INDEX language_descriptions_idx ON language_descriptions (factor_code);

-- +migrate Down
DROP TABLE language_descriptions;
