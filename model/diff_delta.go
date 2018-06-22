// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// Status of a file touched by a commit.
var (
	StatusAdded    = "added"
	StatusDeleted  = "deleted"
	StatusModified = "modified"
	StatusRenamed  = "renamed"
	StatusCopied   = "copied"
)

// DiffDelta represents a delta difference between a commit and its parent.
type DiffDelta struct {
	// Patch represents the difference between a commit and its parent.
	Patch *string `json:"patch,omitempty"`

	// Status gives information about whether the file has been added, deleted
	// modified, renamed of copied.
	Status *string `json:"status,omitempty"`

	// Binary gives information about whether the file is a binary or not.
	Binary *bool `json:"binary,omitempty"`

	// Similarity is a score that indicates how similar the file is from its
	// state in the previous commit.
	// A similarity score is a value between 0 and 100 indicating how similar
	// the old and new files are. The higher the value, the more similar the
	// files are.
	Similarity *uint `json:"similarity,omitempty"`

	// OldFilePath represents the path to the old file.
	OldFilePath *string `json:"old_file_path,omitempty"`

	// NewFilePath represents the path to the new file.
	NewFilePath *string `json:"new_file_path,omitempty"`
}
