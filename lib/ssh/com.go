package ssh

import (
	"net"

	"github.com/gliderlabs/com/log"
)

// Component is the primary export of the ssh package
type Component struct {
	Config

	Log log.InfoLogger `com:"singleton"`

	SessionHandler SessionHandler `com:"singleton"`
	AuthHandlers   []AuthHandler  `com:"extpoint"`

	listener net.Listener
}
