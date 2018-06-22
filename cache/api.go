// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"database/sql"

	"github.com/DevMine/api-server/model"
)

// loadStats loads some database related statistics.
func loadStats(db *sql.DB) error {
	var err error
	var s model.Stats

	err = db.QueryRow(`SELECT COUNT(users.id) FROM users`).Scan(&s.UsersCount)
	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT COUNT(repositories.id) FROM repositories`).Scan(
		&s.RepositoriesCount)
	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT COUNT(commits.id) FROM commits`).Scan(
		&s.CommitsCount)
	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT COUNT(commit_diff_deltas.id) FROM commit_diff_deltas`).Scan(
		&s.CommitDeltasCount)
	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT COUNT(features.id) FROM features`).Scan(
		&s.FeaturesCount)
	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT COUNT(gh_users.id) FROM gh_users`).Scan(
		&s.GhUsersCount)
	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT COUNT(gh_organizations.id) FROM gh_organizations`).Scan(
		&s.GhOrganizationsCount)
	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT COUNT(gh_repositories.id) FROM gh_repositories`).Scan(
		&s.GhRepositoriesCount)
	if err != nil {
		return err
	}

	stats = &s

	return nil
}

// loadFeaturesNames loads the map of features names.
func loadFeaturesNames(db *sql.DB) error {
	feats := make(map[string]struct{})

	rows, err := db.Query(`SELECT features.name FROM features ORDER BY features.name ASC`)
	if err != nil {
		return err
	}

	for rows.Next() {
		var f string
		if err := rows.Scan(&f); err != nil {
			return err
		}

		feats[f] = struct{}{}
	}

	featuresNames = feats

	return nil
}
