package core

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"testing"
)

func TestBhattacharyyaSampledDatasets(t *testing.T) {
	datasets := createPoolBasedDatasets(2000, 20, 2)
	defer cleanDatasets(datasets)
	percentage := 0.2
	est := new(BhattacharyyaEstimator)
	est.datasets = datasets
	conf := map[string]string{
		"concurrency": "10",
		"tree.sr":     fmt.Sprintf("%.2f", percentage),
		"partitions":  "16",
	}
	est.Configure(conf)
	totalNoTuples := 0.0
	for _, d := range datasets {
		totalNoTuples += float64(len(d.Data()))
	}
	s := est.sampledDataset()
	merged := float64(len(s))
	if math.Abs(merged-totalNoTuples*percentage)/merged > .01 {
		t.Logf("Merged dataset is of weird size (%.0f != %.0f)\n", merged, totalNoTuples*percentage)
		t.FailNow()
	}

}

func TestBhattacharyyaCompute(t *testing.T) {
	datasets := createPoolBasedDatasets(20000, 20, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeBhattacharyya, datasets)
	conf := map[string]string{"concurrency": "10"}
	est.Configure(conf)
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	s := est.SimilarityMatrix()
	smSanityCheck(s, t)
	cleanDatasets(datasets)
}

func TestKdTree(t *testing.T) {
	datasets := createPoolBasedDatasets(20000, 1, 5)
	dataset := datasets[0]
	err := dataset.ReadFromFile()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	kd := newKDTreePartition(dataset.Data(), []int{0, 1, 2, 3, 4})
	kd.Prune(kd.Height() - 2)
	ids := kd.GetLeafIndex(dataset.Data())
	sum := 0
	for _, v := range ids {
		sum += v
	}
	if sum != len(dataset.Data()) {
		log.Println("Not all tuples are indexed!!")
		t.FailNow()
	}
	cleanDatasets(datasets)

}

func TestKdPruning(t *testing.T) {
	datasets := createPoolBasedDatasets(20000, 1, 5)
	dataset := datasets[0]
	err := dataset.ReadFromFile()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	height := int(math.Log2(float64(len(dataset.Data())))) + 1
	for i := 0; i < height; i++ {
		kd := newKDTreePartition(dataset.Data(), []int{0, 1, 2, 3, 4})
		if i < kd.Height() {
			kd.Prune(kd.Height() - i)
			ids := kd.GetLeafIndex(dataset.Data())
			sum := .0
			for _, v := range ids {
				sum += math.Sqrt(float64(v * v))
			}
			if sum != math.Sqrt(float64(len(dataset.Data())*len(dataset.Data()))) {
				t.FailNow()
			}
		}
	}

	cleanDatasets(datasets)
}

func TestBhattacharyyaComputeAppxCnt(t *testing.T) {
	datasets := createPoolBasedDatasets(2000, 50, 4)
	defer cleanDatasets(datasets)
	est := NewDatasetSimilarityEstimator(SimilarityTypeBhattacharyya, datasets)
	conf := map[string]string{"concurrency": "10"}
	est.Configure(conf)
	policy := DatasetSimilarityPopulationPolicy{
		PolicyType: PopulationPolicyAprx,
		Parameters: map[string]float64{
			"count": 20,
		},
	}
	est.SetPopulationPolicy(policy)
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	s := est.SimilarityMatrix()
	smSanityCheck(s, t)
}

func TestBhattacharyyaComputeAppxThres(t *testing.T) {
	datasets := createPoolBasedDatasets(2000, 50, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeBhattacharyya, datasets)
	conf := map[string]string{
		"concurrency": "10",
		"partitions":  "256",
		"tree.sr":     "0.1",
	}
	est.Configure(conf)
	policy := DatasetSimilarityPopulationPolicy{
		PolicyType: PopulationPolicyAprx,
		Parameters: map[string]float64{
			"threshold": 0.9999,
		},
	}
	est.SetPopulationPolicy(policy)
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	s := est.SimilarityMatrix()
	smSanityCheck(s, t)
	cleanDatasets(datasets)
}

func TestKdtreeNodeSerialization(t *testing.T) {
	datasets := createPoolBasedDatasets(10000, 1, 2)
	datasets[0].ReadFromFile()
	tree := newKDTreePartition(datasets[0].Data(), []int{0, 1})
	h := float64(tree.Height()) * (rand.NormFloat64()*0.2 + 0.5)
	tree.Prune(int(math.Ceil(h)))
	b := tree.Serialize()
	newTree := new(kdTreeNode)
	newTree.Deserialize(b)
	var dfsTraversal func(treeA, treeB *kdTreeNode) bool
	dfsTraversal = func(treeA, treeB *kdTreeNode) bool {
		if treeA == nil && treeB == nil {
			return true
		}
		if treeA.dim != treeB.dim || treeA.value != treeB.value {
			return false
		}
		return dfsTraversal(treeA.left, treeB.left) && dfsTraversal(treeA.right, treeB.right)
	}
	if !dfsTraversal(tree, newTree) {
		t.Log("Trees not equal")
		t.FailNow()
	}
	cleanDatasets(datasets)
}

func TestBhattacharyyaSerialization(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 50, 4)
	//	est := NewDatasetSimilarityEstimator(SIMILARITY_TYPE_JACCARD, datasets)
	est := *new(BhattacharyyaEstimator)
	est.datasets = datasets
	pol := DatasetSimilarityPopulationPolicy{
		PolicyType: PopulationPolicyFull,
		Parameters: map[string]float64{},
	}
	est.SetPopulationPolicy(pol)
	est.Configure(map[string]string{
		"concurrency": "10",
		"partitions":  "10",
	})
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	bytes := est.Serialize()

	newEst := *new(BhattacharyyaEstimator)
	newEst.Deserialize(bytes)

	estimatorsCheck(newEst.AbstractDatasetSimilarityEstimator,
		est.AbstractDatasetSimilarityEstimator, t)
	for k, v := range est.inverseIndex {
		if newEst.inverseIndex[k] != v {
			t.Log("Inverse index failed")
			t.Fail()
		}
	}

	for i, arr := range est.pointsPerRegion {
		for j, v := range arr {
			if v != newEst.pointsPerRegion[i][j] {
				t.Log("Points per region not the same", i, j)
				t.Fail()
			}
		}
	}

	for i, v := range est.datasetsSize {
		if v != newEst.datasetsSize[i] {
			t.Log("Datasets sizes not the same", i)
			t.Fail()
		}
	}

	cleanDatasets(datasets)
}
