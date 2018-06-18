package zkutils

import (
	log "github.com/nickbruun/gocommons/logging"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"time"
)

// Connection session watcher.
//
// Watches for possible losses of sessions.
type connSessionWatcher struct {
	ec             <-chan zk.Event
	sessionTimeout time.Duration
	recvTimeout    time.Duration
	watchers       []chan bool
	lock           sync.Mutex
}

// New connection session watcher.
func newConnSessionWatcher(ec <-chan zk.Event, sessionTimeout time.Duration, recvTimeout time.Duration) *connSessionWatcher {
	w := &connSessionWatcher{
		ec:             ec,
		sessionTimeout: sessionTimeout,
		recvTimeout:    recvTimeout,
		watchers:       make([]chan bool, 0),
	}

	go w.watch()

	return w
}

// Publish the possible loss of a session.
func (w *connSessionWatcher) publish(expired bool) {
	w.lock.Lock()
	defer w.lock.Unlock()

	for _, ch := range w.watchers {
		ch <- expired
		close(ch)
	}
	w.watchers = make([]chan bool, 0)
}

// Watch a session from the point where it was acquired until it is somehow
// lost.
//
// Returns true if the event channel has been closed.
func (w *connSessionWatcher) watchSession() (done bool) {
	// Wait for the either the session to actively expire or the connection to
	// disconnect.
	for {
		ev, ok := <-w.ec
		if !ok {
			return true
		}

		if ev.Type != zk.EventSession {
			continue
		}

		switch ev.State {
		case zk.StateExpired:
			log.Warn("Session expired")
			w.publish(true)
			return false

		case zk.StateDisconnected:
			// Loop until we either hit the possible expiration deadline or the
			// session is reacquired.
			deadlineDur := w.sessionTimeout - w.recvTimeout
			deadlineChan := time.After(deadlineDur)
			log.Warnf("Connection to ZooKeeper lost, considering session lost if not resolved in %s", deadlineDur)
			resolved := false

			for !resolved {
				select {
				case ev, ok = <-w.ec:
					if !ok {
						return true
					}

					if ev.Type != zk.EventSession {
						break
					}

					switch ev.State {
					case zk.StateHasSession:
						log.Info("Connection to ZooKeeper restablished without session loss")
						resolved = true

					case zk.StateExpired:
						log.Warn("Session expired")
						w.publish(true)
						return false
					}

				case <-deadlineChan:
					log.Warnf("Session considered lost due to lack of reestablished session after %s", deadlineDur)
					w.publish(false)
					return false
				}
			}
		}
	}
}

// Watch.
func (w *connSessionWatcher) watch() {
	// Wait for us to get a session.
	for {
		ev, ok := <-w.ec
		if !ok {
			break
		}

		if ev.Type != zk.EventSession || ev.State != zk.StateHasSession {
			continue
		}

		// Watch the newly established session.
		if w.watchSession() {
			break
		}
	}

	// Close all outstanding watchers by indicating a hard session loss, as the
	// end of the event channel means the connection has been closed.
	log.Debug("Done watching sessions as event channel has closed")

	w.publish(true)
}

// Add watcher.
//
// Returns a one-shot channel which will emit a boolean value indicating
// whether or not the session loss is due to explicit expiration.
func (w *connSessionWatcher) AddWatcher() <-chan bool {
	ch := make(chan bool, 1)

	w.lock.Lock()
	defer w.lock.Unlock()

	w.watchers = append(w.watchers, ch)

	return ch
}
