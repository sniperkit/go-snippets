package zkutils

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

// ZooKeeper connection manager.
//
// ZooKeeper connection wrapper that provides multiplexed access to events as
// well as information about session timeout.
type ConnMan struct {
	Conn           *zk.Conn
	SessionTimeout time.Duration
	RecvTimeout    time.Duration
	PingInterval   time.Duration
	em             *EventMultiplexer
	sw             *connSessionWatcher
}

// Connect as connection manager.
func Connect(servers []string, sessionTimeout time.Duration) (*ConnMan, error) {
	// Create the connection.
	conn, ec, err := zk.Connect(servers, sessionTimeout)
	if err != nil {
		return nil, err
	}

	// Set up the connection manager.
	recvTimeout := sessionTimeout * 2 / 3 // Hardcoded in zk.
	em := NewEventMultiplexer(ec)

	cm := &ConnMan{
		Conn:           conn,
		SessionTimeout: sessionTimeout,
		RecvTimeout:    recvTimeout,
		PingInterval:   recvTimeout / 2, // Hardcoded in zk.
		em:             em,
		sw:             newConnSessionWatcher(em.Subscribe(), sessionTimeout, recvTimeout),
	}

	return cm, nil
}

// Close connection.
func (m *ConnMan) Close() {
	m.Conn.Close()
}

// Watch for session loss.
//
// Session loss is indicated if a session expires or connection to a cluster
// is lost for more than the time it is safe to assume that a session is still
// well and alive. If session loss is indicated, it is recommended that any
// caller strictly relying on ephemeral nodes attempt to remove the ephemeral
// node.
//
// Returns a one-shot channel which will emit a boolean value indicating
// whether or not the session loss is due to explicit expiration.
func (m *ConnMan) WatchSessionLoss() <-chan bool {
	return m.sw.AddWatcher()
}

// Subscribe to all events for the connection.
func (m *ConnMan) Subscribe() <-chan zk.Event {
	return m.em.Subscribe()
}
