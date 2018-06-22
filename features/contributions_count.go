// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package features

import (
	"database/sql"
	"errors"

	"github.com/DevMine/featscomp/util"
)

type contributionsCount struct {
	name string
	db   *sql.DB
}

var _ Feature = (*contributionsCount)(nil)

// NewContributionsCount creates a new contributions count feature.
func NewContributionsCount(name string, db *sql.DB) (*contributionsCount, error) {
	if name == "" {
		return nil, errors.New("contributionsCount: name cannot be empty")
	}

	if db == nil {
		return nil, errors.New("contributionsCount: db cannot be nil")
	}

	return &contributionsCount{name: name, db: db}, nil
}

// Score computes the scores of the contributions count feature.
func (n contributionsCount) Score() error {
	if err := deleteOldScores(n.db, n.name); err != nil {
		return err
	}

	rows, err := n.db.Query(
		`SELECT COUNT(r.id), u.id
         FROM users AS u
         LEFT JOIN users_repositories AS ur ON u.id = ur.user_id
         LEFT JOIN repositories AS r ON ur.repository_id = r.id
         LEFT JOIN gh_repositories AS ghr ON r.id = ghr.repository_id
		 WHERE NOT ghr.fork
		 GROUP BY u.id
		 ORDER BY u.id ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	featID, err := util.FindFeatureID(n.db, n.name)
	if err != nil {
		return err
	}

	var contribs []float64
	var userIDs []int64
	var contribCountMax float64
	for rows.Next() {
		var userID int64
		var contCount *int64
		if err := rows.Scan(&contCount, &userID); err != nil {
			return err
		}
		var contCountF float64
		if contCount != nil {
			contCountF = float64(*contCount)
		}
		if contCountF > contribCountMax {
			contribCountMax = contCountF
		}
		contribs = append(contribs, contCountF)
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	coeff := 1.0 / contribCountMax
	for i, c := range contribs {
		if _, err = n.db.Exec(
			`INSERT INTO scores(user_id, feature_id, score)
             VALUES($1, $2, $3)`,
			userIDs[i], featID, c*coeff); err != nil {
			return err
		}
	}

	return nil
}

// Name returns the name of the feature.
func (n contributionsCount) Name() string {
	return n.name
}
