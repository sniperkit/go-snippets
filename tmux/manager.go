package tmux

import "errors"

var (
	// ErrProjectNotFoundForGivenSessionName is returned by SwitchClient if the
	// selected session (via fzf currently) was not found. This usually means
	// that fzf output was not one of the input.
	ErrProjectNotFoundForGivenSessionName = errors.New("project not found for the given session name")

	// ErrVimSessionFound is returned by KillServer(closeVim bool) if a vim was
	// found running on the server and closeVim is false
	ErrVimSessionFound = errors.New("vim was found, cannot exit server to avoid data loss")
)

// Manager represents a TMUX manager
type Manager interface {
	// SwitchClient switches the TMUX to a different client
	// KiilPane when set, will close the TMUX pane running swm
	SwitchClient(killPane bool) error

	// VimExit will close any running vim for this session, saving any changed
	// file.
	VimExit() error

	// KillServer kills the TMUX server. If closeVim is false, and Vim sessions
	// were found running, ErrVimSessionFound is returned.
	KillServer(closeVim bool) error
}
