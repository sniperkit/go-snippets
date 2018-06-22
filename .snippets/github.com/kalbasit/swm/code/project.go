package code

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/google/go-github/github"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
)

const srcDir = "src"

// gitPath is the PATH to the git binary
var gitPath string

type project struct {
	// story returns the parent story
	story *story

	// importPath is the path of the project relative to the GOPATH/src of the profile/workspace
	importPath string
}

func init() {
	var err error
	gitPath, err = exec.LookPath("git")
	if err != nil {
		log.Fatal().Msgf("error looking up the git executable, is it installed? %s", err)
	}
}

func newProject(s *story, importPath string) *project {
	return &project{
		story:      s,
		importPath: importPath,
	}
}

// Story returns the story to which this project belongs to
func (p *project) Story() Story { return p.story }

// Path returns the absolute path of the project
func (p *project) Path() string { return path.Join(p.story.GoPath(), srcDir, p.importPath) }

// Ensure ensures the project exists on disk, by creating a new worktree from
// the base project or noop if the worktree already exists on disk.
func (p *project) Ensure() error {
	if _, err := AppFS.Stat(p.Path()); os.IsNotExist(err) {
		// get the base project
		baseStory := p.story.Profile().Base()
		baseProject, err := baseStory.Project(p.importPath)
		if err != nil {
			return err
		}
		// run the pre-hooks
		if err = p.runPreHooks(baseProject); err != nil {
			return err
		}
		// create a new worktree for this project based on the base project
		// TODO(kalbasit): switch to using [go-git](https://github.com/src-d/go-git)
		cmd := exec.Command(gitPath, "worktree", "add", "-b", p.story.name, p.Path(), "master")
		cmd.Dir = baseProject.Path()
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error creating a new worktree: %s\nOutput:\n%s", err, string(out))
		}
		// run the post-hooks
		if err = p.runPostHooks(baseProject); err != nil {
			return err
		}
		// add this project to the projects of the story above
		p.story.addProject(p.importPath)
	}

	return nil
}

// ImportPath returns the path under which this project can be imported in Go
func (p *project) ImportPath() string { return p.importPath }

func (p *project) Owner() string {
	parts := strings.Split(p.importPath, "/")
	if len(parts) != 3 {
		return ""
	}
	return parts[1]
}

func (p *project) Repo() string {
	parts := strings.Split(p.importPath, "/")
	if len(parts) != 3 {
		return ""
	}
	return parts[2]
}

func (p *project) ListPullRequests() ([]*github.PullRequest, error) {
	prs, _, err := GithubClient.PullRequests.List(context.Background(), p.Owner(), p.Repo(), nil)
	return prs, err
}

// runPreHooks iterates over the executable files in
// ~/.config/swm/hooks/coder/pre-hook and runs each with the following
// arguments:
// - The GOPATH of the base story
// - The GOPATH of the story
// - The path of the base story
// - The path of the story
func (p *project) runPreHooks(baseProject Project) error {
	// get the hooks directory
	preHooksDir := path.Join(p.hookPath(), "pre-hook")
	// first get the list of the hooks
	hooks, err := afero.ReadDir(AppFS, preHooksDir)
	if err != nil && !os.IsNotExist(err) {
		log.Error().Err(err).Str("pre-hook-dir", preHooksDir).Msg("error reading the directory")
		return err
	} else if os.IsNotExist(err) {
		return nil
	}
	// iterate over the list of hooks and run it
	for _, hook := range hooks {
		// compute the absolute path of the hook
		hookPath := path.Join(preHooksDir, hook.Name())
		log.Debug().
			Str("hook_path", hookPath).
			Bool("executable", hook.Mode().Perm()&0111 != 0).
			Msg("found a pre-hook")
		// is this a file and is executable by the current user?
		if !hook.IsDir() && hook.Mode().Perm()&0111 != 0 {
			cmd := exec.Command(hookPath, p.story.Profile().Base().GoPath(), p.story.GoPath(), baseProject.Path(), p.Path())
			out, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("error running the pre-hook: %s\nOutput:\n%s", err, string(out))
			}
		}
	}

	return nil
}

// runPostHooks iterates over the executable files in
// ~/.config/swm/hooks/coder/post-hook and runs each with the following
// arguments:
// - The GOPATH of the base story
// - The GOPATH of the story
// - The path of the base story
// - The path of the story
func (p *project) runPostHooks(baseProject Project) error {
	// compute the absolute path of the hook
	postHooksDir := path.Join(p.hookPath(), "post-hook")
	// first get the list of the hooks
	hooks, err := afero.ReadDir(AppFS, postHooksDir)
	if err != nil && !os.IsNotExist(err) {
		log.Error().Err(err).Str("post-hook-dir", postHooksDir).Msgf("error reading the directory")
		return err
	} else if os.IsNotExist(err) {
		return nil
	}
	// iterate over the list of hooks and run it
	for _, hook := range hooks {
		hookPath := path.Join(postHooksDir, hook.Name())
		log.Debug().
			Str("hook_path", hookPath).
			Bool("executable", hook.Mode().Perm()&0111 != 0).
			Msg("found a post-hook")
		// is this a file and is executable by the current user?
		if !hook.IsDir() && hook.Mode().Perm()&0111 != 0 {
			cmd := exec.Command(hookPath, p.story.Profile().Base().GoPath(), p.story.GoPath(), baseProject.Path(), p.Path())
			out, err := cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("error running the post-hook: %s\nOutput:\n%s", err, string(out))
			}
		}
	}

	return nil
}

func (p *project) hookPath() string {
	return path.Join(os.Getenv("HOME"), ".config", "swm", "hooks", "coder")
}
