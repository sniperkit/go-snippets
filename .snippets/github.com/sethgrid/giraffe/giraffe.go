package giraffe

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// Error constants
const (
	// ErrKeyExists returns as the error value when an operation would cause a key collision
	ErrKeyExists = "key exists"

	// ErrCircular returns when adding a node relationship would result in a circular relationship
	ErrCircular = "circular relationship"
)

// Graph is a data structure composed of Nodes that have directional relationships to one another
type Graph struct {
	Name string

	duplicateKeys        bool
	circularRelationship bool

	sync.Mutex
	Nodes map[uint64]*Node
	keys  map[string]bool

	topNodeID uint64
}

// NewGraph creates a graph with default properties
func NewGraph(name string) (*Graph, error) {
	return NewConstraintGraph(name, true, true)
}

// NewConstraintGraph allows you to control if duplicate keys or circular relationships can exist
func NewConstraintGraph(name string, duplicateKeys, circularRelationship bool) (*Graph, error) {
	g := &Graph{
		Name:                 name,
		Nodes:                make(map[uint64]*Node),
		keys:                 make(map[string]bool),
		duplicateKeys:        duplicateKeys,
		circularRelationship: circularRelationship,
	}
	// insert root node
	g.Nodes[0] = &Node{ID: 0}

	return g, nil
}

// Root is a simple accessor function to get to the initial root node
func (g *Graph) Root() *Node {
	return g.Nodes[0]
}

// NodeCount returns the number of nodes in the graph
func (g *Graph) NodeCount() int {
	g.Lock()
	defer g.Unlock()
	return len(g.Nodes)
}

// LastNodeID returns the id of the last inserted node
func (g *Graph) LastNodeID() uint64 {
	g.Lock()
	g.Unlock()
	return g.topNodeID
}

// InsertNode inserts an empty default node into the graph. See InsertDataNode().
func (g *Graph) InsertNode() *Node {
	g.Lock()
	g.Unlock()

	n := g.insertNode()

	return n
}

// InsertDataNode is an alternate constructor to InsertNode() allowing you to pass in a key and value
func (g *Graph) InsertDataNode(key string, value []byte) (*Node, error) {
	g.Lock()
	g.Unlock()

	if !g.duplicateKeys {
		if _, ok := g.keys[key]; ok {
			return nil, errors.New(ErrKeyExists)
		}
		g.keys[key] = true
	}

	n := g.insertNode()

	n.Key = key
	n.Value = value
	return n, nil
}

// insertNode contains shared logic for the insert node calls
func (g *Graph) insertNode() *Node {
	n := &Node{
		ID:                   atomic.AddUint64(&g.topNodeID, 1),
		circularRelationship: g.circularRelationship,
	}

	g.Nodes[n.ID] = n

	return n
}

// DeleteNode removes a node and its relationships from the graph
func (g *Graph) DeleteNode(node *Node) error {
	g.Lock()
	defer g.Unlock()

	return g.deleteNodeByID(node.ID)
}

// DeleteNodeByID removes a node (by its ID) and its relationships from the graph
func (g *Graph) DeleteNodeByID(ID uint64) error {
	g.Lock()
	defer g.Unlock()

	return g.deleteNodeByID(ID)
}

// deleteNodeByID is a helper function to keep logic DRY in the delete node endpoints
func (g *Graph) deleteNodeByID(ID uint64) error {
	sourceIDs := extractIDs(g.Nodes[ID].sources)
	destinationIDs := extractIDs(g.Nodes[ID].destinations)

	for _, nodeID := range sourceIDs {
		g.Nodes[nodeID].RemoveRelationship(g.Nodes[ID])
	}
	for _, nodeID := range destinationIDs {
		g.Nodes[ID].RemoveRelationship(g.Nodes[nodeID])
	}

	delete(g.Nodes, ID)

	return nil
}

// FindRoots finds all nodes that do not have a source nodes below them
func (g *Graph) FindRoots() []uint64 {
	g.Lock()
	defer g.Unlock()
	var roots []uint64
	for _, node := range g.Nodes {
		node.Lock()
		if len(node.sources) == 0 {
			roots = append(roots, node.ID)
		}
		node.Unlock()
	}
	return roots
}

// FindNodeIDByKey returns the first node's ID with a matching key
func (g *Graph) FindNodeIDByKey(key string) (uint64, bool) {
	for id, node := range g.Nodes {
		if node.Key == key {
			return id, true
		}
	}

	return 0, false
}

// FindNodeByKey returns the the first node with a matching key
func (g *Graph) FindNodeByKey(key string) (*Node, bool) {
	id, found := g.FindNodeIDByKey(key)
	if !found {
		return nil, false
	}

	return g.Nodes[id], true
}

// ToVisJS generates a simple HTML/Javascript view of the data
// See http://visjs.org/ for information and styles
func (g *Graph) ToVisJS(showID, showKey, showValue bool) string {
	dataSet := ""
	edges := ""

	for id, node := range g.Nodes {
		label := ""
		if showID {
			label += fmt.Sprintf("node %d ", id)
		}
		if showKey {
			label += fmt.Sprintf("%s ", node.Key)
		}
		if showValue {
			label += fmt.Sprintf("%s ", string(node.Value))
		}

		dataSet += fmt.Sprintf(`{id: %d, label: '%s'},`, id, label)
		for _, nID := range extractIDs(node.ListDestinations()) {
			edges += fmt.Sprintf(`{from: %d, to: %d, arrows:'middle',},`, id, nID)
		}
	}

	js := `
<html>
<head>
    <script type="text/javascript" src="http://visjs.org/dist/vis.js"></script>
    <link href="http://visjs.org/dist/vis.css" rel="stylesheet" type="text/css" />

    <style type="text/css">
        #mynetwork {
            width: 600px;
            height: 400px;
            border: 1px solid lightgray;
        }
    </style>
</head>
<body>
<div id="mynetwork"></div>

<script type="text/javascript">
    // create an array with nodes
    var nodes = new vis.DataSet([
        ` + dataSet + `
    ]);

    // create an array with edges
    var edges = new vis.DataSet([
        ` + edges + `
    ]);

    // create a network
    var container = document.getElementById('mynetwork');

    // provide the data in the vis format
    var data = {
        nodes: nodes,
        edges: edges
    };
    var options = {};

    // initialize your network!
    var network = new vis.Network(container, data, options);
</script>
</body>
</html>
    `

	return js
}
