package code

import (
	"regexp"
	"testing"

	"github.com/kalbasit/swm/testhelper"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileCoder(t *testing.T) {
	// create a new profile
	p := &profile{
		name: t.Name(),
		code: &code{
			path: "/code",
		},
	}
	assert.Equal(t, Coder(p.code), p.Coder())
}

func TestProfileName(t *testing.T) {
	// create a new profile
	p := &profile{
		name: t.Name(),
		code: &code{
			path: "/code",
		},
	}
	assert.Equal(t, p.name, p.Name())
}

func TestProfileBase(t *testing.T) {
	// swap the filesystem
	oldAppFS := AppFS
	AppFS = afero.NewMemMapFs()
	defer func() { AppFS = oldAppFS }()
	// create the filesystem we want to scan
	testhelper.CreateProjects(t, AppFS)
	// create a new code
	c := New("/code", regexp.MustCompile("^.snapshots$"))
	// scan now
	if err := c.Scan(); err != nil {
		t.Fatalf("code scan failed: %s", err)
	}
	// get a profile
	p, err := c.Profile(t.Name())
	if assert.NoError(t, err) {
		assert.Equal(t, Story(p.(*profile).getStory(baseStoryName)), p.Base())
	}
}

func TestProfilePath(t *testing.T) {
	// create a new profile
	p := &profile{
		name: t.Name(),
		code: &code{
			path: "/code",
		},
	}
	// assert the Path
	assert.Equal(t, "/code/"+t.Name(), p.Path())
}

func TestProfileStory(t *testing.T) {
	// swap the filesystem
	oldAppFS := AppFS
	AppFS = afero.NewMemMapFs()
	defer func() { AppFS = oldAppFS }()
	// create the filesystem we want to scan
	testhelper.CreateProjects(t, AppFS)
	// create a new code
	c := New("/code", regexp.MustCompile("^.snapshots$"))
	// scan now
	if err := c.Scan(); err != nil {
		t.Fatalf("code scan failed: %s", err)
	}
	// get the profile
	p, err := c.Profile(t.Name())
	require.NoError(t, err)

	// test with a story that does exist

	s := p.Story("STORY-123")
	assert.Equal(t, p.(*profile).getStory("STORY-123"), s)

	// test with a story that does not exist

	s = p.Story("STORY-456")
	assert.Equal(t, "STORY-456", s.(*story).name)
	assert.NotNil(t, s.(*story).projects)
	assert.Empty(t, s.(*story).projects)
	assert.Equal(t, p.(*profile).getStory("STORY-456"), s)
}
