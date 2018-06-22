# Fields for the tables generated from the mysql dump
table_fields = {"commits":
               [
                   "author_id",
                   "committer_id",
                   "project_id",
                   "created_at"
               ],
               "counters":
               [
                   "id",
                   "date",
                   "commit_comments",
                   "commit_parents",
                   "commits",
                   "followers",
                   "organization_members",
                   "projects",
                   "users",
                   "issues",
                   "pull_requests",
                   "issue_comments",
                   "pull_request_comments",
                   "pull_request_history",
                   "watchers",
                   "forks"
               ],
               "followers":
               [
                   "follower_id",
                   "user_id"
               ],
               "forks":
               [
                   "forked_project_id",
                   "forked_from_id",
                   "created_at"
               ],
               "issues":
               [
                   "repo_id",
                   "reporter_id",
                   "assignee_id",
                   "created_at"
               ],
               "organization_members":
               [
                   "org_id",
                   "user_id",
                   "created_at"
               ],
               "project_commits":
               [
                   "project_id",
                   "commit_id"
               ],
               "project_members":
               [
                   "repo_id",
                   "user_id",
                   "created_at"
               ],
               "projects":
               [
                   "id",
                   "owned_id",
                   "name",
                   "language",
                   "created_at",
                   "forked_from",
                   "deleted"
               ],
               "pull_requests":
               [
                   "id",
                   "head_repo_id",
                   "base_repo_id",
                   "head_commit_id",
                   "base_commid_id",
                   "user_id",
                   "merged"
               ],
               "users":
               [
                   "id",
                   "login",
                   "name",
                   "company",
                   "location",
                   "email",
                   "created_at",
                   "type"
               ],
               "watchers":
               [
                   "repo_id",
                   "user_id",
                   "created_at"
               ]}
