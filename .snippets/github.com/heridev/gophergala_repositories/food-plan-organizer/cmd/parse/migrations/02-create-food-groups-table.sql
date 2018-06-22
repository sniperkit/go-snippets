-- +migrate Up

CREATE TABLE food_groups (
  id   integer,
  name text
);
CREATE INDEX food_groups_idx ON food_groups (id);

-- +migrate Down
DROP TABLE food_groups;
