package core

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// BhattacharyyaEstimator is the similarity estimator that quantifies the similarity
// of the distribution between the datasets.
type BhattacharyyaEstimator struct {
	AbstractDatasetSimilarityEstimator

	// struct used to map dataset paths to indexes
	inverseIndex map[string]int
	// determines the height of the kd tree to be used
	maxPartitions int
	// hold the portion of the data examined for constructing the tree
	kdTreeSamplePerc float64
	// kd tree, utilized for dataset partitioning
	kdTree *kdTreeNode
	// holds the number of points for each dataset region
	pointsPerRegion [][]int
	// holds the total number of points for each dataset
	datasetsSize []int
	// columns to investigate for partitioning
	columns []int
}

// Compute method constructs the Similarity Matrix
func (e *BhattacharyyaEstimator) Compute() error {
	return datasetSimilarityEstimatorCompute(e)
}

// Similarity returns the similarity between two datasets
func (e *BhattacharyyaEstimator) Similarity(a, b *Dataset) float64 {
	var indexA, indexB []int
	var countA, countB int
	if idx, ok := e.inverseIndex[a.Path()]; ok {
		indexA = e.pointsPerRegion[idx]
		countA = e.datasetsSize[idx]
	} else {
		err := a.ReadFromFile()
		if err != nil {
			log.Println(err)
		}
		indexA = e.kdTree.GetLeafIndex(a.Data())
		countA = len(a.Data())
	}

	if idx, ok := e.inverseIndex[b.Path()]; ok {
		indexB = e.pointsPerRegion[idx]
		countB = e.datasetsSize[idx]
	} else {
		err := b.ReadFromFile()
		if err != nil {
			log.Println(err)
		}
		indexB = e.kdTree.GetLeafIndex(b.Data())
		countB = len(b.Data())
	}
	return e.getValue(indexA, indexB, countA, countB)
}

// Configure sets a the configuration parameters of the estimator
func (e *BhattacharyyaEstimator) Configure(conf map[string]string) {
	if val, ok := conf["concurrency"]; ok {
		conv, err := strconv.ParseInt(val, 10, 32)
		e.concurrency = int(conv)
		if err != nil {
			log.Println(err)
		}
	} else {
		e.concurrency = 1
	}
	if val, ok := conf["partitions"]; ok {
		//conv, err := strconv.ParseInt(val, 10, 32)
		conv, err := strconv.ParseInt(val, 10, 32)
		e.maxPartitions = int(conv)
		if err != nil {
			log.Println(err)
		}
	} else {
		e.maxPartitions = 32
	}

	if val, ok := conf["tree.sr"]; ok {
		//conv, err := strconv.ParseInt(val, 10, 32)
		conv, err := strconv.ParseFloat(val, 64)
		e.kdTreeSamplePerc = conv
		if err != nil {
			log.Println(err)
		}
	} else {
		e.kdTreeSamplePerc = 0.1
	}

	if val, ok := conf["columns"]; ok && val != "all" {
		arr := strings.Split(val, ",")
		for _, d := range arr {
			v, err := strconv.ParseInt(d, 10, 32)
			if err != nil {
				log.Println(err)
			} else {
				e.columns = append(e.columns, int(v))
			}
		}
	} else {
		e.columns = nil
	}

	e.init()

}

func (e *BhattacharyyaEstimator) init() {
	// initialization step
	e.inverseIndex = make(map[string]int)
	for i, d := range e.datasets {
		e.inverseIndex[d.Path()] = i
	}
	for _, d := range e.datasets {
		d.ReadFromFile()
	}
	log.Println("Estimating a KD-tree partition")
	//e.kdTree = newKDTreePartition(e.datasets[0].Data())
	s := e.sampledDataset()
	if e.columns == nil {
		log.Println("Setting all columns to examine for the kd-tree")
		for i := 0; i < len(s[0].Data); i++ {
			e.columns = append(e.columns, i)
		}
	}
	e.kdTree = newKDTreePartition(s, e.columns)
	temp, height := e.maxPartitions, 0
	for temp > 0 {
		temp = temp >> 1
		height++
	}
	e.kdTree.Prune(height)
	log.Printf("Prunning the tree to height %d created %d partitions",
		height,
		len(e.kdTree.Leaves()))
	e.pointsPerRegion = make([][]int, len(e.datasets))
	e.datasetsSize = make([]int, len(e.datasets))
	for i, d := range e.datasets {
		e.pointsPerRegion[i] = e.kdTree.GetLeafIndex(d.Data())
		e.datasetsSize[i] = len(d.Data())
	}
}

