// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// SearchResult represents a search result, as from the results of a user query
// which the Rank() function from score package computes.
type SearchResult struct {
	User
	Rank float64 `json:"rank"`
}

// SearchResults is used to store search results and is sortable
// by SearchResult.Rank.
type SearchResults []SearchResult

// Len is the number of elements in a collection of search results.
func (sr SearchResults) Len() int {
	return len(sr)
}

// Less reports whether the element with index i should sort before the element
// with index j.
func (sr SearchResults) Less(i, j int) bool {
	return sr[i].Rank < sr[j].Rank
}

// Swap swaps the elements with indexes i and j.
func (sr SearchResults) Swap(i, j int) {
	sr[i], sr[j] = sr[j], sr[i]
}
