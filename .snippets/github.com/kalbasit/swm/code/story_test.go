package code

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/kalbasit/swm/testhelper"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoryProfile(t *testing.T) {
	// create a new story
	s := &story{
		name: "base",
		profile: &profile{
			name: t.Name(),
			code: &code{
				path: "/code",
			},
		},
	}
	assert.Equal(t, s.profile, s.Profile())
}

func TestStoryName(t *testing.T) {
	// create a new story
	s := &story{
		name: "base",
		profile: &profile{
			name: t.Name(),
			code: &code{
				path: "/code",
			},
		},
	}
	assert.Equal(t, s.name, s.Name())
}

func TestStoryBase(t *testing.T) {
	// test with a base story

	// create a new story
	s := &story{
		name: "base",
		profile: &profile{
			name: t.Name(),
			code: &code{
				path: "/code",
			},
		},
	}
	assert.True(t, s.Base())

	// test with a real story

	// create a new story
	s = &story{
		name: "STORY-123",
		profile: &profile{
			name: t.Name(),
			code: &code{
				path: "/code",
			},
		},
	}
	assert.False(t, s.Base())
}

func TestStoryGoPath(t *testing.T) {
	// testing the case of base

	// create a new story
	s := &story{
		name: "base",
		profile: &profile{
			name: t.Name(),
			code: &code{
				path: "/code",
			},
		},
	}
	// assert the Path
	assert.Equal(t, "/code/"+t.Name()+"/base", s.GoPath())

	// testing the case of a story

	// create a new story
	s = &story{
		name: "STORY-123",
		profile: &profile{
			name: t.Name(),
			code: &code{
				path: "/code",
			},
		},
	}
	// assert the Path
	assert.Equal(t, "/code/"+t.Name()+"/stories/STORY-123", s.GoPath())
}

func TestStoryExist(t *testing.T) {
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

	t.Run("testing the case of a story existing", func(t *testing.T) {
		parts := strings.Split(t.Name(), "/")
		p, err := c.Profile(parts[0])
		require.NoError(t, err)
		s := p.Story("STORY-123")
		assert.True(t, s.Exists())
	})

	t.Run("testing the case of a story", func(t *testing.T) {
		parts := strings.Split(t.Name(), "/")
		p, err := c.Profile(parts[0])
		require.NoError(t, err)
		s := p.Story("not-existing")
		assert.False(t, s.Exists())
	})
}

func TestStoryProjects(t *testing.T) {
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
		// get a story
		s := p.Base()
		// get the projects
		prjs := s.Projects()
		// assert they are the same as the projects
		for _, prj2 := range s.(*story).projects {
			assert.Contains(t, prjs, Project(prj2))
		}
	}
}

func TestStoryProject(t *testing.T) {
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

		// testing with a base story

		s := p.Base()
		for importPath, expectedPrj := range s.(*story).projects {
			prj, err := s.Project(importPath)
			if assert.NoError(t, err) {
				assert.Equal(t, Project(expectedPrj), prj)
				_, err := AppFS.Stat(prj.Path())
				assert.NoError(t, err)
			}
		}

		// testing with story that does exist

		s = p.Story("STORY-123")
		for importPath, expectedPrj := range s.(*story).projects {
			prj, err := s.Project(importPath)
			if assert.NoError(t, err) {
				assert.Equal(t, Project(expectedPrj), prj)
				_, err := AppFS.Stat(prj.Path())
				assert.NoError(t, err)
			}
		}

		// testing with a story that does not exist (no test for ensure here)

		s = p.Story("STORY-123")
		// prepare the expected
		expectedPrj, err := p.Base().Project("github.com/kalbasit/workflow")
		if assert.NoError(t, err) {
			expectedPrj.(*project).story = s.(*story)
			prj, err := s.Project("github.com/kalbasit/workflow")
			if assert.NoError(t, err) {
				assert.Equal(t, Project(expectedPrj), prj)
				_, err := AppFS.Stat(prj.Path())
				assert.True(t, os.IsNotExist(err))
			}
		}
	}
}