// Options returns a list of parameters that can be set by the user
func (e *BhattacharyyaEstimator) Options() map[string]string {
	return map[string]string{
		"concurrency": "max num of threads used (int)",
		"partitions":  "max number of partitions to be used for the estimation (default is 32)",
		"tree.sr":     "determines the portion of datasets to sample for the kd tree construction",
		"columns":     "comma separated values of column indices to consider (starting from 0)  or all (default)",
		//"tree.scale":  "determines the portion of the kd-tree to be used",
	}
}

func (e *BhattacharyyaEstimator) getValue(indA, indB []int, countA, countB int) float64 {
	sum := 0.0
	for k := 0; k < len(indA); k++ {
		sum += math.Sqrt(float64(indA[k] * indB[k]))
	}
	sum /= math.Sqrt(float64(countA * countB))
	return sum
}

// Serialize returns a byte array containing a serialized form of the estimator
func (e *BhattacharyyaEstimator) Serialize() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(getBytesInt(int(SimilarityTypeBhattacharyya)))

	bytes := datasetSimilarityEstimatorSerialize(e.AbstractDatasetSimilarityEstimator)
	buffer.Write(bytes)
	buffer.Write(getBytesInt(e.maxPartitions))

	// write points per region
	buffer.Write(getBytesInt(len(e.pointsPerRegion[0])))
	for i := range e.pointsPerRegion {
		for j := range e.pointsPerRegion[i] {
			buffer.Write(getBytesInt(e.pointsPerRegion[i][j]))
		}
	}
	// write datasets size
	for _, s := range e.datasetsSize {
		buffer.Write(getBytesInt(s))
	}

	// write kdtree
	serializedTree := e.kdTree.Serialize()
	buffer.Write(getBytesInt(len(serializedTree)))
	buffer.Write(e.kdTree.Serialize())
	return buffer.Bytes()
}

// Deserialize constructs a similarity object based on the byte stream
func (e *BhattacharyyaEstimator) Deserialize(b []byte) {
	buffer := bytes.NewBuffer(b)
	tempInt := make([]byte, 4)
	tempFloat := make([]byte, 8)
	buffer.Read(tempInt) // contains estimator type

	var count int
	buffer.Read(tempInt)
	absEstBytes := make([]byte, getIntBytes(tempInt))
	buffer.Read(absEstBytes)
	e.AbstractDatasetSimilarityEstimator =
		*datasetSimilarityEstimatorDeserialize(absEstBytes)

	buffer.Read(tempInt)
	e.maxPartitions = getIntBytes(tempFloat)

	e.inverseIndex = make(map[string]int)
	for i := range e.datasets {
		e.inverseIndex[e.datasets[i].Path()] = i
	}

	buffer.Read(tempInt)
	count = getIntBytes(tempInt)
	e.pointsPerRegion = make([][]int, len(e.datasets))
	for i := range e.pointsPerRegion {
		e.pointsPerRegion[i] = make([]int, count)
		for j := range e.pointsPerRegion[i] {
			buffer.Read(tempInt)
			e.pointsPerRegion[i][j] = getIntBytes(tempInt)
		}
	}

	e.datasetsSize = make([]int, len(e.datasets))
	for i := range e.datasetsSize {
		buffer.Read(tempInt)
		e.datasetsSize[i] = getIntBytes(tempInt)
	}

	buffer.Read(tempInt)
	count = getIntBytes(tempInt)
	tempCustom := make([]byte, count)
	buffer.Read(tempCustom)
	e.kdTree = new(kdTreeNode)
	e.kdTree.Deserialize(tempCustom)

}

