package zkutils

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func getPort() int {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}

func getClusterServerAddrs(c *zk.TestCluster) (serverAddrs []string) {
	serverAddrs = make([]string, len(c.Servers))

	for i, s := range c.Servers {
		serverAddrs[i] = fmt.Sprintf("127.0.0.1:%d", s.Port)
	}

	return
}

func StartTestCluster(size int, stdout, stderr io.Writer) (*zk.TestCluster, error) {
	tmpPath, err := ioutil.TempDir("", "gozk")
	if err != nil {
		return nil, err
	}
	success := false
	cluster := &zk.TestCluster{Path: tmpPath}
	defer func() {
		if !success {
			cluster.Stop()
		}
	}()
	for serverN := 0; serverN < size; serverN++ {
		srvPath := filepath.Join(tmpPath, fmt.Sprintf("srv%d", serverN))
		if err := os.Mkdir(srvPath, 0700); err != nil {
			return nil, err
		}
		port := getPort()
		cfg := zk.ServerConfig{
			ClientPort: port,
			DataDir:    srvPath,
		}
		for i := 0; i < size; i++ {
			cfg.Servers = append(cfg.Servers, zk.ServerConfigServer{
				ID:                 i + 1,
				Host:               "127.0.0.1",
				PeerPort:           getPort(),
				LeaderElectionPort: getPort(),
			})
		}
		cfgPath := filepath.Join(srvPath, "zoo.cfg")
		fi, err := os.Create(cfgPath)
		if err != nil {
			return nil, err
		}
		err = cfg.Marshall(fi)
		fi.Close()
		if err != nil {
			return nil, err
		}

		fi, err = os.Create(filepath.Join(srvPath, "myid"))
		if err != nil {
			return nil, err
		}
		_, err = fmt.Fprintf(fi, "%d\n", serverN+1)
		fi.Close()
		if err != nil {
			return nil, err
		}

		srv := &zk.Server{
			ConfigPath: cfgPath,
			Stdout:     stdout,
			Stderr:     stderr,
		}
		if err := srv.Start(); err != nil {
			return nil, err
		}
		cluster.Servers = append(cluster.Servers, zk.TestServer{
			Path: srvPath,
			Port: cfg.ClientPort,
			Srv:  srv,
		})
	}
	success = true

	// Wait for the cluster to become active.
	var probeConn *zk.Conn
	if probeConn, _, err = cluster.ConnectAll(); err != nil {
		cluster.Stop()
		return nil, err
	}
	defer probeConn.Close()

	for i := 0; i < 1000; i++ {
		if _, _, err := probeConn.Exists("/"); err == nil {
			return cluster, nil
		}

		time.Sleep(10 * time.Millisecond)
	}

	cluster.Stop()
	return nil, ErrTimeout
}

// Create a test cluster of a given size.
func CreateTestCluster(t *testing.T, size int) (testCluster *zk.TestCluster, serverAddrs []string) {
	var err error

	// Create the test cluster.
	testCluster, err = StartTestCluster(size, os.Stdout, os.Stdout)
	if err != nil {
		t.Fatalf("Failed to create test cluster: %v", err)
	}

	// Construct server addresses and create a connection.
	serverAddrs = getClusterServerAddrs(testCluster)
	return
}

// Create a test cluster of a given size and a connection to the cluster.
func CreateTestClusterAndConnMan(t *testing.T, size int) (testCluster *zk.TestCluster, connMan *ConnMan) {
	var serverAddrs []string
	testCluster, serverAddrs = CreateTestCluster(t, size)

	// Create a connection.
	var err error
	connMan, err = Connect(serverAddrs, 10*time.Second)
	if err != nil {
		testCluster.Stop()
		t.Fatalf("Failed to create connection to test cluster: %v", err)
	}

	return
}
