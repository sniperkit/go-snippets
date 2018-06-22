// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package features

import (
	"database/sql"
	"errors"

	"github.com/DevMine/featscomp/util"
)

type followersCount struct {
	name string
	db   *sql.DB
}

var _ Feature = (*followersCount)(nil)

// NewFollowersCount creates a new followers count feature.
func NewFollowersCount(name string, db *sql.DB) (*followersCount, error) {
	if name == "" {
		return nil, errors.New("nbFollwers: name cannot be empty")
	}

	if db == nil {
		return nil, errors.New("followersCount: db cannot be nil")
	}

	return &followersCount{name: name, db: db}, nil
}

// Score computes the scores of the followers count feature.
func (n followersCount) Score() error {
	if err := deleteOldScores(n.db, n.name); err != nil {
		return err
	}

	rows, err := n.db.Query(
		`SELECT u.id, ghu.followers_count
         FROM users AS u
         LEFT JOIN gh_users AS ghu ON u.id=ghu.user_id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	featID, err := util.FindFeatureID(n.db, n.name)
	if err != nil {
		return err
	}

	var followers []float64
	var userIDs []int64
	var maxVal float64
	for rows.Next() {
		var userID int64
		var followersCount *int64
		if err := rows.Scan(&userID, &followersCount); err != nil {
			return err
		}

		var followersCountF float64
		if followersCount != nil {
			followersCountF = float64(*followersCount)
		}

		if followersCountF > maxVal {
			maxVal = followersCountF
		}

		followers = append(followers, followersCountF)
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	coeff := 1.0 / maxVal
	for i, f := range followers {
		if _, err = n.db.Exec(
			`INSERT INTO scores(user_id, feature_id, score)
			 VALUES($1, $2, $3)`,
			userIDs[i], featID, f*coeff); err != nil {
			return err
		}
	}

	return nil
}

// Name returns the name of the feature.
func (n followersCount) Name() string {
	return n.name
}