// sampledDataset returns a custom dataset that consist of the tuples of the
// previous
func (e *BhattacharyyaEstimator) sampledDataset() []DatasetTuple {
	log.Println("Generating a sampled and merged dataset with all tuples")
	var result []DatasetTuple
	for _, d := range e.datasets {
		tuplesToChoose := int(math.Floor(float64(len(d.Data())) * e.kdTreeSamplePerc))
		log.Printf("%d/%d tuples chosen for %s\n", tuplesToChoose, len(d.Data()), d.path)
		tuplesIdx := make(map[int]bool)
		for len(tuplesIdx) < tuplesToChoose {
			tuplesIdx[rand.Int()%len(d.Data())] = true
		}
		for k := range tuplesIdx {
			result = append(result, d.Data()[k])
		}
	}
	return result
}

type kdTreeNode struct {
	dim   int
	value float64
	right *kdTreeNode
	left  *kdTreeNode
}

func (r *kdTreeNode) Serialize() []byte {
	var countNodes func(*kdTreeNode) int
	countNodes = func(node *kdTreeNode) int {
		if node == nil {
			return 0
		}
		count := 0
		count += countNodes(node.left)
		count += countNodes(node.right)
		return count + 1
	}
	dataSize := 8 + 4
	//nodes := countNodes(r)
	nodes := 2<<uint(r.Height()-1) - 1
	buf := make([]byte, dataSize*nodes)
	valid := make([]byte, nodes)
	var serializeNode func(*kdTreeNode, int)
	serializeNode = func(node *kdTreeNode, offset int) {
		if node == nil {
			if offset < len(valid) {
				valid[offset] = '0'
			}
			return
		}
		valid[offset] = '1'
		for i, b := range getBytesInt(node.dim) {
			buf[offset*dataSize+i] = b
		}
		count := len(getBytesInt(node.dim))
		for i, b := range getBytesFloat(node.value) {
			buf[offset*dataSize+count+i] = b
		}
		serializeNode(node.left, (2*offset + 1))
		serializeNode(node.right, (2*offset + 2))
	}
	serializeNode(r, 0)
	buffer := new(bytes.Buffer)
	buffer.Write(getBytesInt(len(buf)))
	buffer.Write(buf)
	buffer.Write(getBytesInt(len(valid)))
	buffer.Write(valid)
	return buffer.Bytes()
}

func (r *kdTreeNode) Deserialize(b []byte) {
	buf := bytes.NewBuffer(b)
	tempIntBuff := make([]byte, 4)
	buf.Read(tempIntBuff)
	count := getIntBytes(tempIntBuff)
	treeBuf := make([]byte, count)
	buf.Read(treeBuf)

	buf.Read(tempIntBuff)
	count = getIntBytes(tempIntBuff)
	valid := make([]byte, count)
	buf.Read(valid)

	dataSize := 8 + 4
	var deserialize func(int) *kdTreeNode
	deserialize = func(offset int) *kdTreeNode {
		if offset*dataSize >= len(treeBuf) || valid[offset] == '0' {
			return nil
		}
		node := new(kdTreeNode)
		tempInt := make([]byte, 4)
		for i := range tempInt {
			tempInt[i] = treeBuf[offset*dataSize+i]
		}
		node.dim = getIntBytes(tempInt)
		tempFloat := make([]byte, 8)
		for i := range tempFloat {
			tempFloat[i] = treeBuf[offset*dataSize+4+i]
		}
		node.value = getFloatBytes(tempFloat)
		node.left = deserialize(2*offset + 1)
		node.right = deserialize(2*offset + 2)
		return node
	}
	n := deserialize(0)
	r.dim = n.dim
	r.value = n.value
	r.left = n.left
	r.right = n.right
}

