// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// User represents a user.
type User struct {
	ID       *int64  `json:"id"`
	Username *string `json:"username"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	GhUser   *GhUser `json:"gh_user,omitempty"`
}
