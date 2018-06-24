package main

const qRepos = `
SELECT
  id, created, updated, url, branch, latest_commit
FROM repos
ORDER BY $1
LIMIT $2 OFFSET $3;`

const qRepoReadById = `
SELECT 
  id, created, updated, url, branch, latest_commit
FROM repos
WHERE id = $1;`

const qRepoInsert = `
INSERT INTO repos
  (id, created, updated, url, branch, latest_commit)
VALUES
  ($1, $2, $3, $4, $5, $6);`

const qRepoUpdate = `
UPDATE repos SET
  created = $2, updated = $3, url = $4, branch = $5, latest_commit = $6
WHERE id = $1;`

const qRepoDelete = `DELETE FROM repos WHERE id = $1;`
