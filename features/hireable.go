// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package features

import (
	"database/sql"
	"errors"

	"github.com/DevMine/featscomp/util"
)

type hireable struct {
	name string
	db   *sql.DB
}

var _ Feature = (*hireable)(nil)

// NewHireable creates a new hireable feature.
func NewHireable(name string, db *sql.DB) (*hireable, error) {
	if name == "" {
		return nil, errors.New("hireable: name cannot be empty")
	}
	if db == nil {
		return nil, errors.New("hireable: db cannot be nil")
	}

	return &hireable{name: name, db: db}, nil
}

// Score computes the scores of the hireable feature.
func (h hireable) Score() error {
	if err := deleteOldScores(h.db, h.name); err != nil {
		return err
	}

	rows, err := h.db.Query(
		`SELECT u.id, ghu.hireable
         FROM users AS u
         LEFT JOIN gh_users AS ghu ON u.id=ghu.user_id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	featID, err := util.FindFeatureID(h.db, h.name)
	if err != nil {
		return err
	}

	for rows.Next() {
		var userID int64
		var hireable *bool
		if err := rows.Scan(&userID, &hireable); err != nil {
			return err
		}

		var hire float64
		if hireable != nil && *hireable {
			hire = 1.0
		}

		if _, err = h.db.Exec(
			`INSERT INTO scores(user_id, feature_id, score)
             VALUES($1, $2, $3)`,
			userID, featID, hire); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// Name returns the name of the feature.
func (h hireable) Name() string {
	return h.name
}
