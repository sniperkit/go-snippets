-- name: delete-repos
DELETE FROM repos;
-- name: insert-repos
INSERT INTO repos
  (id, created, updated, url, branch, latest_commit)
VALUES
  ('ebfae531-3e35-4d2b-9b91-3f71bf6df664', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'https://github.com/ipfs/distributed-wikipedia-mirror', 'master', '1dd318d678fd5496e5c93ad15aecb16dc128d36d');

-- name: delete-sources
DELETE FROM sources;
-- name: insert-sources
INSERT INTO sources
  (id, created, updated, title, url, checksum, meta)
VALUES
  ('bac6f89b-703e-4751-8109-b14d604df746', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'wikipedia ab full' , 'http://download.kiwix.org/zim/wikipedia_ab_all.zim', 'cefb808c89d55a1085966efc47df3d38', null),
  ('78f9b2a4-cf20-438e-adc8-77396d831327', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'wikipedia ace full' , 'http://download.kiwix.org/zim/wikipedia_ace_all.zim', 'eb43e54729670f7116ace499aee2f692', null);

-- name: delete-repo_sources
DELETE FROM repo_sources;
-- name: insert-repo_sources
INSERT INTO repo_sources
  (repo_id, source_id)
VALUES
  ('ebfae531-3e35-4d2b-9b91-3f71bf6df664', 'bac6f89b-703e-4751-8109-b14d604df746'),
  ('ebfae531-3e35-4d2b-9b91-3f71bf6df664', '78f9b2a4-cf20-438e-adc8-77396d831327');

-- name: delete-tasks
DELETE FROM tasks;
-- name: insert-tasks
INSERT INTO tasks
  (id, created, updated, title, user_id, type, params, status, error, enqueued, started, succeeded, failed)
  -- (id, created, updated, title, request, success, fail, repo_url, repo_commit, source_url, source_checksum, result_url, result_hash, message)
VALUES
  ('57220705-4954-4a42-9e02-e6aa53b6908e', '2017-01-01 00:00:01', '2017-01-01 00:00:01', 'Add a url to IPFS', '', 'ipfs.add', null, '', '', null, null, null,null);
