package core

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func TestClustering(t *testing.T) {
	datasets := make([]*Dataset, 500)
	for i := 0; i < len(datasets); i++ {
		datasets[i] = NewDataset(fmt.Sprintf("data-%d", i))
	}
	sim := NewDatasetSimilarities(len(datasets))
	for i := range datasets {
		for j := range datasets {
			sim.Set(i, j, rand.Float64())
		}
	}

	cluster := NewClustering(sim, datasets)
	cluster.SetConcurrency(10)
	err := cluster.Compute()
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	max, min := cluster.Results().Heights()
	for i := 0; i <= max+10; i++ {
		sum := 0
		for _, cl := range cluster.Results().GetClusters(i) {
			sum += len(cl)
		}
		if sum != len(datasets) {
			log.Println("Datasets lost in clusters")
			t.FailNow()
		}
	}
	unary := false
	for _, cl := range cluster.Results().GetClusters(min) {
		unary = unary || (len(cl) == 1)
	}
	if !unary {
		log.Println("No unary cluster found for min height")
		t.FailNow()
	}
	unary = true
	for _, cl := range cluster.Results().GetClusters(max) {
		unary = unary && (len(cl) == 1)
	}
	if !unary {
		log.Println("Cluster with more than 1 datasets found, although at max height")
		t.FailNow()
	}

}

func TestNewDendrogram(t *testing.T) {
	datasets := make([]*Dataset, 40)
	for i := 0; i < len(datasets); i++ {
		datasets[i] = NewDataset(fmt.Sprintf("data-%d", i))
	}
	d := NewDendrogram(datasets)
	for unm := d.getUnmerged(); d.root == nil; unm = d.getUnmerged() {
		err := d.merge(unm[0], unm[1])
		if err != nil {
			log.Println(err)
			t.FailNow()
		}
	}
	if len(d.unmerged) > 0 {
		log.Println("We still have unmerged nodes")
		t.FailNow()
	}
}