func (r *kdTreeNode) Height() int {
	var treeHeight func(*kdTreeNode) int
	treeHeight = func(node *kdTreeNode) int {
		if node == nil {
			return 0
		}
		left := 1 + treeHeight(node.left)
		right := 1 + treeHeight(node.right)
		if left > right {
			return left
		}
		return right
	}
	return treeHeight(r)
}

func (r *kdTreeNode) MinHeight() int {
	var treeHeight func(*kdTreeNode) int
	treeHeight = func(node *kdTreeNode) int {
		if node == nil {
			return 0
		}
		left := 1 + treeHeight(node.left)
		right := 1 + treeHeight(node.right)
		if left < right {
			return left
		}
		return right
	}
	return treeHeight(r)
}

func (r *kdTreeNode) Prune(level int) {
	target := level - 1
	var dfs func(*kdTreeNode, int)
	dfs = func(node *kdTreeNode, level int) {
		if node == nil {
			return
		}
		if level == target {
			node.left = nil
			node.right = nil
		} else {
			dfs(node.left, level+1)
			dfs(node.right, level+1)
		}
	}
	dfs(r, 0)
}

func (r *kdTreeNode) GetLeafIndex(tuples []DatasetTuple) []int {
	treeHeight := r.Height()
	var dfs func(*kdTreeNode, int, int, DatasetTuple) int
	dfs = func(node *kdTreeNode, id, level int, tup DatasetTuple) int {
		if node == nil { // leaf
			return id << uint(treeHeight-level)
		}
		if tup.Data[node.dim] <= node.value {
			return dfs(node.left, (id<<1)|0, level+1, tup)
		}
		return dfs(node.right, (id<<1)|1, level+1, tup)
	}
	results := make([]int, 2<<uint(treeHeight-1))
	for _, tup := range tuples {
		results[dfs(r, 0, 0, tup)]++
	}
	return results
}

func (r kdTreeNode) String() string {
	var myString func(*kdTreeNode, string) string
	myString = func(node *kdTreeNode, pad string) string {
		if node != nil {
			return fmt.Sprintf("%s(%d - %.5f)\n",
				pad, node.dim, node.value) +
				myString(node.left, pad+"\t") +
				myString(node.right, pad+"\t")
		}
		return ""
	}
	return myString(&r, "")
}

func (r *kdTreeNode) Leaves() []*kdTreeNode {
	var dfs func(*kdTreeNode)
	var leaves []*kdTreeNode
	dfs = func(n *kdTreeNode) {
		if n != nil {
			if n.left == nil && n.right == nil {
				leaves = append(leaves, n)
			} else {
				dfs(n.left)
				dfs(n.right)
			}
		}
	}
	dfs(r)
	return leaves
}

// partitions the tuples and stores the tree structure in the kdTreeNode ptr
func newKDTreePartition(tuples []DatasetTuple, cols []int) *kdTreeNode {
	findMedian :=
		func(tuples []DatasetTuple, dim int) (float64, []DatasetTuple, []DatasetTuple) {
			values := make([]float64, 0)
			for _, t := range tuples {
				values = append(values, t.Data[dim])
			}
			sort.Float64s(values)
			median := values[len(values)/2]
			left, right := make([]DatasetTuple, 0), make([]DatasetTuple, 0)
			for _, t := range tuples {
				if t.Data[dim] <= median {
					left = append(left, t)
				} else {
					right = append(right, t)
				}
			}
			return median, left, right
		}
	var partition func([]DatasetTuple, int, *kdTreeNode)
	partition = func(tuples []DatasetTuple, dim int, node *kdTreeNode) {
		median, left, right := findMedian(tuples, cols[dim%len(cols)])
		node.dim = cols[dim%len(cols)]
		node.value = median
		if len(left) > 0 && len(right) > 0 {
			node.right = new(kdTreeNode)
			node.left = new(kdTreeNode)
			partition(left, dim+1, node.left)
			partition(right, dim+1, node.right)
		}
	}
	node := new(kdTreeNode)
	partition(tuples, cols[0], node)
	return node
}
