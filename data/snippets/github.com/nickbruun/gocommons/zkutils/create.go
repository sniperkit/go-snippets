package zkutils

import (
	log "github.com/nickbruun/gocommons/logging"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

// Recursively create nodes with no data if they do not exist.
//
// Does not return any error if the node path already exist.
func CreateRecursively(conn *zk.Conn, path string, acl []zk.ACL) (err error) {
	// Test if the node already exists for efficiency reasons.
	var exists bool
	if exists, _, err = conn.Exists(path); exists || err != nil {
		return
	}

	// Start from the root.
	comps := strings.Split(path, "/")
	for i := 1; i < len(comps); i++ {
		path := strings.Join(comps[:i+1], "/")

		_, err = conn.Create(path, nil, 0, acl)
		if err == zk.ErrNodeExists {
			log.Debugf("Node already exists: %s", path)
			err = nil
		} else if err != nil {
			return
		}

		log.Debugf("Created node: %s", path)
	}

	return
}
