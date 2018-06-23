// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package srv provides utilities and functions to make the server running
// possible.
package srv

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"

	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/DevMine/api-server/api"
	"github.com/DevMine/api-server/api/features"
	repos "github.com/DevMine/api-server/api/repositories"
	"github.com/DevMine/api-server/api/search"
	"github.com/DevMine/api-server/api/stats"
	"github.com/DevMine/api-server/api/users"
	"github.com/DevMine/api-server/config"
	"github.com/DevMine/api-server/srv/context"
	"github.com/DevMine/api-server/util/httputil"
)

type handler func(c *context.Context, w http.ResponseWriter, r *http.Request)

// makeHandler creates the handler function prototype
func makeHandler(db *sql.DB, h handler, cors bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				glog.Error(err)
				ise := http.StatusInternalServerError
				he := httputil.NewResponseError(http.StatusText(ise))
				http.Error(w, he.JSON(), ise)
			}
		}()

		c, err := context.NewContext(db, r)
		if err != nil {
			panic(err)
		}

		// we only serve JSON
		w.Header().Set("Content-Type", "application/json")

		// the server only accepts GET requests
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		// enable Cross Origin Resource Sharing
		if cors {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers",
				"Origin, Accept, Content-Type, X-Requested-With, X-CSRF-Token")
		}

		requestURI, err := url.QueryUnescape(r.RequestURI)
		if err != nil {
			requestURI = r.RequestURI
		}
		glog.Infof("%s %s from %s", r.Method, requestURI, r.RemoteAddr)
		h(c, w, r)
	}
}

// notFoundHandler handles all Not Found errors (ie 404 errors).
func notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf(
			"Mmmh... It looks like you're lost in the void... Here, read this scroll: %s",
			api.DocURL)
		he := httputil.NewResponseError(msg)
		http.Error(w, he.JSON(), http.StatusNotFound)
	})
}

// SetupRouter creates API routes.
// The cors parameter is used to specify whether to enable
// Cross Origin Resource Sharing or not.
func SetupRouter(db *sql.DB, cors bool) *mux.Router {
	r := mux.NewRouter()

	// 404 Not Found routes
	r.NotFoundHandler = notFoundHandler()

	// default route
	r.HandleFunc("/",
		makeHandler(db, api.Index, cors)).Methods("GET")

	// features
	r.HandleFunc("/features",
		makeHandler(db, features.Index, cors)).Methods("GET")
	r.HandleFunc("/features/by_category/{category:[a-zA-Z]+}",
		makeHandler(db, features.ByCategory, cors)).Methods("GET")
	r.HandleFunc("/features/{name:[a-zA-Z0-9_]+}/scores",
		makeHandler(db, features.ShowScores, cors)).Methods("GET")

	// repositories
	r.HandleFunc("/repositories",
		makeHandler(db, repos.Index, cors)).Methods("GET")
	r.HandleFunc("/repositories/{name:[a-zA-Z0-9\\-_\\.]+}",
		makeHandler(db, repos.Show, cors)).Methods("GET")

	// search
	r.HandleFunc("/search/{query}",
		makeHandler(db, search.Query, cors)).Methods("GET")

	// stats
	r.HandleFunc("/stats",
		makeHandler(db, stats.Index, cors)).Methods("GET")

	// users
	r.HandleFunc("/users",
		makeHandler(db, users.Index, cors)).Methods("GET")
	r.HandleFunc("/users/{username:[a-zA-Z0-9-_\\.]+}",
		makeHandler(db, users.Show, cors)).Methods("GET")
	r.HandleFunc("/users/{username:[a-zA-Z0-9-_\\.]+}/commits",
		makeHandler(db, users.ShowCommits, cors)).Methods("GET")
	r.HandleFunc("/users/{username:[a-zA-Z0-9-_\\.]+}/repositories",
		makeHandler(db, users.ShowRepositories, cors)).Methods("GET")
	r.HandleFunc("/users/{username:[a-zA-Z0-9-_\\.]+}/scores",
		makeHandler(db, users.ShowScores, cors)).Methods("GET")

	return r
}

// OpenDBSession creates a session to the database.
func OpenDBSession(cfg config.DatabaseConfig) (*sql.DB, error) {
	dbURL := fmt.Sprintf(
		"user='%s' password='%s' host='%s' port=%d dbname='%s' sslmode='%s'",
		cfg.UserName, cfg.Password, cfg.HostName, cfg.Port, cfg.DBName, cfg.SSLMode)

	return sql.Open("postgres", dbURL)
}
