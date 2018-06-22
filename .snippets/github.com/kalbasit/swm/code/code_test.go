package code

import (
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/kalbasit/swm/testhelper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// discard logs
	log.Logger = zerolog.New(ioutil.Discard)
}

func TestCodeScan(t *testing.T) {
	// swap the filesystem
	oldAppFS := AppFS
	AppFS = afero.NewMemMapFs()
	defer func() { AppFS = oldAppFS }()
	// create the filesystem we want to scan
	testhelper.CreateProjects(t, AppFS)
	// create a code
	c := New("/code", regexp.MustCompile("^.snapshots$"))
	// define the assertion function
	assertFn := func() {
		// create the expected structs
		p := newProfile(c.(*code), t.Name())
		baseStory := p.addStory(baseStoryName)
		baseStory.addProject("github.com/kalbasit/swm")
		baseStory.addProject("github.com/kalbasit/dotfiles")
		baseStory.addProject("github.com/kalbasit/workflow")
		story123 := p.addStory("STORY-123")
		story123.addProject("github.com/kalbasit/swm")
		story123.addProject("github.com/kalbasit/dotfiles")

		// get the profile
		pTest, err := c.(*code).getProfile(t.Name())
		require.NoError(t, err)

		// assert the base story
		assert.Equal(t, baseStory.name, pTest.getStory("base").name)
		assert.Equal(t, baseStory.profile.name, pTest.getStory("base").profile.name)
		for _, importPath := range []string{"github.com/kalbasit/swm", "github.com/kalbasit/dotfiles", "github.com/kalbasit/workflow"} {
			expected, err := baseStory.getProject(importPath)
			require.NoError(t, err)
			got, err := pTest.getStory("base").getProject(importPath)
			require.NoError(t, err)

			assert.Equal(t, expected.importPath, got.importPath)
		}

		// assert the STORY-123 story
		assert.Equal(t, story123.name, pTest.getStory("STORY-123").name)
		assert.Equal(t, story123.profile.name, pTest.getStory("STORY-123").profile.name)
		for _, importPath := range []string{"github.com/kalbasit/swm", "github.com/kalbasit/dotfiles"} {
			expected, err := story123.getProject(importPath)
			require.NoError(t, err)
			got, err := pTest.getStory("STORY-123").getProject(importPath)
			require.NoError(t, err)

			assert.Equal(t, expected.importPath, got.importPath)
		}
	}
	// scan now
	require.NoError(t, c.Scan())
	assertFn()
}

func TestCodeProfile(t *testing.T) {
	// swap the filesystem
	oldAppFS := AppFS
	AppFS = afero.NewMemMapFs()
	defer func() { AppFS = oldAppFS }()
	// create the filesystem we want to scan
	testhelper.CreateProjects(t, AppFS)
	// create a code
	c := New("/code", regexp.MustCompile("^.snapshots$"))
	// define the assertion function
	assertFn := func(pTest *profile) {
		// create the expected structs
		p := newProfile(c.(*code), t.Name())
		baseStory := p.addStory(baseStoryName)
		baseStory.addProject("github.com/kalbasit/swm")
		baseStory.addProject("github.com/kalbasit/dotfiles")
		baseStory.addProject("github.com/kalbasit/workflow")
		story123 := p.addStory("STORY-123")
		story123.addProject("github.com/kalbasit/swm")
		story123.addProject("github.com/kalbasit/dotfiles")

		// assert the base story
		assert.Equal(t, baseStory.name, pTest.getStory("base").name)
		assert.Equal(t, baseStory.profile.name, pTest.getStory("base").profile.name)
		for _, importPath := range []string{"github.com/kalbasit/swm", "github.com/kalbasit/dotfiles", "github.com/kalbasit/workflow"} {
			expected, err := baseStory.getProject(importPath)
			require.NoError(t, err)
			got, err := pTest.getStory("base").getProject(importPath)
			require.NoError(t, err)

			assert.Equal(t, expected.importPath, got.importPath)
		}

		// assert the STORY-123 story
		assert.Equal(t, story123.name, pTest.getStory("STORY-123").name)
		assert.Equal(t, story123.profile.name, pTest.getStory("STORY-123").profile.name)
		for _, importPath := range []string{"github.com/kalbasit/swm", "github.com/kalbasit/dotfiles"} {
			expected, err := story123.getProject(importPath)
			require.NoError(t, err)
			got, err := pTest.getStory("STORY-123").getProject(importPath)
			require.NoError(t, err)

			assert.Equal(t, expected.importPath, got.importPath)
		}
	}
	// assert it throws an error before scanning
	_, err := c.Profile(t.Name())
	assert.EqualError(t, err, ErrCoderNotScanned.Error())
	// scan now
	require.NoError(t, c.Scan())
	// get the profile
	p, err := c.Profile(t.Name())
	require.NoError(t, c.Scan())
	assertFn(p.(*profile))
}

func TestPath(t *testing.T) {
	c := &code{path: "/code"}
	assert.Equal(t, "/code", c.Path())
}

func TestProjectByAbsolutePath(t *testing.T) {
	// compute the test name
	testName := t.Name()
	// swap the filesystem
	oldAppFS := AppFS
	AppFS = afero.NewMemMapFs()
	defer func() { AppFS = oldAppFS }()
	// create the filesystem we want to scan
	testhelper.CreateProjects(t, AppFS)
	// create a code
	c := New("/code", regexp.MustCompile("^.snapshots$"))
	// scan now
	require.NoError(t, c.Scan())

	type desc struct {
		profile    string
		story      string
		importPath string
	}
	tests := map[string]desc{
		"/code/" + testName + "/base/src/github.com/kalbasit/swm": desc{
			profile:    testName,
			story:      baseStoryName,
			importPath: "github.com/kalbasit/swm",
		},

		"/code/" + testName + "/base/src/github.com/kalbasit/swm/cmd": desc{
			profile:    testName,
			story:      baseStoryName,
			importPath: "github.com/kalbasit/swm",
		},

		"/code/" + testName + "/stories/STORY-123/src/github.com/kalbasit/swm": desc{
			profile:    testName,
			story:      "STORY-123",
			importPath: "github.com/kalbasit/swm",
		},

		"/code/" + testName + "/stories/STORY-123/src/github.com/kalbasit/swm/cmd": desc{
			profile:    testName,
			story:      "STORY-123",
			importPath: "github.com/kalbasit/swm",
		},
	}

	for p, expdec := range tests {
		prj, err := c.ProjectByAbsolutePath(p)
		if assert.NoError(t, err) {
			assert.Equal(t, expdec.profile, prj.Story().Profile().Name())
			assert.Equal(t, expdec.story, prj.Story().Name())
			assert.Equal(t, expdec.importPath, prj.ImportPath())
		}
	}

	var err error
	_, err = c.ProjectByAbsolutePath("/code/not-existing/base")
	assert.EqualError(t, err, ErrPathIsInvalid.Error())
	_, err = c.ProjectByAbsolutePath("/code/not-existing/base/src/github.com/kalbasit/swm")
	assert.EqualError(t, err, ErrProfileNoFound.Error())
	_, err = c.ProjectByAbsolutePath("/code/" + testName + "/stories/NOT_HERE/src/github.com/kalbasit/swm")
	assert.EqualError(t, err, ErrStoryNotFound.Error())
	_, err = c.ProjectByAbsolutePath("/code/" + testName + "/base/src/github.com/user/repo")
	assert.EqualError(t, err, ErrProjectNotFound.Error())
}
