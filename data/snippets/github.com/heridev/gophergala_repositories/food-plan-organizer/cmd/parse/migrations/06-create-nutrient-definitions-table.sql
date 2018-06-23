-- +migrate Up

CREATE TABLE nutrient_definitions (
  nutrient_id    integer,
  units          TEXT,
  tagname        TEXT,
  description    TEXT,
  decimal_places integer
);
CREATE INDEX nutrient_definitions_idx ON nutrient_definitions (nutrient_id);

-- +migrate Down
DROP TABLE nutrient_definitions;
