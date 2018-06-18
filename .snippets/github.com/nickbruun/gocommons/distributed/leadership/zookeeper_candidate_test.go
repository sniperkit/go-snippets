package leadership

import (
	"fmt"
	logrus "github.com/Sirupsen/logrus"
	log "github.com/nickbruun/gocommons/logging"
	"github.com/nickbruun/gocommons/zkutils"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
	"testing"
	"time"
)

// ZooKeeper candidate test rig configuration.
type zkCandidateTestRigConfig struct {
	// ZooKeeper test cluster.
	TestCluster *zk.TestCluster

	// ZooKeeper test cluster server addresses.
	Servers []string

	// Candidate count.
	NumCandidates int

	// Create candidates serially.
	CreateCandidatesSerially bool

	// Path prefix.
	Path string
}

// ZooKeeper candidate test rig candidate.
type zkCandidateTestRigCandidate struct {
	zk         *zkutils.ConnMan
	cand       Candidate
	resign     chan struct{}
	isLeader   bool
	leaderTerm int
}

// ZooKeeper candidate test rig.
type zkCandidateTestRig struct {
	// Test case.
	t *testing.T

	// ZooKeeper test cluster.
	cluster *zk.TestCluster

	// ZooKeeper test cluster server addresses.
	servers []string

	// Internal ZooKeeper connection.
	zk *zk.Conn

	// Path prefix.
	path string

	// Candidate lock.
	candsLock sync.Mutex

	// Candidates.
	cands []*zkCandidateTestRigCandidate

	// Leader change.
	leaderChange chan struct{}
}

// Clean up the ZooKeeper candidate test rig.
func (r *zkCandidateTestRig) CleanUp() {
	if r.cands != nil {
		for _, c := range r.cands {
			c.cand.Stop()
			c.zk.Close()
		}

		r.cands = nil
	}

	if r.zk != nil {
		r.zk.Close()
		r.zk = nil
	}
}

// List candidate nodes.
func (r *zkCandidateTestRig) listCandidateNodes() []zkutils.SequenceNode {
	// List the child nodes of the prefix.
	children, _, err := r.zk.Children(r.path)
	if err == zk.ErrNoNode {
		return []zkutils.SequenceNode{}
	} else if err != nil {
		r.t.Fatalf("Failed to list candidate nodes: %v", err)
	}

	return zkutils.ParseSequenceNodes(children, "candidate")
}

// Await a leader.
func (r *zkCandidateTestRig) AwaitLeader() {
	<-r.leaderChange
}

// Get the current leader index.
//
// Returns -1 if no candidate is currently the leader.
func (r *zkCandidateTestRig) CurrentLeader() (candidate int, term int) {
	r.candsLock.Lock()
	defer r.candsLock.Unlock()

	for i, c := range r.cands {
		if c.isLeader {
			return i, c.leaderTerm
		}
	}

	return -1, -1
}

// Resign the current leader.
func (r *zkCandidateTestRig) ResignCurrentLeader() {
	r.candsLock.Lock()
	defer r.candsLock.Unlock()

	for _, c := range r.cands {
		if c.isLeader {
			c.resign <- struct{}{}
			return
		}
	}
}

// Await a new candidate node.
func (r *zkCandidateTestRig) AwaitNewCandidateNode() <-chan error {
	await := make(chan error, 1)

	go func() {
		var seqNumGt int32
		first := true

		for {
			// List the candidates.
			children, _, event, err := r.zk.ChildrenW(r.path)
			var nodes []zkutils.SequenceNode

			if err == zk.ErrNoNode {
				log.Debugf("Awaiting existence of %s", r.path)
				if err = <-zkutils.AwaitExists(r.zk, r.path); err != nil {
					await <- err
					return
				}
				log.Debugf("%s exists, retrying", r.path)

				seqNumGt = -1
				first = false

				continue
			} else if err != nil {
				await <- err
				return
			} else {
				nodes = zkutils.ParseSequenceNodes(children, "candidate")
			}

			// Determine the sequence number to wait for if the is the first
			// shot.
			if first {
				if len(nodes) > 0 {
					seqNumGt = nodes[len(nodes)-1].SequenceNumber
				} else {
					seqNumGt = -1
				}

				first = false

				log.Debugf("Awaiting candidate node with sequence number greater than %d", seqNumGt)
			} else {
				for _, n := range nodes {
					if n.SequenceNumber > seqNumGt {
						await <- nil
						return
					}
				}
			}

			log.Debugf("Awaiting event for children of %s", r.path)
			<-event
			log.Debugf("Event occured for children of %s", r.path)
		}
	}()

	return await
}

