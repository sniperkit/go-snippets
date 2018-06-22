package giraffe

import (
	"errors"
	"sync"
)

// Node is a key value pair that has directional relationships with other Nodes
type Node struct {
	ID    uint64
	Key   string
	Value []byte

	circularRelationship bool

	sync.RWMutex
	destinations []*Node
	sources      []*Node

	// used for encoding/decoding
	sourceIDs []uint64
}

// AddRelationship adds a newNode as a destination of this node
func (n *Node) AddRelationship(newNode *Node) error {
	n.Lock()
	defer n.Unlock()

	if !n.circularRelationship && newNode.DepthFirstSearch(n) {
		return errors.New(ErrCircular)
	}

	n.destinations = append(n.destinations, newNode)
	newNode.addSource(n)

	return nil
}

// addSource provides a way to more easily traverse the graph.
// addSource is not exposed to prevent circular logic as nodes are added
// to both destination and source lists
func (n *Node) addSource(newNode *Node) {
	n.Lock()
	defer n.Unlock()

	n.sources = append(n.sources, newNode)
}

// RemoveRelationship removes the edge/relationship between a source node and its destination node
func (n *Node) RemoveRelationship(oldNode *Node) error {
	n.Lock()
	defer n.Unlock()

	// while typical use would dictate that any node would only have one
	// relationship to another node, we cannot be sure. remove all relationships.
	var removeSources []*Node
	for _, n := range n.sources {
		if n.ID == oldNode.ID {
			removeSources = append(removeSources, n)
		}
	}
	var removeDestinations []*Node
	for _, n := range n.destinations {
		if n.ID == oldNode.ID {
			removeDestinations = append(removeDestinations, n)
		}
	}

	n.sources = difference(n.sources, removeSources)
	n.destinations = difference(n.destinations, removeDestinations)

	return nil
}

// ListDestinations lists all nodes that this node points towards
func (n *Node) ListDestinations() []*Node {
	n.Lock()
	defer n.Unlock()

	return n.destinations
}

// ListSources lists all nodes that point to this node
func (n *Node) ListSources() []*Node {
	n.Lock()
	defer n.Unlock()

	return n.sources
}

// DepthFirstSearch traverses the graph starting at this node to find the otherNode
func (n *Node) DepthFirstSearch(otherNode *Node) bool {
	for _, node := range n.destinations {
		if node.ID == otherNode.ID || node.DepthFirstSearch(otherNode) {
			return true
		}
	}
	return false
}

// BreadthFirstSearch traverses the graph starting at this node to find the otherNode
func (n *Node) BreadthFirstSearch(otherNode *Node) bool {
	// initialize the bfs queue
	return bfs(otherNode, n.destinations)
}

// bfs maintains the search queue and is the logic for BreadthFirstSearch()
func bfs(otherNode *Node, queue []*Node) bool {
	if len(queue) == 0 {
		return false
	}

	var nextQueue []*Node
	for _, node := range queue {
		if node.ID == otherNode.ID {
			return true
		}
		nextQueue = append(nextQueue, node.destinations...)
	}

	if bfs(otherNode, nextQueue) {
		return true
	}
	return false
}

// extractIDs grabs all the ids for the nodes. useful for debugging and in tests.
func extractIDs(nodes []*Node) []uint64 {
	IDs := make([]uint64, len(nodes))
	for i, node := range nodes {
		IDs[i] = node.ID
	}
	return IDs
}

// difference returns the result of removing subtrahend elements from the minuend
func difference(minuend, subtrahend []*Node) []*Node {
	var difference []*Node
	for _, m := range minuend {
		if !inList(subtrahend, m) {
			difference = append(difference, m)
		}
	}
	return difference
}

// inList returns true if the target node is in the node list
func inList(list []*Node, target *Node) bool {
	for _, node := range list {
		if node.ID == target.ID {
			return true
		}
	}
	return false
}
