package locking

import (
	logrus "github.com/Sirupsen/logrus"
	log "github.com/nickbruun/gocommons/logging"
	"github.com/nickbruun/gocommons/unittest"
	"github.com/nickbruun/gocommons/zkutils"
	"github.com/samuel/go-zookeeper/zk"
	"testing"
	"time"
)

type ZooKeeperLockTestSuite struct {
	LockTestSuite

	TestCluster *zk.TestCluster
	ConnMan     *zkutils.ConnMan
}

func (s *ZooKeeperLockTestSuite) SetUpSuite() {
	s.TestCluster, s.ConnMan = zkutils.CreateTestClusterAndConnMan(s.T(), 1)
	s.LockProvider = NewZooKeeperLockProvider(s.ConnMan, zk.WorldACL(zk.PermAll))
}

func (s *ZooKeeperLockTestSuite) TearDownSuite() {
	s.ConnMan.Close()
	s.TestCluster.Stop()
}

func (s *ZooKeeperLockTestSuite) TestLockFailureOnNodeRemoval() {
	l := s.LockProvider.GetLock("/lock/path")
	failCh := s.AssertLock(l)
	log.Debug("Acquried lock, removing node")

	children, _, err := s.ConnMan.Conn.Children("/lock/path")
	if err != nil {
		s.Fatalf("Failed to list children of /lock/path: %v", err)
	}

	for _, c := range children {
		s.ConnMan.Conn.Delete("/lock/path/"+c, -1)
	}

	select {
	case <-failCh:
	case <-time.After(time.Second):
		s.Fatalf("Lock did not fail after 1 s")
	}

	if err := l.Unlock(); err != ErrNotLocked {
		s.Fatalf("Expected attempt at unlocking to return ErrNotLocked, but returned: %v", err)
	}
}

func TestZooKeeperLock(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	unittest.RunTestSuite(&ZooKeeperLockTestSuite{}, t)
}
