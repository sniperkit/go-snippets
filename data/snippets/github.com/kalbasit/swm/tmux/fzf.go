package tmux

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	// fzfPath is the PATH to the fzf executable
	fzfPath string
)

func init() {
	var err error
	fzfPath, err = exec.LookPath("fzf")
	if err != nil {
		log.Fatal().Msgf("error looking up the fzf executable, is it installed? %s", err)
	}
}

// withFilter filters input using fzf
func (t *tmux) withFilter(input func(in io.WriteCloser)) (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	cmd := exec.Command(shell, "-c", fzfPath)
	cmd.Stderr = os.Stderr
	in, _ := cmd.StdinPipe()
	go func() {
		input(in)
		in.Close()
	}()
	result, err := cmd.Output()
	if err != nil {
		return "", err
	}
	res := strings.Split(string(result), "\n")
	if res[len(res)-1] == "" {
		res = res[0 : len(res)-1]
	}
	return res[0], nil
}
