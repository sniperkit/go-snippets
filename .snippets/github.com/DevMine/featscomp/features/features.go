// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package features provides various features implementations.
package features

import (
	"database/sql"

	"github.com/DevMine/featscomp/util"
)

// Feature interface must be implemented by all features.
type Feature interface {
	// Score computes the scores of a feature.
	Score() error

	// Name returns the name of a feature.
	Name() string
}

func deleteOldScores(db *sql.DB, featureName string) error {
	featureID, err := util.FindFeatureID(db, featureName)
	if err != nil {
		return err
	}

	if _, err := db.Exec(`DELETE FROM scores WHERE feature_id = $1`, featureID); err != nil {
		return err
	}
	return nil
}
