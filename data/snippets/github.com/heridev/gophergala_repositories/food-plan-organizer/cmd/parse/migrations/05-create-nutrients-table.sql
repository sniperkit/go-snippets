-- +migrate Up

CREATE TABLE nutrients (
  food_id            integer,
  nutrient_id        integer,
  nutrient_value     REAL,
  min                REAL,
  max                REAL,
  degrees_of_freedom INTEGER,
  lower_error_bound  REAL,
  upper_error_bound  REAL
);
CREATE INDEX nutrients_food_idx ON nutrients (food_id);
CREATE INDEX nutrients_nutrient_idx ON nutrients (nutrient_id);

-- +migrate Down
DROP TABLE nutrients;
