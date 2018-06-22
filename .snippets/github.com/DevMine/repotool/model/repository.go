// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package model

// Repository represents a source code repository.
type Repository struct {
	// Name is the name of the repository.
	Name string `json:"name"`

	// VCS is the VCS type of the repository (Git, Mercurial, ...).
	VCS string `json:"vcs"`

	// CloneURL represents the URL from which the repository was cloned.
	CloneURL string `json:"clone_url"`

	// ClonePath is the absolute path to which the repository was cloned
	// on the file system.
	ClonePath string `json:"clone_path"`

	// DefaultBranch is the branch that was active when the repository
	// information were obtained..
	DefaultBranch string `json:"default_branch"`

	// Commits is the list of commits of a repository.
	// Note that only the commit of the default branch are retrieved.
	Commits []Commit `json:"commits"`
}
