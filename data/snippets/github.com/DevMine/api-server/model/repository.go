// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// Repository represents a repository.
type Repository struct {
	ID              *int64        `json:"id"`
	Name            *string       `json:"name"`
	PrimaryLanguage *string       `json:"primary_language"`
	CloneURL        *string       `json:"clone_url"`
	ClonePath       *string       `json:"clone_path"`
	VCS             *string       `json:"vcs"`
	GhRepository    *GhRepository `json:"gh_repository,omitempty"`
}
