// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

import "time"

// Commit is a representation of a VCS commit.
type Commit struct {
	// VCSID represents the VCS interanl identification of a commit.
	// In the case of Git, for instance, it is the commit SHA.
	VCSID string `json:"vcs_id"`

	// Message represents the commits message.
	Message string `json:"message"`

	// Author represents the developer that authored the changes made
	// in the commit.
	Author Developer `json:"author"`

	// Committer represents the developer that commited the changes.
	// Most of the time, the committer is also the author.
	Committer Developer `json:"committer"`

	// AuthorDate represents the date when the commit was created.
	AuthorDate time.Time `json:"author_date"`

	// CommitDate represents the date when the commit was committed.
	CommitDate time.Time `json:"commit_date"`

	// DiffDelta represents the changes maed by the commit.
	DiffDelta []DiffDelta `json:"diff_delta,omitempty"`

	// FileChangedCount reprensents how many files have been touched by the
	// commit.
	FileChangedCount int `json:"file_changed_count"`

	// InsertionsCount represents how many new lines have been added.
	InsertionsCount int `json:"insertions_count"`

	// DeletionsCount represents how many lines have been removed.
	DeletionsCount int `json:"deletions_count"`
}
