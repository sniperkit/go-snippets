package tasks

const qTaskCreateTable = `
CREATE TABLE tasks (
  id               UUID NOT NULL PRIMARY KEY,
  created          timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  updated          timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
  title            text NOT NULL DEFAULT '',
  user_id          text NOT NULL DEFAULT '',
  type             text NOT NULL DEFAULT '',
  params           json,
  status           text NOT NULL DEFAULT '',
  error            text NOT NULL DEFAULT '',
  enqueued         timestamp,
  started          timestamp,
  succeeded        timestamp,
  failed           timestamp
);`

// an available task a source.Checksum && repo.LatestCommit combination that doesn't
// have a task model already created.
// TODO - this is a carry-over from the former task_mgmt, need to rethink
const qAvailableTasks = `
WITH t AS (
  SELECT
    repos.url as repo_url,
    repos.latest_commit as repo_commit,
    sources.title as source_title,
    sources.url as source_url,
    sources.checksum as source_checksum
  FROM sources, repos, repo_sources
  WHERE 
    sources.id = repo_sources.source_id AND
    repos.id = repo_sources.repo_id
)
SELECT
  t.repo_url, t.repo_commit, t.source_title, t.source_url, t.source_checksum
FROM t LEFT OUTER JOIN tasks ON (t.source_url = tasks.source_url)
WHERE
  tasks.repo_commit is null OR
  tasks.source_checksum is null;`

const qTasks = `
SELECT
  id, created, updated, title, user_id, type,
  params, status, error, enqueued, started, succeeded, failed
FROM tasks
ORDER BY created DESC
LIMIT $1 OFFSET $2;`

const qTaskExists = `SELECT exists(SELECT 1 FROM tasks WHERE id = $1);`

const qTaskReadById = `
SELECT 
  id, created, updated, title, user_id, type,
  params, status, error, enqueued, started, succeeded, failed
FROM tasks
WHERE id = $1;`

const qTaskInsert = `
INSERT INTO tasks
  (id, created, updated, title, user_id, type,
   params, status, error, enqueued, started, succeeded, failed)
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`

const qTaskUpdate = `
UPDATE tasks SET
  created = $2, updated = $3, title = $4, user_id = $5, type = $6,
  params = $7, status = $8, error = $9, enqueued = $10, started = $11, succeeded = $12, failed = $13
WHERE id = $1;`

const qTaskDelete = `DELETE FROM tasks WHERE id = $1;`
