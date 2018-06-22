package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectStory(t *testing.T) {
	// create a story
	s := &story{}
	// create a project
	p := &project{story: s}

	assert.Equal(t, s, p.Story())
}

func TestProjectPath(t *testing.T) {
	// create a new project
	p := &project{
		story: &story{
			name: "base",
			profile: &profile{
				name: "personal",
				code: &code{
					path: "/home/kalbasit/code",
				},
			},
		},

		importPath: "github.com/kalbasit/swm",
	}
	// assert the Path
	assert.Equal(t, "/home/kalbasit/code/personal/base/src/github.com/kalbasit/swm", p.Path())
}

func TestProjectEnsure(t *testing.T) { t.Skip("not implemented yet") }

func TestProjectImportPath(t *testing.T) {
	assert.Equal(t, "github.com/kalbasit/swm", (&project{importPath: "github.com/kalbasit/swm"}).ImportPath())
}

func TestProjectOwner(t *testing.T) {
	// create a new project
	p := &project{
		story: &story{
			name: "base",
			profile: &profile{
				name: "personal",
				code: &code{
					path: "/home/kalbasit/code",
				},
			},
		},

		importPath: "github.com/kalbasit/swm",
	}
	// assert the owner
	assert.Equal(t, "kalbasit", p.Owner())
}

func TestProjectRepo(t *testing.T) {
	// create a new project
	p := &project{
		story: &story{
			name: "base",
			profile: &profile{
				name: "personal",
				code: &code{
					path: "/home/kalbasit/code",
				},
			},
		},

		importPath: "github.com/kalbasit/swm",
	}
	// assert the repo
	assert.Equal(t, "swm", p.Repo())
}