// Create a ZooKeeper candidate test rig.
func newZkCandidateTestRig(t *testing.T, config zkCandidateTestRigConfig) *zkCandidateTestRig {
	// Create an internal connection.
	zkConn, _, err := zk.Connect(config.Servers, 10*time.Second)
	if err != nil {
		t.Fatalf("Failed to create internal rig connection: %v", err)
	}

	// Create the rig.
	rig := &zkCandidateTestRig{
		t:            t,
		cluster:      config.TestCluster,
		servers:      config.Servers,
		zk:           zkConn,
		path:         config.Path,
		cands:        make([]*zkCandidateTestRigCandidate, 0, config.NumCandidates),
		leaderChange: make(chan struct{}, 128),
	}

	// Create the candidates, waiting for a new node to be created before
	// creating the next, if requested.
	rig.candsLock.Lock()

	for i := 0; i < config.NumCandidates; i++ {
		var await <-chan error

		if config.CreateCandidatesSerially {
			await = rig.AwaitNewCandidateNode()
		}

		candZkConn, err := zkutils.Connect(config.Servers, 10*time.Second)
		if err != nil {
			t.Fatalf("Failed to create candidate rig connection: %v", err)
		}

		candIdx := i

		provider := NewZooKeeperLeadershipProvider(candZkConn, rig.path, zk.WorldACL(zk.PermAll))
		cand := provider.GetCandidate([]byte(fmt.Sprintf("%d", i)), func(endChan <-chan struct{}) {
			// Mark the candidate as the leader.
			rig.candsLock.Lock()
			rig.leaderChange <- struct{}{}
			rig.cands[candIdx].isLeader = true
			rig.cands[candIdx].leaderTerm++
			rig.candsLock.Unlock()

			// Wait for the candidate to be stopped.
			select {
			case <-endChan:
				break
			case <-rig.cands[candIdx].resign:
				break
			}

			// Mark the candidate as no longer being the leader.
			rig.candsLock.Lock()
			rig.cands[candIdx].isLeader = false
			rig.candsLock.Unlock()
		})

		rig.cands = append(rig.cands, &zkCandidateTestRigCandidate{
			zk:     candZkConn,
			cand:   cand,
			resign: make(chan struct{}, 10),
		})

		if config.CreateCandidatesSerially {
			if err = <-await; err != nil {
				t.Fatalf("Failed to await creation of candidate: %v", err)
			}
		}
	}

	rig.candsLock.Unlock()

	return rig
}

func TestZooKeeperCandidate(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	// Set up the test cluster.
	testCluster, servers := zkutils.CreateTestCluster(t, 1)
	defer testCluster.Stop()

	// Ensure that with a single candidate and a fresh cluster, the candidate
	// becomes the leader. We do this twice to ensure repeatability.
	for i := 0; i < 2; i++ {
		rig := newZkCandidateTestRig(t, zkCandidateTestRigConfig{
			TestCluster:              testCluster,
			Servers:                  servers,
			Path:                     "/my/test/election",
			NumCandidates:            1,
			CreateCandidatesSerially: true,
		})

		// Ensure that the candidate becomes the leader.
		rig.AwaitLeader()

		leader, term := rig.CurrentLeader()
		if leader != 0 {
			t.Errorf("Expected candidate to be the current leader")
		}

		if term != 1 {
			t.Errorf("Expected candidate term to be 1, but is %d", term)
		}

		// Resign the candidate.
		rig.ResignCurrentLeader()

		// Ensure that the candidate becomes the leader once again.
		rig.AwaitLeader()

		leader, term = rig.CurrentLeader()
		if leader != 0 {
			t.Errorf("Expected candidate to be the current leader")
		}

		if term != 2 {
			t.Errorf("Expected candidate term to be 2, but is %d", term)
		}

		rig.CleanUp()

		// Ensure that there is no leader.
		leader, term = rig.CurrentLeader()
		if leader != -1 {
			t.Errorf("Expected no leader, but %d is leader", leader)
		}
	}

	// Ensure that with three candidates, the candidates become leaders in
	// sequential order.
	{
		rig := newZkCandidateTestRig(t, zkCandidateTestRigConfig{
			TestCluster:              testCluster,
			Servers:                  servers,
			Path:                     "/my/test/election",
			NumCandidates:            3,
			CreateCandidatesSerially: true,
		})

		for term := 0; term < 15; term++ {
			// Wait for a leader to be elected.
			rig.AwaitLeader()

			leader, leaderTerm := rig.CurrentLeader()
			if leader != term%3 {
				t.Errorf("Expected candidate %d to be the current leader, but %d is", term%3, leader)
			}

			if leaderTerm != 1+term/3 {
				t.Errorf("Expected candidate term to be %d, but is %d", 1+term/3, leaderTerm)
			}

			// Resign the leader.
			rig.ResignCurrentLeader()
		}

		// Clean up and ensure that there is no leader.
		rig.CleanUp()

		leader, _ := rig.CurrentLeader()
		if leader != -1 {
			t.Errorf("Expected no leader, but %d is leader", leader)
		}
	}
}
