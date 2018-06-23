-- +migrate Up

CREATE TABLE ingredients (
  id          integer primary key asc,
  recipe_id   integer,
  food_id     integer,
  unit        text,
  volume      real
);

-- +migrate Down
DROP TABLE ingredients;
