package locking

import (
	"fmt"
	log "github.com/nickbruun/gocommons/logging"
	"github.com/nickbruun/gocommons/zkutils"
	"github.com/samuel/go-zookeeper/zk"
	"path"
	"sync"
	"time"
)

// ZooKeeper lock.
type zkLock struct {
	// Path.
	path string

	// ZooKeeper connection.
	cm *zkutils.ConnMan

	// ZooKeeper ACL.
	acl []zk.ACL

	// Lock mutex.
	lock sync.Mutex

	// Locked.
	locked bool

	// Release request channel.
	releaseCh chan struct{}

	// Release wait group.
	releaseWg sync.WaitGroup
}

// Test if cancellation has been requested.
func isCancelled(cancelCh <-chan interface{}) bool {
	select {
	case <-cancelCh:
		return true
	default:
		return false
	}
}

// Test if cancellation has been requested before a certain timeout.
func isCancelledBefore(cancelCh <-chan interface{}, timeout time.Duration) bool {
	select {
	case <-time.After(timeout):
		return false
	case <-cancelCh:
		return true
	}
}

// Get node path from name.
func (l *zkLock) nodePath(name string) string {
	return fmt.Sprintf("%s/%s", l.path, name)
}

// List all locks.
func (l *zkLock) listLockNodes() ([]zkutils.SequenceNode, error) {
	// List the child nodes of the path.
	children, _, err := l.cm.Conn.Children(l.path)
	if err == zk.ErrNoNode {
		return []zkutils.SequenceNode{}, nil
	} else if err != nil {
		return nil, err
	}

	return zkutils.ParseSequenceNodes(children, "lock"), nil
}

// Create lock node.
func (l *zkLock) createNode() (n zkutils.SequenceNode, err error) {
	// Attempt to create the lock node.
	lockPath := fmt.Sprintf("%s/lock", l.path)

	var path string
	path, err = l.cm.Conn.CreateProtectedEphemeralSequential(lockPath, nil, l.acl)

	if err == nil {
		n, err = zkutils.ParseSequenceNode(path, "lock")
		return
	} else if err == zk.ErrNoNode {
		if err = zkutils.CreateRecursively(l.cm.Conn, l.path, l.acl); err != nil {
			return
		}
		return l.createNode()
	} else {
		log.Warnf("Error creating epehermeral sequential node for lock: %v", err)
		return
	}
}

// Retain lock.
func (l *zkLock) retainLock(lockNode zkutils.SequenceNode) (<-chan struct{}, error) {
	lockPath := l.nodePath(lockNode.Name)
	log.Infof("Acquired lock: %s", l.path)

	// Watch for session or node loss.
	lost := make(chan struct{}, 1)
	lostEnd := make(chan struct{}, 1)

	go func() {
		sessionLoss := l.cm.WatchSessionLoss()
		exists, _, ec, err := l.cm.Conn.ExistsW(lockPath)

		if err != nil {
			log.Warnf("Error checking for existence of lock node %s, releasing lock: %v", lockPath, err)
			lost <- struct{}{}
			return
		} else if !exists {
			log.Warnf("Lock node no longer exists, releasing lock: %s", lockPath)
			lost <- struct{}{}
			return
		}

		select {
		case <-ec:
			log.Warnf("Lock node removed, releasing lock: %s", lockPath)
		case <-sessionLoss:
			log.Warnf("Session likely lost, releasing lock: %s", lockPath)
		case <-lostEnd:
		}

		lost <- struct{}{}
	}()

	// Set up the channels and wait for release or failure.
	failCh := make(chan struct{}, 1)
	l.releaseCh = make(chan struct{}, 1)
	l.releaseWg.Add(1)
	l.locked = true

	go func() {
		// Remove the lock node once we're done.
		defer zkutils.DeleteSafely(l.cm.Conn, lockPath)

		// Wait for the lock to be released or the lock node to disappear for
		// some reason.
		select {
		case <-l.releaseCh:
			// If the lock was released, stop watching the node.
			lostEnd <- struct{}{}
			log.Debugf("Released lock: %s", l.path)

		case <-lost:
			// If the lock node is removed (or presumed removed due to a
			// timeout), send a fail message and unlock the lock.
			log.Debug("Lock node removed or check failed, signalling failure")
			failCh <- struct{}{}
		}

		l.lock.Lock()
		l.locked = false
		l.releaseWg.Done()
		l.lock.Unlock()
	}()

	return failCh, nil
}

