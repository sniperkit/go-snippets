package tmux

import "os/exec"

func (t *tmux) KillServer(closeVim bool) error {
	// find out if we have any running vim session and if we do, act on closeVim
	targets, err := t.getTargetsRunningVim()
	if err != nil {
		return err
	}
	if len(targets) > 0 {
		if !closeVim {
			return ErrVimSessionFound
		}
		// ask vim to exit
		if err := t.VimExit(); err != nil {
			return err
		}
	}

	return exec.Command(tmuxPath, "-L", t.socketName(), "kill-server").Run()
}
