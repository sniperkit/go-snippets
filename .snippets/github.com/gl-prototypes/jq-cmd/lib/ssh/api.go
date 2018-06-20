package ssh

import "github.com/gliderlabs/ssh"

// SessionHandler is an extension point interface for handling SSH sessions.
type SessionHandler interface {
	HandleSSH(ssh.Session)
}

// AuthHandler is an extension point interface for authenticating SSH connections.
type AuthHandler interface {
	HandleAuth(ssh.Context, ssh.PublicKey) bool
}