// Wait to acquire the lock.
func (l *zkLock) awaitLock(lockNode zkutils.SequenceNode, cancelCh <-chan interface{}) (<-chan struct{}, error) {
	for {
		// Test if acquisition has been cancelled.
		if isCancelled(cancelCh) {
			return nil, ErrCancelled
		}

		// List all lock nodes.
		lockNodes, err := l.listLockNodes()
		if err != nil {
			if !zkutils.IsErrorRecoverable(err) {
				log.Errorf("Unrecoverable ZooKeeper error occured: %v", err)
				return nil, err
			}

			log.Warn("Error retrieving leadership locks, waiting 100 ms to retry: %v", err)

			if isCancelledBefore(cancelCh, 100*time.Millisecond) {
				return nil, ErrCancelled
			} else {
				continue
			}
		}

		// Sort the list of lock nodes.
		zkutils.SortSequenceNodes(lockNodes)

		// Ensure that our lock node is present in the list of
		// locks.
		lockIdx := -1

		for i, n := range lockNodes {
			if n.SequenceNumber == lockNode.SequenceNumber {
				lockIdx = i
				break
			}
		}

		if lockIdx == -1 {
			log.Warn("Lock node has gone away")
			return nil, nil
		}

		// If the first node is our lock, we have obtained the lock.
		if lockIdx == 0 {
			return l.retainLock(lockNode)
		} else {
			// Wait for the lock node prior to our lock to disappear.
			path := l.nodePath(lockNodes[lockIdx-1].Name)
			log.Debugf("Waiting for change in %s", path)

			exists, _, eventChan, err := l.cm.Conn.ExistsW(path)
			if !exists || err != nil {
				if !zkutils.IsErrorRecoverable(err) {
					log.Errorf("Unrecoverable ZooKeeper error occured: %v", err)
					return nil, err
				}

				continue
			}

			select {
			case <-eventChan:
				break
			case <-cancelCh:
				return nil, ErrCancelled
			}
		}
	}
}

func (l *zkLock) LockWithCancel(cancelCh <-chan interface{}) (<-chan struct{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.locked {
		return nil, ErrLocked
	}

	for {
		// Test if acqusition has been stopped.
		if isCancelled(cancelCh) {
			return nil, ErrCancelled
		}

		// Create a lock node.
		lockNode, err := l.createNode()
		if err != nil {
			if !zkutils.IsErrorRecoverable(err) {
				log.Errorf("Unrecoverable ZooKeeper error occured: %v", err)
				return nil, err
			}

			log.Warnf("Error creating lock node, waiting 100 ms to retry: %v", err)

			if isCancelledBefore(cancelCh, 100*time.Millisecond) {
				return nil, ErrCancelled
			} else {
				continue
			}
		}

		_, lockNode.Name = path.Split(lockNode.Name)

		// Await acqusition of the lock.
		failCh, err := l.awaitLock(lockNode, cancelCh)
		if failCh == nil {
			defer zkutils.DeleteSafely(l.cm.Conn, l.nodePath(lockNode.Name))
		}

		if failCh != nil || err != nil {
			return failCh, err
		}
	}
}

func (l *zkLock) Lock() (<-chan struct{}, error) {
	return l.LockWithCancel(make(chan interface{}))
}

func (l *zkLock) LockTimeout(timeout time.Duration) (failCh <-chan struct{}, err error) {
	cancelCh := make(chan interface{}, 1)
	go func() {
		log.Debugf("Timing out %v after %s", l, timeout)
		<-time.After(timeout)
		log.Debugf("Timeout occured for %v", l)
		cancelCh <- struct{}{}
	}()

	failCh, err = l.LockWithCancel(cancelCh)
	if err == ErrCancelled {
		err = ErrTimeout
	}
	return
}

func (l *zkLock) Unlock() error {
	l.lock.Lock()

	if !l.locked {
		l.lock.Unlock()
		return ErrNotLocked
	}

	l.releaseCh <- struct{}{}
	l.lock.Unlock()

	l.releaseWg.Wait()
	return nil
}
