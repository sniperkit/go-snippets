package leadership

import (
	"github.com/nickbruun/gocommons/zkutils"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

// ZooKeeper leadership provider.
type zkLeadershipProvider struct {
	cm   *zkutils.ConnMan
	path string
	acl  []zk.ACL
}

// New ZooKeeper leadership provider.
func NewZooKeeperLeadershipProvider(connMan *zkutils.ConnMan, path string, acl []zk.ACL) LeadershipProvider {
	return &zkLeadershipProvider{
		cm:   connMan,
		path: path,
		acl:  acl,
	}
}

func (p *zkLeadershipProvider) GetCandidate(data []byte, leadershipHandler LeadershipHandler) Candidate {
	zc := &zkCandidate{
		cm:   p.cm,
		pp:   strings.TrimRight(p.path, "/"),
		acl:  p.acl,
		data: data,
		lh:   leadershipHandler,
		stop: make(chan struct{}, 10),
	}

	zc.done.Add(1)
	go zc.run()

	return zc
}
