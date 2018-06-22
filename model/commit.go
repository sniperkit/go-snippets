// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// Commit is a representation of a VCS commit.
type Commit struct {
	ID               *int64      `json:"id"`
	Repository       *Repository `json:"repository,omitempty"`
	Message          *string     `json:"message"`
	Author           *User       `json:"author,omitempty"`
	Committer        *User       `json:"committer,omitempty"`
	AuthorDate       *time.Time  `json:"author_date"`
	CommitDate       *time.Time  `json:"commit_date"`
	FileChangedCount *int64      `json:"file_changed_count"`
	InsertionsCount  *int64      `json:"insertions_count"`
	DeletionsCount   *int64      `json:"deletions_count"`
}
