// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package features handles /features... routes.
package features

import (
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/DevMine/api-server/model"
	"github.com/DevMine/api-server/srv/context"
	"github.com/DevMine/api-server/util/json"
)

const selectFeatures = `
SELECT
    f.id, f.name, f.category, f.default_weight
FROM features AS f `

type user struct {
	ID       *int64   `json:"id"`
	Username *string  `json:"username"`
	Score    *float64 `json:"score"`
}

// Index handles "/features" route.
func Index(c *context.Context, w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query(selectFeatures+`
		WHERE f.id >= $1
        ORDER BY f.id ASC
        LIMIT $2`,
		c.SinceID,
		c.PerPage)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	features := make([]model.Feature, 0)

	for rows.Next() {
		var f model.Feature
		if err := rows.Scan(
			&f.ID, &f.Name, &f.Category, &f.DefaultWeight); err != nil {
			glog.Error(err)
			continue
		}

		features = append(features, f)
	}

	w.Write(json.MarshalPanic(features))
}

// ByCategory handles "/features/by_category" route.
func ByCategory(c *context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category := vars["category"]

	rows, err := c.DB.Query(selectFeatures+`
		WHERE f.id >= $1
        AND LOWER(f.category) = LOWER($2)
        ORDER BY f.id ASC
        LIMIT $3`,
		c.SinceID,
		category,
		c.PerPage)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	features := make([]model.Feature, 0)

	for rows.Next() {
		var f model.Feature
		if err := rows.Scan(&f.ID, &f.Name, &f.Category, &f.DefaultWeight); err != nil {
			glog.Error(err)
			continue
		}

		features = append(features, f)
	}

	w.Write(json.MarshalIndentPanic(features))
}

// ShowScores handles "/features/{name:[a-zA-Z0-9_]+}/scores" route.
func ShowScores(c *context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	rows, err := c.DB.Query(`
		SELECT s.score, u.id, u.username
		FROM scores AS s
		INNER JOIN users AS u ON s.user_id = u.id
		INNER JOIN features AS f ON f.id = s.feature_id
		WHERE LOWER(f.name) = LOWER($1)
		AND u.id >= $2
		ORDER BY u.id
		LIMIT $3`,
		name,
		c.SinceID,
		c.PerPage)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := make([]user, 0)

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.Score, &u.ID, &u.Username); err != nil {
			glog.Error(err)
			continue
		}

		users = append(users, u)
	}

	w.Write(json.MarshalPanic(users))
}
