// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package search handles /search... routes.
package search

import (
	stdjson "encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/DevMine/api-server/cache"
	"github.com/DevMine/api-server/score"
	"github.com/DevMine/api-server/srv/context"
	"github.com/DevMine/api-server/util/httputil"
	"github.com/DevMine/api-server/util/json"
)

// numberOfResults represents the number of results to return
// from a search query.
const numberOfResults = 1000

// Query handles "/search/{query}" route.
func Query(c *context.Context, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := map[string]int64{}

	if err := stdjson.Unmarshal([]byte(vars["query"]), &query); err != nil {
		he := httputil.NewResponseError("invalid JSON input")
		http.Error(w, he.JSON(), http.StatusBadRequest)
		return
	}

	featuresNames := cache.GetFeaturesNames()

	for feat, weight := range query {
		if _, ok := featuresNames[feat]; !ok {
			he := httputil.NewResponseError(fmt.Sprintf("non existing feature: %s", feat))
			http.Error(w, he.JSON(), http.StatusBadRequest)
			return
		}

		if weight < 0 {
			he := httputil.NewResponseError("negative weight given")
			http.Error(w, he.JSON(), http.StatusBadRequest)
			return
		}
	}

	ranks, err := score.Rank(c.DB, query)
	if err != nil {
		panic(err)
	}

	// return only 1000 first results
	w.Write(json.MarshalPanic(ranks[:numberOfResults]))
}
