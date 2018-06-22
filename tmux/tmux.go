package tmux

import (
	"fmt"
	"os/exec"

	"github.com/kalbasit/swm/code"
	"github.com/rs/zerolog/log"
)

const (
	dotChar   = "\u2022"
	colonChar = "\uFF1A"
)

var (
	// tmuxPath is the PATH to the tmux executable
	tmuxPath string
)

func init() {
	var err error
	tmuxPath, err = exec.LookPath("tmux")
	if err != nil {
		log.Fatal().Msgf("error looking up the tmux executable, is it installed? %s", err)
	}
}

// Options configures the tmux manager
type Options struct {
	// Coder represents the coder instance
	Coder code.Coder

	// Profile represents the profile we are going to use to compute the list of
	// available projects as well as the ACTIVE_PROFILE of new sessions.
	Profile string

	// Story represents the story we are going to use to compute the list of
	// available projects.
	Story string
}

// tmux implements the Manager interface
type tmux struct{ options *Options }

// socketName returns the session name
func (t *tmux) socketName() string { return fmt.Sprintf("%s@%s", t.options.Profile, t.options.Story) }

// New returns a new tmux manager
func New(opts *Options) Manager {
	return &tmux{options: opts}
}
