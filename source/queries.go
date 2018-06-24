package source

const qSourceCreateTable = `
CREATE TABLE sources (
  id               UUID NOT NULL PRIMARY KEY,
  created          timestamp NOT NULL DEFAULT (now() at time zone 'utc'), 
  updated          timestamp NOT NULL DEFAULT (now() at time zone 'utc'), 
  title            text NOT NULL DEFAULT '',
  url              text NOT NULL,
  checksum         text NOT NULL DEFAULT '', 
  meta             json
);`

const qSourcesList = `
SELECT
  id, created, updated, title, url, checksum, meta
FROM sources
ORDER BY created DESC
LIMIT $1 OFFSET $2;`

const qSourceReadById = `
SELECT 
  id, created, updated, title, url, checksum, meta
FROM sources
WHERE 
  id = $1;`

const qSourceExists = `SELECT exists(SELECT 1 FROM sources WHERE id = $1);`

const qSourceInsert = `
INSERT INTO sources
  (id, created, updated, title, url, checksum, meta)
VALUES
  ($1, $2, $3, $4, $5, $6, $7);`

const qSourceUpdate = `
UPDATE sources SET
  created = $2, updated = $3, title = $4, url = $5, checksum = $6, meta = $7
WHERE id = $1;`

const qSourceDelete = `DELETE FROM sources WHERE id = $1;`
