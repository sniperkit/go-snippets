// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// GhRepository represents a GitHub repository.
type GhRepository struct {
	ID               *int64     `json:"id"`
	RepositoryID     *int64     `json:"-"`
	GithubID         *int64     `json:"github_id"`
	FullName         *string    `json:"full_name"`
	Description      *string    `json:"description"`
	Homepage         *string    `json:"homepage"`
	Fork             *bool      `json:"fork"`
	DefaultBranch    *string    `json:"default_branch"`
	MasterBranch     *string    `json:"master_branch"`
	HTMLURL          *string    `json:"html_url"`
	ForksCount       *int64     `json:"forks_count"`
	OpenIssuesCount  *int64     `json:"open_issues_count"`
	StargazersCount  *int64     `json:"stargazers_count"`
	SubscribersCount *int64     `json:"subscribers_count"`
	WatchersCount    *int64     `json:"watchers_count"`
	SizeInKb         *int64     `json:"size_in_kb"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
	PushedAt         *time.Time `json:"pushed_at"`
}
