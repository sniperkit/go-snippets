// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package context creates a context to pass to each api-server route
// handlers functions.
package context

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/DevMine/api-server/util/typeutil"
)

// Context represents a context of a query to the API server that is meant
// to be used by route handlers functions.
type Context struct {
	// DB is an active connection to the database.
	DB *sql.DB

	// SinceID corresponds to an ID since which to show results.
	SinceID uint64

	// PerPage corresponds to the number of results to show per page.
	PerPage uint64

	// PageNumber shall be used to paginate results when necessary.
	PageNumber uint64
}

// NewContext initializes a Context structure.
func NewContext(db *sql.DB, r *http.Request) (*Context, error) {
	if db == nil {
		return nil, errors.New("database session cannot be nil")
	}

	if r == nil {
		return nil, errors.New("http request cannot be nil")
	}

	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	params := r.Form

	var sinceID uint64
	sinceID, _ = typeutil.StrToUint(params.Get("since"))

	var perPage uint64
	perPage, _ = typeutil.StrToUint(params.Get("per_page"))
	if perPage > 100 {
		perPage = 100
	} else if perPage == 0 {
		perPage = 30
	}

	var pageNumber uint64
	pageNumber, _ = typeutil.StrToUint(params.Get("page"))
	if pageNumber == 0 {
		pageNumber = 1
	}

	return &Context{DB: db, SinceID: sinceID, PerPage: perPage, PageNumber: pageNumber}, nil
}
