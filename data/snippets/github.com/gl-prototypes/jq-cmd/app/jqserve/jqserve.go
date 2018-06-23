package jqserve

import (
	"os/exec"

	"github.com/gliderlabs/com"
	"github.com/gliderlabs/ssh"
)

func init() {
	com.DefaultRegistry.Register(&com.Object{Value: &Component{}})
}

type Component struct{}

func (com *Component) HandleSSH(sess ssh.Session) {
	cmd := exec.Command("jq", sess.Command()...)
	help := false
	for _, arg := range sess.Command() {
		if arg == "-h" {
			help = true
		}
	}
	if !help {
		cmd.Stdin = sess
	}
	cmd.Stdout = sess
	cmd.Stderr = sess.Stderr()
	if err := cmd.Run(); err != nil {
		sess.Exit(1)
	}
}

func (com *Component) HandleAuth(ssh.Context, ssh.PublicKey) bool {
	return true
}
