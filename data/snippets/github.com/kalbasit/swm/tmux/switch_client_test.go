package tmux

import (
	"regexp"
	"sort"
	"testing"

	"github.com/kalbasit/swm/code"
	"github.com/kalbasit/swm/testhelper"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSessionProjectsBase(t *testing.T) {
	// swap the filesystem
	oldAppFS := code.AppFS
	code.AppFS = afero.NewMemMapFs()
	defer func() { code.AppFS = oldAppFS }()
	// create the filesystem we want to scan
	testhelper.CreateProjects(t, code.AppFS)
	// create a code and scan
	c := code.New("/code", regexp.MustCompile("^.snapshots$"))
	require.NoError(t, c.Scan())
	// create the tmux client
	tmx := &tmux{options: &Options{
		Coder:   c,
		Profile: t.Name(),
		Story:   "base",
	}}
	// get the session map
	sessionNameProjects, err := tmx.getSessionNameProjects()
	require.NoError(t, err)
	assert.NotEmpty(t, sessionNameProjects)
	// assert the keys in the map
	var keys []string
	for k, _ := range sessionNameProjects {
		keys = append(keys, k)
	}
	expectedKeys := []string{
		"github" + dotChar + "com/kalbasit/swm",
		"github" + dotChar + "com/kalbasit/dotfiles",
		"github" + dotChar + "com/kalbasit/workflow",
	}
	sort.Strings(keys)
	sort.Strings(expectedKeys)
	require.Equal(t, expectedKeys, keys)
	// assert the correct project
	for name, prj := range sessionNameProjects {
		switch name {
		case "github" + dotChar + "com/kalbasit/swm":
			assert.Equal(t, "github.com/kalbasit/swm", prj.ImportPath())
		case "github" + dotChar + "com/kalbasit/dotfiles":
			assert.Equal(t, "github.com/kalbasit/dotfiles", prj.ImportPath())
		case "github" + dotChar + "com/kalbasit/workflow":
			assert.Equal(t, "github.com/kalbasit/workflow", prj.ImportPath())
		}
	}
}

func TestGetSessionProjectsStory123(t *testing.T) {
	// swap the filesystem
	oldAppFS := code.AppFS
	code.AppFS = afero.NewMemMapFs()
	defer func() { code.AppFS = oldAppFS }()
	// create the filesystem we want to scan
	testhelper.CreateProjects(t, code.AppFS)
	// create a code and scan
	c := code.New("/code", regexp.MustCompile("^.snapshots$"))
	require.NoError(t, c.Scan())
	// create the tmux client
	tmx := &tmux{options: &Options{
		Coder:   c,
		Profile: t.Name(),
		Story:   "STORY-123",
	}}
	// get the session map
	sessionNameProjects, err := tmx.getSessionNameProjects()
	require.NoError(t, err)
	assert.NotEmpty(t, sessionNameProjects)
	// assert the keys in the map
	var keys []string
	for k, _ := range sessionNameProjects {
		keys = append(keys, k)
	}
	expectedKeys := []string{
		"github" + dotChar + "com/kalbasit/swm",
		"github" + dotChar + "com/kalbasit/dotfiles",
		"github" + dotChar + "com/kalbasit/workflow",
	}
	sort.Strings(keys)
	sort.Strings(expectedKeys)
	require.Equal(t, expectedKeys, keys)
	// assert the correct project
	for name, prj := range sessionNameProjects {
		switch name {
		case "github" + dotChar + "com/kalbasit/swm":
			assert.Equal(t, "github.com/kalbasit/swm", prj.ImportPath())
		case "github" + dotChar + "com/kalbasit/dotfiles":
			assert.Equal(t, "github.com/kalbasit/dotfiles", prj.ImportPath())
		case "github" + dotChar + "com/kalbasit/workflow":
			assert.Equal(t, "github.com/kalbasit/workflow", prj.ImportPath())
		}
	}
}
