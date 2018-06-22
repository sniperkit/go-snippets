package code

import "github.com/google/go-github/github"

// Coder defines the coder interface
// code/
// |-- profile1
// |   |-- base
// |   |   |-- src
// |   |       |-- go.import.path
// |   |-- stories
// |   |   |-- STORY-123
// |   |       |-- src
// |   |           |-- go.import.path
// |-- profile2
// |   |-- base
// |   |   |-- src
// |   |       |-- go.import.path
// |   |-- stories
// |   |   |-- STORY-123
// |   |       |-- src
// |   |           |-- go.import.path
type Coder interface {
	// Path returns the absolute path of this coder
	Path() string

	// Profile returns the profile given it's name or an error if no profile with
	// this name was found
	Profile(profile string) (Profile, error)

	// ProjectByAbsolutePath returns the project corresponding to the absolute
	// path.
	ProjectByAbsolutePath(p string) (Project, error)

	// Scan scans the code path
	Scan() error
}

// Profile defines the profile interface
type Profile interface {
	// Coder returns the coder under which this exists
	Coder() Coder

	// Name returns the name of the profile
	Name() string

	// Base returns the base Story
	Base() Story

	// Path returns the absolute path to this profile
	Path() string

	// Story returns the story given it's name or an error if no story with this
	// name was found
	Story(name string) Story
}

// Story defines the story interface
type Story interface {
	// Profile returns the profile under which this story exists
	Profile() Profile

	// Name returns the name of the profile
	Name() string

	// Base returns true if this story is the base story
	Base() bool

	// GoPath returns the absolute GOPATH of this story.
	GoPath() string

	// Projects returns all the projects that are available for this story as
	// well as all the projects for this profile in the base story (with no
	// duplicates). All projects returned from the base story will be a copy of
	// the base project with the story changed. The caller must call Ensure() on
	// a project to make sure it exists (as a worktree) before using it.
	Projects() []Project

	// Project returns the project given the importPath or an error if no project
	// exists. If the project does not exist for this story but does exist in the
	// Base story, it will be copied and story changed. The caller must call
	// Ensure() on the project to make sure it exists (as a worktree) before
	// using it.
	Project(importPath string) (Project, error)

	// AddProject clones url as the new project. Will automatically compute the
	// import path from the given URL.
	AddProject(url string) error

	// Exists returns true if the story does exist on disk
	Exists() bool
}

// Project defines the project interface
type Project interface {
	// Story returns the story to which this project belongs to
	Story() Story

	// Ensure ensures the project exists on disk, by creating a new worktree from
	// the base project, or noop if the worktree already exists on disk.
	Ensure() error

	// Path returns the absolute path to this project
	Path() string

	// ImportPath returns the path under which this project can be imported in Go
	ImportPath() string

	// ListPullRequests returns the list of pull requests
	ListPullRequests() ([]*github.PullRequest, error)
}
