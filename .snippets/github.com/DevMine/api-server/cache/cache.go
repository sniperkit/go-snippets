// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cache provides caching. All data that can be loaded once into memory
// is dealt with in this package.
// To use it, a call to LoadCache() shall be made once at the start of the
// application. After the data is loaded into memory, it can be accessed with
// the getters functions.
package cache

import (
	"database/sql"
	"errors"

	mx "code.google.com/p/biogo.matrix"

	"github.com/DevMine/api-server/model"
)

var (
	features      []model.Feature
	featuresNames map[string]struct{}
	scoresMatrix  *mx.Sparse
	stats         *model.Stats
	usersVector   []model.User

	errCacheNotLoaded = errors.New("cache not loaded")
)

// LoadCache loads all cacheable data into memory.
func LoadCache(db *sql.DB) error {
	if err := loadStats(db); err != nil {
		return err
	}

	if err := loadFeatures(db); err != nil {
		return err
	}

	if err := loadFeaturesNames(db); err != nil {
		return err
	}

	if err := loadScoresAndUsers(db); err != nil {
		return err
	}

	return nil
}

// GetStats provides database statistics.
func GetStats() model.Stats {
	if stats == nil {
		panic(errCacheNotLoaded)
	}
	return *stats
}

// GetFeatures returns a slice containing all features.
func GetFeatures() []model.Feature {
	if features == nil {
		panic(errCacheNotLoaded)
	}
	return features
}

// GetFeaturesNames returns a map of features names. It can be used to check in
// O(1) if a feature exists by providing its name as a key to the map.
func GetFeaturesNames() map[string]struct{} {
	if featuresNames == nil {
		panic(errCacheNotLoaded)
	}
	return featuresNames
}

// GetScoresMatrix returns the scores matrix from cache.
// Each row of the matrix corresponds to a user whereas each column corresponds
// to a feature. The scores matrix is closely related to the users vector.
// Row 'i' of the scores matrix corresponds to the scores for each feature for
// the user at row 'i' in the users vector.
func GetScoresMatrix() *mx.Sparse {
	if scoresMatrix == nil {
		panic(errCacheNotLoaded)
	}
	return scoresMatrix
}

// GetUsersVector returns the vector of users from cache.
// Important: the vector is sorted, matching the rows of the scores matrix.
// In other words, user at position 'i' in the slice corresponds as the scores
// for each feature at row 'i' of the scores matrix.
func GetUsersVector() []model.User {
	if usersVector == nil {
		panic(errCacheNotLoaded)
	}
	return usersVector
}
