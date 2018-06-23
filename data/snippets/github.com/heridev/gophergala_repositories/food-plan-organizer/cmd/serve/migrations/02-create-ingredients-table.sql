-- +migrate Up

CREATE TABLE recipes (
  id             integer PRIMARY KEY ASC,
  name           TEXT,
  description    TEXT
);

-- +migrate Down
DROP TABLE recipes;
