// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package features

import (
	"database/sql"
	"errors"

	"github.com/DevMine/featscomp/util"
)

type commitsCount struct {
	name string
	db   *sql.DB
}

var _ Feature = (*commitsCount)(nil)

// NewCommitsCount creates a new commits count feature.
func NewCommitsCount(name string, db *sql.DB) (*commitsCount, error) {
	if name == "" {
		return nil, errors.New("commitsCount: name cannot be empty")
	}

	if db == nil {
		return nil, errors.New("commitsCount: db cannot be nil")
	}

	return &commitsCount{name: name, db: db}, nil
}

// Score computes the scores of the commits count feature.
func (cc commitsCount) Score() error {
	if err := deleteOldScores(cc.db, cc.name); err != nil {
		return err
	}

	rows, err := cc.db.Query(`
        SELECT u.id, COUNT(c.id)
        FROM users AS u
        LEFT JOIN commits AS c ON c.author_id = u.id
        GROUP BY c.author_id, u.id
        ORDER BY u.id ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	featID, err := util.FindFeatureID(cc.db, cc.name)
	if err != nil {
		return err
	}

	var commits []float64
	var userIDs []int64
	var ciCountMax float64
	for rows.Next() {
		var userID int64
		var ciCount *int64
		if err := rows.Scan(&userID, &ciCount); err != nil {
			return err
		}

		var ciCountF float64
		if ciCount != nil {
			ciCountF = float64(*ciCount)
		}

		if ciCountF > ciCountMax {
			ciCountMax = ciCountF
		}
		commits = append(commits, ciCountF)
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	coeff := 1.0 / ciCountMax
	for i, ci := range commits {
		if _, err = cc.db.Exec(
			`INSERT INTO scores(user_id, feature_id, score)
             VALUES($1, $2, $3)`,
			userIDs[i], featID, ci*coeff); err != nil {
			return err
		}
	}

	return nil
}

// Name returns the name of the feature.
func (cc commitsCount) Name() string {
	return cc.name
}
