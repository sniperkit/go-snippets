package testhelper

import (
	"testing"

	"github.com/spf13/afero"
)

// CreateProjects creates projects in a filesystem to prepare a coder
func CreateProjects(t *testing.T, appFS afero.Fs) {
	name := t.Name()

	// projects outside GOPATH
	appFS.MkdirAll("/code/"+name+"/invalid/repo/.git", 0755)

	// projects under base
	appFS.MkdirAll("/code/"+name+"/base/src/github.com/kalbasit/swm/.git", 0755)
	appFS.MkdirAll("/code/"+name+"/base/src/github.com/kalbasit/dotfiles/.git", 0755)
	appFS.MkdirAll("/code/"+name+"/base/src/github.com/kalbasit/workflow/.git", 0755)

	// projects under STORY-123
	appFS.MkdirAll("/code/"+name+"/stories/STORY-123/src/github.com/kalbasit/dotfiles", 0755)
	afero.WriteFile(appFS, "/code/"+name+"/stories/STORY-123/src/github.com/kalbasit/dotfiles/.git", []byte(
		"gitdir: /code/"+name+"/base/src/github.com/kalbasit/.git/worktrees/dotfiles",
	), 0644)
	appFS.MkdirAll("/code/"+name+"/stories/STORY-123/src/github.com/kalbasit/swm", 0755)
	afero.WriteFile(appFS, "/code/"+name+"/stories/STORY-123/src/github.com/kalbasit/swm/.git", []byte(
		"gitdir: /code/"+name+"/base/src/github.com/kalbasit/.git/worktrees/swm",
	), 0644)

	// projects ignored

	// projects outside GOPATH
	appFS.MkdirAll("/code/.snapshots/"+name+"/invalid/repo/.git", 0755)

	// projects under base
	appFS.MkdirAll("/code/.snapshots/"+name+"/base/src/github.com/kalbasit/swm/.git", 0755)
	appFS.MkdirAll("/code/.snapshots/"+name+"/base/src/github.com/kalbasit/dotfiles/.git", 0755)

	// projects under STORY-123
	appFS.MkdirAll("/code/.snapshots/"+name+"/stories/STORY-123/src/github.com/kalbasit/dotfiles", 0755)
	afero.WriteFile(appFS, "/code/.snapshots/"+name+"/stories/STORY-123/src/github.com/kalbasit/dotfiles/.git", []byte(
		"gitdir: /code/"+name+"/base/src/github.com/kalbasit/.git/worktrees/dotfiles",
	), 0644)
}
