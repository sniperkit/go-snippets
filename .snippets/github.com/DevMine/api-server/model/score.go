// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// Score represents a score.
type Score struct {
	ID        *int64   `json:"id"`
	UserID    *int64   `json:"-"`
	FeatureID *int64   `json:"-"`
	Score     *float64 `json:"score"`
}
