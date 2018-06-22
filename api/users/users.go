// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package users handles all /users... routes.
package users

import (
	"database/sql"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/DevMine/api-server/model"
	"github.com/DevMine/api-server/srv/context"
	"github.com/DevMine/api-server/util/apiutil"
	"github.com/DevMine/api-server/util/json"
)

const selectUsers = `
SELECT
    u.id, u.username, u.name, u.email,
    ghu.id, ghu.github_id, ghu.login, ghu.bio, ghu.blog, ghu.company,
    ghu.email, ghu.hireable, ghu.location, ghu.avatar_url, ghu.html_url,
    ghu.followers_count, ghu.following_count, ghu.collaborators_count,
    ghu.created_at, ghu.updated_at,
    array_agg(DISTINCT row(gho.id, gho.github_id, gho.login, gho.avatar_url,
        gho.html_url, gho.name, gho.company, gho.blog, gho.location, gho.email,
        gho.collaborators_count, gho.created_at, gho.updated_at)) AS gh_orgs
FROM users AS u
LEFT OUTER JOIN gh_users AS ghu ON u.id = ghu.user_id
JOIN gh_users_organizations AS ghuo ON ghu.id = ghuo.gh_user_id
LEFT OUTER JOIN gh_organizations AS gho ON ghuo.gh_organization_id = gho.id `

// Index handles "/users" route.
func Index(c *context.Context, w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query(selectUsers+
		`WHERE u.id >= $1
         GROUP BY ghu.id, u.id
         ORDER BY u.id ASC
         LIMIT $2`,
		c.SinceID,
		c.PerPage)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var u model.User
		var ghu model.GhUser
		var ghOrgsArray string

		if err := rows.Scan(
			&u.ID, &u.Username, &u.Name, &u.Email,
			&ghu.ID, &ghu.GithubID, &ghu.Login,
			&ghu.Bio, &ghu.Blog, &ghu.Company, &ghu.Email,
			&ghu.Hireable, &ghu.Location, &ghu.AvatarURL,
			&ghu.HTMLURL, &ghu.FollowersCount, &ghu.FollowingCount,
			&ghu.CollaboratorsCount, &ghu.CreatedAt, &ghu.UpdatedAt,
			&ghOrgsArray); err != nil {
			glog.Error(err)
			continue
		}
		ghu.GhOrganizations = apiutil.CreateGhOrgsFromPGArray(ghOrgsArray)
		u.GhUser = &ghu

		users = append(users, u)
	}

	w.Write(json.MarshalPanic(users))
}

// Show handles "/users/{username:[a-zA-Z0-9\\-_\\.]+}" route.
func Show(c *context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	var u model.User
	var ghu model.GhUser
	var ghOrgsArray string
	err := c.DB.QueryRow(selectUsers+
		`WHERE LOWER(u.username) = LOWER($1)
         GROUP BY ghu.id, u.id`,
		username).Scan(
		&u.ID, &u.Username, &u.Name, &u.Email,
		&ghu.ID, &ghu.GithubID, &ghu.Login,
		&ghu.Bio, &ghu.Blog, &ghu.Company, &ghu.Email,
		&ghu.Hireable, &ghu.Location, &ghu.AvatarURL,
		&ghu.HTMLURL, &ghu.FollowersCount, &ghu.FollowingCount,
		&ghu.CollaboratorsCount, &ghu.CreatedAt, &ghu.UpdatedAt,
		&ghOrgsArray)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			w.Write([]byte("{}"))
			return
		default:
			panic(err)
		}
	}
	ghu.GhOrganizations = apiutil.CreateGhOrgsFromPGArray(ghOrgsArray)
	u.GhUser = &ghu

	w.Write(json.MarshalIndentPanic(u))
}

