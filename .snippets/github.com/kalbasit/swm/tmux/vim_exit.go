package tmux

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	psPath string

	psGrepRegexp = regexp.MustCompile(`(?m)^(?:[^TXZ ])+ +(?:((?:\S)+/))?g?(n?vim?x?)$`)
)

func init() {
	var err error
	psPath, err = exec.LookPath("ps")
	if err != nil {
		log.Fatal().Msgf("error looking up the ps executable, is it installed? %s", err)
	}
}

func (t *tmux) VimExit() error {
	// get the list of panes that are running vim
	targets, err := t.getTargetsRunningVim()
	if err != nil {
		return err
	}
	log.Debug().Msgf("found the following tmux targets running vim: %v", targets)
	// iterate over all the panes that has vim, and ask it to close itself
	for _, target := range targets {
		// Send the escape key, in the case we are in a vim like program. This is
		// repeated because the send-key command is not waiting for vim to complete
		// its action.
		// Credit: https://gist.github.com/debugish/2773454
		for i := 0; i < 25; i++ {
			if err := exec.Command(tmuxPath, "-L", t.socketName(), "send-keys", "-t", target, "C-[").Run(); err != nil {
				return err
			}
		}
		// ask Vim to exit
		if err := exec.Command(tmuxPath, "-L", t.socketName(), "send-keys", "-t", target, ":", "x", "a", "C-m").Run(); err != nil {
			return err
		}
	}

	return nil
}

func (t *tmux) getTargetsRunningVim() ([]string, error) {
	var targets []string

	// get the list of sessions
	var sessionNames []string
	{
		cmd := exec.Command(tmuxPath, "-L", t.socketName(), "list-sessions", "-F", "#{session_name}")
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(string(out), "\n") {
			if line != "" {
				sessionNames = append(sessionNames, line)
			}
		}
	}
	log.Debug().Msgf("found the following tmux sessions: %v", sessionNames)
	// iterate over the list of sessions, and for each session iterate over the
	// list of windows, then over the panes and check what they are running
	for _, sessionName := range sessionNames {
		// get the list of windows for this session
		var windowIDs []string
		{
			cmd := exec.Command(tmuxPath, "-L", t.socketName(), "list-windows", "-t", sessionName, "-F", "#I")
			out, err := cmd.Output()
			if err != nil {
				return nil, err
			}
			for _, line := range strings.Split(string(out), "\n") {
				if line != "" {
					windowIDs = append(windowIDs, line)
				}
			}
		}
		// iterate over the list of windows, get the list of panes, and for each pane
		// find out if it is running vim, nvim then will build targets
		for _, windowID := range windowIDs {
			// build the map of pane to tty
			paneTTY := make(map[string]string)
			{
				cmd := exec.Command(tmuxPath, "-L", t.socketName(), "list-panes", "-t", sessionName+":"+windowID, "-F", "#P@#{pane_tty}")
				out, err := cmd.Output()
				if err != nil {
					return nil, err
				}
				paneTTYs := strings.Split(string(out), "\n")
				for _, paneTTYStr := range paneTTYs {
					paneTTYArr := strings.Split(paneTTYStr, "@")
					if len(paneTTYArr) < 2 {
						continue
					}
					paneTTY[paneTTYArr[0]] = paneTTYArr[1]
				}
			}
			// now iterate over the pane/tty, check what is running on that TTY
			for paneID, ttyPath := range paneTTY {
				// TODO: replace this sub-exec with real /proc parsing library for processes
				// see https://github.com/mitchellh/go-ps/blob/master/process_linux.go
				cmd := exec.Command(psPath, "-o", "state=", "-o", "comm=", "-t", ttyPath)
				out, err := cmd.Output()
				if err != nil {
					return nil, err
				}
				// test the output against the psGrepRegexp
				if psGrepRegexp.Match(out) {
					targets = append(targets, fmt.Sprintf("%s:%s.%s", sessionName, windowID, paneID))
				}
			}
		}

	}

	return targets, nil
}
