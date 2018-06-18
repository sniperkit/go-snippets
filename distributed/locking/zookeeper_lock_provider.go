package locking

import (
	"github.com/nickbruun/gocommons/zkutils"
	"github.com/samuel/go-zookeeper/zk"
)

// ZooKeeper lock provider.
type zkLockProvider struct {
	// ZooKeeper connection.
	cm *zkutils.ConnMan

	// ZooKeeper ACL.
	acl []zk.ACL
}

// New ZooKeeper lock provider.
func NewZooKeeperLockProvider(connMan *zkutils.ConnMan, acl []zk.ACL) LockProvider {
	return &zkLockProvider{
		cm:  connMan,
		acl: acl,
	}
}

func (p *zkLockProvider) GetLock(path string) Lock {
	return &zkLock{
		path: path,
		cm:   p.cm,
		acl:  p.acl,
	}
}
