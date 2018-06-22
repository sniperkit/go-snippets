-- insert repos into repositories and gh_repositories from tmp_gh_repositories table
CREATE OR REPLACE FUNCTION insert_repos() RETURNS void AS
$BODY$
DECLARE
    repo_id repositories.id%TYPE;
    repo tmp_gh_repositories%ROWTYPE;
BEGIN
    -- disable constraints
    ALTER TABLE ONLY repositories DROP CONSTRAINT repositories_unique_clone_path;
    ALTER TABLE ONLY repositories DROP CONSTRAINT repositories_unique_clone_url;
    ALTER TABLE ONLY gh_repositories DROP CONSTRAINT gh_repositories_fk_repositories;

    FOR repo IN
        -- get all non already inserted repositories, without duplicates
        SELECT DISTINCT
            tgr.name,
            tgr.primary_language,
            tgr.clone_url,
            tgr.clone_path,
            tgr.vcs,
            tgr.github_id,
            tgr.full_name,
            tgr.description,
            tgr.homepage,
            tgr.fork,
            tgr.default_branch,
            tgr.master_branch,
            tgr.html_url,
            tgr.forks_count,
            tgr.open_issues_count,
            tgr.stargazers_count,
            tgr.subscribers_count,
            tgr.watchers_count,
            tgr.size_in_kb,
            tgr.created_at,
            tgr.updated_at,
            tgr.pushed_at
        FROM tmp_gh_repositories AS tgr
        INNER JOIN (
            SELECT
                clone_path,
                max(updated_at) AS max_updated_at,
                max(pushed_at) AS max_pushed_at,
                min(open_issues_count) AS max_open_issues_count
            FROM tmp_gh_repositories
            GROUP BY clone_path) AS tmp ON (
                tmp.clone_path = tgr.clone_path AND
                tmp.max_updated_at = tgr.updated_at AND
                tmp.max_pushed_at = tgr.pushed_at AND
                tmp.max_open_issues_count = tgr.open_issues_count
            )
        LEFT JOIN gh_repositories AS gr ON tgr.github_id = gr.github_id
        LEFT JOIN repositories AS r ON (tgr.clone_path = r.clone_path AND tgr.primary_language = r.primary_language)
        WHERE gr.id IS NULL AND r.id IS NULL AND tgr.clone_url <> '' AND tgr.clone_path <> '' AND tgr.primary_language <> ''
    LOOP
        -- raise notice 'Value: %', repo;

        -- create repositories
        INSERT INTO repositories (name, primary_language, clone_url, clone_path, vcs)
        VALUES (repo.name, repo.primary_language, repo.clone_url, repo.clone_path, repo.vcs)
        RETURNING id INTO repo_id;

        -- create gh_repositories
        INSERT INTO gh_repositories (repository_id, github_id, full_name, description, homepage, fork, default_branch, master_branch, html_url, forks_count, open_issues_count, stargazers_count, subscribers_count, watchers_count, size_in_kb, created_at, updated_at, pushed_at)
        VALUES(
            repo_id,
            repo.github_id,
            repo.full_name,
            repo.description,
            repo.homepage,
            repo.fork,
            repo.default_branch,
            repo.master_branch,
            repo.html_url,
            repo.forks_count,
            repo.open_issues_count,
            repo.stargazers_count,
            repo.subscribers_count,
            repo.watchers_count,
            repo.size_in_kb,
            repo.created_at,
            repo.updated_at,
            repo.pushed_at
        );
    END LOOP;

    -- re-enable constraints
    ALTER TABLE ONLY repositories ADD CONSTRAINT repositories_unique_clone_path UNIQUE (clone_path);
    ALTER TABLE ONLY repositories ADD CONSTRAINT repositories_unique_clone_url UNIQUE (clone_url);
    ALTER TABLE ONLY gh_repositories ADD CONSTRAINT gh_repositories_fk_repositories FOREIGN KEY (repository_id) REFERENCES repositories(id);
END
$BODY$
LANGUAGE plpgsql;

SELECT insert_repos();
