// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// Stats represents various statistics about users, repositories and so on.
type Stats struct {
	UsersCount           *int64 `json:"users_count"`
	RepositoriesCount    *int64 `json:"repositories_count"`
	CommitsCount         *int64 `json:"commits_count"`
	CommitDeltasCount    *int64 `json:"commit_deltas_count"`
	FeaturesCount        *int64 `json:"features_count"`
	GhUsersCount         *int64 `json:"gh_users_count"`
	GhOrganizationsCount *int64 `json:"gh_organizations_count"`
	GhRepositoriesCount  *int64 `json:"gh_repositories_count"`
}
