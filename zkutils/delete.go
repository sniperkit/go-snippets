package zkutils

import (
	log "github.com/nickbruun/gocommons/logging"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

// Delete a node safely.
//
// Will attempt to delete a node continuously until either an unrecoverable
// error is encountered, or the node is gone.
func DeleteSafely(conn *zk.Conn, path string) {
	for {
		err := conn.Delete(path, -1)

		if err == nil {
			log.Debugf("Removed node: %s", path)
			return
		} else if err == zk.ErrNoNode {
			log.Debugf("Node no longer exists: %s", path)
			return
		} else if !IsErrorRecoverable(err) {
			log.Errorf("Unrecoverable error trying to remove node %s: %v", path, err)
			return
		} else {
			log.Warnf("Failed to remove node %s, waiting 100 ms to retry: %v", path, err)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
