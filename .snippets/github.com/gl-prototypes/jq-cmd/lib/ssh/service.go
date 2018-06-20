package ssh

import (
	"net"

	"github.com/gliderlabs/ssh"
)

// NewServer returns an SSH server configured with authentication and session
// handlers provided by other components.
func NewServer(c *Component) *ssh.Server {
	server := &ssh.Server{}
	server.SetOption(ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		for _, handler := range c.AuthHandlers {
			if handler.HandleAuth(ctx, key) {
				return true
			}
		}
		return false
	}))
	server.Handle(func(sess ssh.Session) {
		c.SessionHandler.HandleSSH(sess)
	})
	return server
}

// Stop will close the listener for the SSH server
func (c *Component) Stop() {
	if c.listener != nil {
		c.listener.Close()
	}
}

// Serve uses NewServer to create a server configured by component configuration
// and starts listening for connections.
func (c *Component) Serve() {
	server := NewServer(c)
	server.SetOption(ssh.HostKeyFile(c.HostKeyPEM))
	server.IdleTimeout = c.IdleTimeout
	server.MaxTimeout = c.MaxTimeout
	var err error
	c.listener, err = net.Listen("tcp", c.ListenAddr)
	if err != nil {
		panic(err)
	}
	c.Log.Infof("listening on %s ...", c.ListenAddr)
	server.Serve(c.listener)
}
