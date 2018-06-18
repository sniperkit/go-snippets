package leadership

import (
	"fmt"
	log "github.com/nickbruun/gocommons/logging"
	"github.com/nickbruun/gocommons/zkutils"
	"github.com/samuel/go-zookeeper/zk"
	"path"
	"sync"
	"time"
)

// ZooKeeper candidate.
type zkCandidate struct {
	// ZooKeeper connection.
	cm *zkutils.ConnMan

	// ZooKeeper path prefix.
	pp string

	// ZooKeeper ACL.
	acl []zk.ACL

	// Leadership handler.
	lh LeadershipHandler

	// Stop channel.
	stop chan struct{}

	// Done.
	done sync.WaitGroup

	// Leader data.
	data []byte
}

// Test if the candidature has been stopped.
func (c *zkCandidate) isStopped() bool {
	select {
	case <-c.stop:
		return true
	default:
		return false
	}
}

// Test if the candidature has been stopped before a certain amount of time.
func (c *zkCandidate) isStoppedBefore(d time.Duration) bool {
	select {
	case <-c.stop:
		return true
	case <-time.After(d):
		return false
	}
}

// Get candidate node path from name.
func (c *zkCandidate) nodePath(name string) string {
	return fmt.Sprintf("%s/%s", c.pp, name)
}

// List all candidates.
func (c *zkCandidate) listCandidateNodes() ([]zkutils.SequenceNode, error) {
	// List the child nodes of the prefix.
	children, _, err := c.cm.Conn.Children(c.pp)
	if err == zk.ErrNoNode {
		return []zkutils.SequenceNode{}, nil
	} else if err != nil {
		return nil, err
	}

	return zkutils.ParseSequenceNodes(children, "candidate"), nil
}

// Create candidate node.
func (c *zkCandidate) createCandidateNode() (n zkutils.SequenceNode, err error) {
	// Attempt to create the candidate node.
	candidatePath := fmt.Sprintf("%s/candidate", c.pp)

	var path string
	path, err = c.cm.Conn.CreateProtectedEphemeralSequential(candidatePath, c.data, c.acl)

	if err == nil {
		n, err = zkutils.ParseSequenceNode(path, "candidate")
		return
	} else if err == zk.ErrNoNode {
		if err = zkutils.CreateRecursively(c.cm.Conn, c.pp, c.acl); err != nil {
			return
		}
		return c.createCandidateNode()
	} else {
		log.Warnf("Error creating epehermeral sequential node for candidate: %v", err)
		return
	}
}

// Assume leadership.
func (c *zkCandidate) assumeLeadership(candidateNode zkutils.SequenceNode) (stopped bool) {
	// Remove the candidate node once we're done.
	candidatePath := c.nodePath(candidateNode.Name)
	defer zkutils.DeleteSafely(c.cm.Conn, candidatePath)

	log.Info("Became leader")

	// Invoke the leadership handler.
	end := make(chan struct{}, 1)
	resigned := make(chan struct{}, 1)

	go func() {
		c.lh(end)
		resigned <- struct{}{}
	}()

	// Watch for session or node loss.
	lost := make(chan struct{}, 1)
	lostEnd := make(chan struct{}, 1)

	go func() {
		sessionLoss := c.cm.WatchSessionLoss()
		exists, _, ec, err := c.cm.Conn.ExistsW(candidatePath)

		if err != nil {
			log.Warnf("Error checking for existence of leader node %s, resigning as leader: %v", candidatePath, err)
			lost <- struct{}{}
			return
		} else if !exists {
			log.Warnf("Candidate node no longer exists, resigning as leader: %s", candidatePath)
			lost <- struct{}{}
			return
		}

		select {
		case <-ec:
			log.Warnf("Candidate node removed, resigning as leader: %s", candidatePath)
		case <-sessionLoss:
			log.Warnf("Session likely lost, resigning as leader: %s", candidatePath)
		case <-lostEnd:
		}

		lost <- struct{}{}
	}()

	// Wait for the leadership handler to stop, the candidate to be stopped
	// externally or the candidate node to disappear for some reason.
	stopped = false

	select {
	case <-resigned:
		lostEnd <- struct{}{}
		break

	case <-c.stop:
		// If we are stopped, send an end message to the leadership handler
		// and wait for the handler to resign.
		stopped = true
		lostEnd <- struct{}{}
		end <- struct{}{}
		log.Debug("Candidate stopped, awaiting leadership handler return")
		<-resigned

	case <-lost:
		// If the candidate node is removed (or presumed removed due to a
		// timeout), send an end message to the leadership handler and wait for
		// the handler to resign.
		end <- struct{}{}
		log.Debug("Candidate node removed or check failed, awaiting leadership handler return")
		<-resigned
	}

	log.Info("Resigned leadership")

	return stopped || c.isStopped()
}

// Wait to become the leader.
func (c *zkCandidate) awaitLeadership(candidateNode zkutils.SequenceNode) (stopped bool) {
	for {
		// Test if the candidate has been stopped.
		if c.isStopped() {
			return true
		}

		// List all candidate nodes.
		candidateNodes, err := c.listCandidateNodes()
		if err != nil {
			log.Warn("Error retrieving leadership candidates, waiting 100 ms to retry: %v", err)

			if c.isStoppedBefore(100 * time.Millisecond) {
				return true
			} else {
				continue
			}
		}

		// Sort the list of candidate nodes.
		zkutils.SortSequenceNodes(candidateNodes)

		// Ensure that our candidate node is present in the list of
		// candidates.
		candidateIdx := -1

		for i, n := range candidateNodes {
			if n.SequenceNumber == candidateNode.SequenceNumber {
				candidateIdx = i
				break
			}
		}

		if candidateIdx == -1 {
			log.Warn("Candidate node has gone away")
			return false
		}

		// If the first candidate is our candidate, we have obtained
		// leadership.
		if candidateIdx == 0 {
			return c.assumeLeadership(candidateNode)
		} else {
			// Wait for the candidate node prior to our candidate to disappear.
			path := c.nodePath(candidateNodes[candidateIdx-1].Name)
			log.Debugf("Waiting for change in %s", path)

			exists, _, eventChan, err := c.cm.Conn.ExistsW(path)
			if !exists || err != nil {
				continue
			}

			select {
			case <-eventChan:
				break
			case <-c.stop:
				return true
			}
		}
	}
}

// Run the election process until stopped.
func (c *zkCandidate) run() {
	for {
		// Test if the candidate has been stopped.
		if c.isStopped() {
			break
		}

		// Create a candidate node.
		candidateNode, err := c.createCandidateNode()
		if err != nil {
			log.Warnf("Error creating candidate node, waiting 100 ms to retry: %v", err)
			if c.isStoppedBefore(100 * time.Millisecond) {
				break
			} else {
				continue
			}
		}

		_, candidateNode.Name = path.Split(candidateNode.Name)

		// Await leadership.
		if stopped := c.awaitLeadership(candidateNode); stopped {
			break
		}
	}

	c.done.Done()
	log.Debug("Done running candidate")
}

func (c *zkCandidate) Stop() {
	c.stop <- struct{}{}
	log.Debug("Sent stop signal, waiting for running candidate")
	c.done.Wait()
	log.Debug("Done waiting for running candidate")
}
