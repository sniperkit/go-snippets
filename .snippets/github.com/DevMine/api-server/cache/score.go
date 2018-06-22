// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cache

import (
	"database/sql"
	"strings"

	mx "code.google.com/p/biogo.matrix"

	"github.com/DevMine/api-server/model"
	"github.com/DevMine/api-server/util/typeutil"
)

// loadFeatures loads all features into memory.
func loadFeatures(db *sql.DB) error {
	rows, err := db.Query(
		`SELECT f.id, f.name, f.category, f.default_weight
         FROM features AS f
         ORDER BY f.name ASC`)
	if err != nil {
		return err
	}

	var feats []model.Feature
	for rows.Next() {
		var f model.Feature
		if err := rows.Scan(&f.ID, &f.Name, &f.Category, &f.DefaultWeight); err != nil {
			return err
		}
		feats = append(feats, f)
	}

	features = feats

	return nil
}

// loadScoresAndUsers loads the scores matrix and the users vector into memory.
func loadScoresAndUsers(db *sql.DB) error {

	var nbUsers uint
	if err := db.QueryRow(`SELECT COUNT(users.id) FROM users`).Scan(&nbUsers); err != nil {
		return err
	}

	var nbFeats uint
	if err := db.QueryRow(`SELECT COUNT(features.id) FROM features`).Scan(&nbFeats); err != nil {
		return err
	}

	rows, err := db.Query(
		`SELECT users.id, users.username, users.name, users.email,
				array_to_string(array_agg(s.score ORDER BY s.feature_id), ' ')
			FROM scores AS s
            JOIN users ON s.user_id=users.id
			GROUP BY users.id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	m := make([][]float64, nbUsers)
	users := make([]model.User, nbUsers)

	for i := 0; rows.Next(); i++ {
		var u model.User
		var scores string
		m[i] = make([]float64, nbFeats)

		if err := rows.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &scores); err != nil {
			return err
		}

		users[i] = u
		for j, s := range strings.Split(scores, " ") {
			if m[i][j], err = typeutil.StrToFloat(s); err != nil {
				return err
			}
		}
	}

	scores, err := mx.NewSparse(m)
	if err != nil {
		return err
	}

	scoresMatrix = scores
	usersVector = users

	return nil
}
