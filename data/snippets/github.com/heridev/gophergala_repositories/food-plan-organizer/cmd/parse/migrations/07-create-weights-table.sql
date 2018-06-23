-- +migrate Up

CREATE TABLE weights (
  nutrient_id integer,
  seq         text,
  amount      real,
  description text,
  gram_weight real
);
CREATE INDEX weights_idx ON weights (nutrient_id);

-- +migrate Down
DROP TABLE nutrient_definitions;
