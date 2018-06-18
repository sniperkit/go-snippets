package zkutils

import (
	"github.com/samuel/go-zookeeper/zk"
	"testing"
	"time"
)

func AssertNoAwaitResult(t *testing.T, path string, await <-chan error) {
	select {
	case err := <-await:
		if err != nil {
			t.Fatalf("Error awaiting node existence for %s: %v", path, err)
		} else {
			t.Fatalf("Node unexpectedly exists: %s", path)
		}

	case <-time.After(100 * time.Millisecond):
		break
	}
}

func AssertAwaitResult(t *testing.T, path string, await <-chan error) {
	select {
	case err := <-await:
		if err != nil {
			t.Fatalf("Error awaiting node existence for %s: %v", path, err)
		}

	case <-time.After(time.Second):
		t.Fatalf("Timeout waiting for node to exist: %s", path)
	}
}

func TestAwaitExists(t *testing.T) {
	testCluster, conn := CreateTestClusterAndConnMan(t, 1)
	defer testCluster.Stop()
	defer conn.Close()

	// Test for root path.
	await := AwaitExists(conn.Conn, "/")
	AssertAwaitResult(t, "/", await)

	// Test for single-level path.
	await = AwaitExists(conn.Conn, "/one")

	AssertNoAwaitResult(t, "/one", await)

	if err := CreateRecursively(conn.Conn, "/one", zk.WorldACL(zk.PermAll)); err != nil {
		t.Fatalf("Error creaeting /one: %v", err)
	}

	AssertAwaitResult(t, "/one", await)

	// Test for multi-level path.
	await = AwaitExists(conn.Conn, "/two/three")

	AssertNoAwaitResult(t, "/two/three", await)

	if err := CreateRecursively(conn.Conn, "/two", zk.WorldACL(zk.PermAll)); err != nil {
		t.Fatalf("Error creating /two: %v", err)
	}

	AssertNoAwaitResult(t, "/two/three", await)

	if err := CreateRecursively(conn.Conn, "/two/three", zk.WorldACL(zk.PermAll)); err != nil {
		t.Fatalf("Error creating /two/three: %v", err)
	}

	AssertAwaitResult(t, "/two/three", await)
}
