// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// Feature represents a feature.
type Feature struct {
	ID            *int64  `json:"id"`
	Name          *string `json:"name"`
	Category      *string `json:"category"`
	DefaultWeight *int64  `json:"default_weight"`
}