// ShowCommits handles "/users/{username:[a-zA-Z0-9\\-_\\.]+}/commits" route.
func ShowCommits(c *context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	rows, err := c.DB.Query(`
        SELECT
            c.id, c.repository_id, c.author_id, c.committer_id,
            c.message, c.author_date, c.commit_date,
            c.file_changed_count, c.insertions_count, c.deletions_count
        FROM commits AS c
        INNER JOIN users AS u
        ON u.id=c.author_id
        WHERE c.id >= $1
        AND LOWER(u.username) = LOWER($2)
        ORDER BY c.id ASC
        LIMIT $3`,
		c.SinceID,
		username,
		c.PerPage)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	commits := make([]model.Commit, 0)

	for rows.Next() {
		var co model.Commit
		var authorID, committerID, repoID *int64

		if err := rows.Scan(
			&co.ID, &repoID, &authorID, &committerID,
			&co.Message, &co.AuthorDate, &co.CommitDate,
			&co.FileChangedCount, &co.InsertionsCount, &co.DeletionsCount); err != nil {
			glog.Error(err)
			continue
		}

		co.Author = apiutil.FetchUser(c.DB, authorID)
		co.Committer = apiutil.FetchUser(c.DB, committerID)
		co.Repository = apiutil.FetchRepository(c.DB, repoID)

		commits = append(commits, co)
	}

	w.Write(json.MarshalPanic(commits))
}

// ShowRepositories handles "/users/{username:[a-zA-Z0-9\\-_\\.]+}/repositories"
// route.
func ShowRepositories(c *context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	rows, err := c.DB.Query(`
		SELECT
			r.id, r.name, r.primary_language, r.clone_url, r.clone_path, r.vcs,
			ghr.id, ghr.github_id, ghr.full_name, ghr.description, ghr.homepage,
			ghr.fork, ghr.default_branch, ghr.master_branch, ghr.html_url,
			ghr.forks_count, ghr.open_issues_count, ghr.stargazers_count,
			ghr.subscribers_count, ghr.watchers_count, ghr.size_in_kb,
			ghr.created_at, ghr.updated_at, ghr.pushed_at
		FROM repositories AS r
		LEFT OUTER JOIN gh_repositories AS ghr
		ON ghr.repository_id = r.id
		INNER JOIN users_repositories AS ur
		ON r.id = ur.repository_id
		INNER JOIN users AS u
		ON ur.user_id = u.id
		WHERE LOWER(u.username) = LOWER($1)
		GROUP BY r.id, ghr.id
		LIMIT $2`,
		username,
		c.PerPage)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	repositories := make([]model.Repository, 0)

	for rows.Next() {
		var r model.Repository
		var ghr model.GhRepository

		if err := rows.Scan(
			&r.ID, &r.Name, &r.PrimaryLanguage, &r.CloneURL, &r.ClonePath,
			&r.VCS, &ghr.ID, &ghr.GithubID, &ghr.FullName, &ghr.Description,
			&ghr.Homepage, &ghr.Fork, &ghr.DefaultBranch, &ghr.MasterBranch,
			&ghr.HTMLURL, &ghr.ForksCount, &ghr.OpenIssuesCount, &ghr.StargazersCount,
			&ghr.SubscribersCount, &ghr.WatchersCount, &ghr.SizeInKb, &ghr.CreatedAt,
			&ghr.UpdatedAt, &ghr.PushedAt); err != nil {
			glog.Error(err)
			continue
		}
		if ghr.ID != nil {
			r.GhRepository = &ghr
		}

		repositories = append(repositories, r)
	}

	w.Write(json.MarshalPanic(repositories))
}

// ShowScores handles "/users/{username:[a-zA-Z0-9\\-_\\.]+}/scores" route.
func ShowScores(c *context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	rows, err := c.DB.Query(`
		SELECT f.name, s.score
		FROM scores AS s
		INNER JOIN features AS f ON s.feature_id = f.id
		INNER JOIN users AS u ON s.user_id = u.id
		WHERE LOWER(u.username) = LOWER($1)
		ORDER BY s.feature_id
		LIMIT $2`,
		username,
		c.PerPage)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	m := make(map[string]float64)

	for rows.Next() {
		var k string
		var v float64

		if err := rows.Scan(&k, &v); err != nil {
			glog.Error(err)
			continue
		}
		m[k] = v
	}
	w.Write(json.MarshalIndentPanic(m))
}
