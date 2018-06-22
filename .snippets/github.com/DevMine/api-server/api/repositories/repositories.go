// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package repositories handles /repositories... routes.
package repositories

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/DevMine/api-server/model"
	"github.com/DevMine/api-server/srv/context"
	"github.com/DevMine/api-server/util/json"
)

const selectRepositories = `
SELECT
	r.id, r.name, r.primary_language, r.clone_url, r.clone_path, r.vcs,
	ghr.id, ghr.github_id, ghr.full_name, ghr.description, ghr.homepage,
	ghr.fork, ghr.default_branch, ghr.master_branch, ghr.html_url,
	ghr.forks_count, ghr.open_issues_count, ghr.stargazers_count,
	ghr.subscribers_count, ghr.watchers_count, ghr.size_in_kb,
	ghr.created_at, ghr.updated_at, ghr.pushed_at
FROM repositories AS r
LEFT OUTER JOIN gh_repositories AS ghr
ON ghr.repository_id = r.id`

// Index handles "/repositories" route.
func Index(c *context.Context, w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query(selectRepositories+`
		WHERE r.id >= $1
		GROUP BY ghr.id, r.id
		ORDER BY r.id ASC
		LIMIT $2`,
		c.SinceID,
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

// Show handles "/repositories/{name:[a-zA-Z0-9\\-_\\.]+}" route.
func Show(c *context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	rows, err := c.DB.Query(selectRepositories+`
		WHERE LOWER(r.name) = LOWER($1)
		AND r.id >= $2
		GROUP BY ghr.id, r.id
		ORDER BY r.id ASC
		LIMIT $3`,
		name,
		c.SinceID,
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
