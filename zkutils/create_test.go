package zkutils

import (
	"github.com/samuel/go-zookeeper/zk"
	"testing"
)

func AssertCreateRecursivelyCreates(t *testing.T, conn *zk.Conn, path string) {
	if err := CreateRecursively(conn, path, zk.WorldACL(zk.PermAll)); err != nil {
		t.Errorf("Failed to create path %s recursively: %v", path, err)
		return
	}

	exists, _, err := conn.Exists(path)
	if err != nil {
		t.Errorf("Failed to test if path %s exists: %v", path, err)
		return
	}

	if !exists {
		t.Errorf("Expected path %s to exist", path)
	}
}

func TestCreateRecursively(t *testing.T) {
	testCluster, conn := CreateTestClusterAndConnMan(t, 1)
	defer testCluster.Stop()
	defer conn.Close()

	AssertCreateRecursivelyCreates(t, conn.Conn, "/")
	AssertCreateRecursivelyCreates(t, conn.Conn, "/one")
	AssertCreateRecursivelyCreates(t, conn.Conn, "/one")
	AssertCreateRecursivelyCreates(t, conn.Conn, "/two/three/four/five/six")
	AssertCreateRecursivelyCreates(t, conn.Conn, "/two/three/four/five/six")
}
