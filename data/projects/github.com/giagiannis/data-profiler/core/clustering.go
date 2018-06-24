package core

import (
	"bytes"
	"errors"
	"fmt"
)

// Clustering struct is responsible to execute the necessary actions in order
// to cluster the datasets based on their availability
type Clustering struct {
	datasets     []*Dataset               // list of datasets
	similarities *DatasetSimilarityMatrix // dataset similarities
	results      *Dendrogram              // holds the clustering results
	concurrency  int
}

// NewClustering is the the constructor for creating a Clustering object,
// providing a DatasetSimilarities object
func NewClustering(similarities *DatasetSimilarityMatrix, datasets []*Dataset) *Clustering {
	c := new(Clustering)
	c.similarities = similarities
	c.concurrency = 1
	c.datasets = datasets
	return c
}

// SetConcurrency sets the number of threads to be used
func (c *Clustering) SetConcurrency(concurrency int) {
	c.concurrency = concurrency
}

// Compute executes the clustering
func (c *Clustering) Compute() error {
	c.results = NewDendrogram(c.datasets)
	for c.results.hasUnmerged() {
		unmerged := c.results.getUnmerged()
		c.mergeClosestPairs(unmerged)
	}
	return nil
}

// Results returns the results
func (c *Clustering) Results() *Dendrogram {
	return c.results
}

// Returns the similarity between two different clusters of datasets
func (c *Clustering) getClustersSimilarity(a, b []*Dataset) float64 {
	sum := 0.0
	for i := range a {
		for j := range b {
			sum += c.similarities.Get(i, j)
		}
	}
	return sum / float64(len(a)*len(b))
}

// Evaluates and merges the most similar clusters at each turn
func (c *Clustering) mergeClosestPairs(unmerged []*DendrogramNode) {
	getDistanceForLine := func(i int) int {
		maxSimilarity, closestNode := -1.0, -1
		for j := 0; j < len(unmerged); j++ {
			n1, n2 := unmerged[i], unmerged[j]
			val := c.getClustersSimilarity(n1.datasets, n2.datasets)
			if (i != j) && (maxSimilarity < val) {
				maxSimilarity = val
				closestNode = j
			}
		}
		return closestNode
	}

	// parallel stuff
	type Pair struct {
		o, d int
	}
	resChannel := make(chan Pair)
	ch := make(chan bool, c.concurrency)
	for i := 0; i < c.concurrency; i++ {
		ch <- true
	}
	for i := 0; i < len(unmerged); i++ {
		go func(i int, done chan Pair, ch chan bool) {
			<-ch
			ret := getDistanceForLine(i)
			ch <- true
			done <- Pair{i, ret}
		}(i, resChannel, ch)
	}
	closestNode := make(map[int]int)
	for i := 0; i < len(unmerged); i++ {
		p := <-resChannel
		closestNode[p.o] = p.d
	}

	// serially merge all the unmerged nodes
	for k, v := range closestNode {
		if closestNode[v] == k {
			c.results.merge(unmerged[k], unmerged[v])
		}
	}

}

// Dendrogram represents the results of the ClusterApp objects
type Dendrogram struct {
	root     *DendrogramNode         // the root of the tree
	unmerged map[int]*DendrogramNode // Dendrogram leaf nodes, indexed by their id
}

// DendrogramNode is the node of the Dendrogram
type DendrogramNode struct {
	id          int
	datasets    []*Dataset
	father      *DendrogramNode
	left, right *DendrogramNode
}

func (n DendrogramNode) String() string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, d := range n.datasets {
		buf.WriteString(fmt.Sprintf("%s", (*d).Path()))
		if i < len(n.datasets)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString("}")
	return buf.String()
}

// NewDendrogram is the constructor for a Dendrogram struct
func NewDendrogram(datasets []*Dataset) *Dendrogram {
	d := new(Dendrogram)
	d.root = nil
	d.unmerged = make(map[int]*DendrogramNode)
	for i := 0; i < len(datasets); i++ {
		d.unmerged[i] = &DendrogramNode{i,
			[]*Dataset{datasets[i]},
			nil, nil, nil}
	}

	return d
}

// GetClusters function returns a slice containing the clusters of datasets
// for the specified dendrogram level
func (d *Dendrogram) GetClusters(level int) [][]*Dataset {
	var dfs func(*DendrogramNode, int) [][]*Dataset
	dfs = func(node *DendrogramNode, level int) [][]*Dataset {
		res := make([][]*Dataset, 0)
		if level > 0 && !node.isLeaf() {
			left := dfs(node.left, level-1)
			right := dfs(node.right, level-1)
			res = append(res, left...)
			res = append(res, right...)
			return res
		}
		res = append(res, node.datasets)
		return res
	}
	return dfs(d.root, level)
}

// Heights function returns the tree heights (max, min)
func (d *Dendrogram) Heights() (int, int) {
	var dfs func(*DendrogramNode) (int, int)
	dfs = func(node *DendrogramNode) (int, int) {
		if node.isLeaf() {
			return 0, 0
		}
		leftMax, leftMin := dfs(node.left)
		rightMax, rightMin := dfs(node.right)
		totalMax := leftMax
		if rightMax > leftMax {
			totalMax = rightMax
		}
		totalMin := leftMin
		if rightMin < leftMin {
			totalMin = rightMin
		}
		return totalMax + 1, totalMin + 1
	}
	return dfs(d.root)
}

func (d *Dendrogram) hasUnmerged() bool {
	return d.root == nil
}

// Returns a list of the unmerged nodes.
func (d *Dendrogram) getUnmerged() []*DendrogramNode {
	res := make([]*DendrogramNode, len(d.unmerged))
	i := 0
	for _, node := range d.unmerged {
		res[i] = node
		i++
	}
	return res
}

// Merge method is used to merge nodes of the tree that have not been merged yet
func (d *Dendrogram) merge(a, b *DendrogramNode) error {
	if _, ok := d.unmerged[a.id]; !ok {
		return errors.New("Node already merged or not known")
	} else if _, ok := d.unmerged[b.id]; !ok {
		return errors.New("Node already merged or not known")
	}
	delete(d.unmerged, a.id)
	delete(d.unmerged, b.id)
	newNode := new(DendrogramNode)
	newNode.id = a.id
	newNode.left, newNode.right, newNode.father = a, b, nil
	a.father, b.father = newNode, newNode
	newNode.datasets = append(newNode.datasets, a.datasets...)
	newNode.datasets = append(newNode.datasets, b.datasets...)
	if len(d.unmerged) == 0 { // the new on is the root
		d.root = newNode
	} else {
		d.unmerged[newNode.id] = newNode
	}
	return nil
}

func (d *Dendrogram) String() string {
	spa := "\t"
	var buf bytes.Buffer
	var dfs func(*DendrogramNode, string)
	dfs = func(node *DendrogramNode, pad string) {
		buf.WriteString(fmt.Sprintf("%s%s\n", pad, node.String()))
		if node.left != nil && node.right != nil {
			dfs(node.left, pad+spa)
			dfs(node.right, pad+spa)
		}
	}
	dfs(d.root, "")
	return buf.String()
}

func (n *DendrogramNode) isLeaf() bool {
	return n.left == nil || n.right == nil
}
