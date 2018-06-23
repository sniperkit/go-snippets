// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package score provides a way to rank users based on their features scores
// and a user query.
package score

import (
	"database/sql"
	"sort"

	mx "code.google.com/p/biogo.matrix"

	"github.com/DevMine/api-server/cache"
	"github.com/DevMine/api-server/model"
)

// constructWeightVector creates the weight vector from default weight values
// for each features from the database. Default weight values are overwritten
// with the values in "featsWeightQuery". Key value of featsWeightQuery
// corresponds to the 'name' column in the features table.
func constructWeightVector(featsWeightQuery map[string]int64) (*mx.Dense, error) {
	features := cache.GetFeatures()

	weightVector := make([][]float64, len(features))
	for i, f := range features {
		w, ok := featsWeightQuery[*f.Name]
		if !ok {
			w = *f.DefaultWeight
		}
		weightVector[i] = []float64{float64(w)}
	}

	w, err := mx.NewDense(weightVector)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// computeRanks computes the ranks vector using a weighted sum
// by dot product between matrix a and vector w.
func computeRanks(scores mx.Matrix, weightVector mx.Matrix) *mx.Dense {
	res := new(mx.Dense)
	_ = scores.Dot(weightVector, res)

	return res
}

// Rank returns search results, sorted by rank.
func Rank(db *sql.DB, featsWeightQuery map[string]int64) (model.SearchResults, error) {
	sm := cache.GetScoresMatrix()
	uv := cache.GetUsersVector()

	// weight vector
	w, err := constructWeightVector(featsWeightQuery)
	if err != nil {
		return nil, err
	}

	// rank matrix
	rm := computeRanks(sm, w)

	rows, _ := rm.Dims()
	results := make(model.SearchResults, rows)

	for i := 0; i < rows; i++ {
		results[i].User = uv[i]
		results[i].Rank = rm.Row(i)[0]
	}

	sort.Sort(sort.Reverse(results))

	return results, nil
}
