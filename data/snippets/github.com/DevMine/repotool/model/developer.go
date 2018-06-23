// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// Developer represents someone linked to a source code repository, be it
// either as a commiter or commit author (which is not mutually exclusive).
type Developer struct {
	// Name represents the name of a developer.
	Name string `json:"name"`

	// Email is the email of a developer.
	Email string `json:"email"`
}
