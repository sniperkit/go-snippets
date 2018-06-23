-- +migrate Up

CREATE TABLE foods (
  id                   integer,
  food_group_id        integer,
  name                 text,
  short_name           text,
  common_name          text,
  scientific_name      text,
  nitrogen_factor      real,
  protein_factor       real,
  fat_factor           real,
  carbonhydrate_factor real
);

CREATE INDEX foods_idx ON foods (id);
CREATE INDEX foods_food_group_idx ON foods (food_group_id);
CREATE INDEX foods_name ON foods (name);
CREATE INDEX foods_short_name ON foods (short_name);
CREATE INDEX foods_common_name ON foods (common_name);
CREATE INDEX foods_scientific_name ON foods (scientific_name);

-- +migrate Down
DROP TABLE foods;
