package core

import (
	"log"
	"math"
	"math/rand"
	"testing"
)

func TestDatasetSimilarityPolicySerialization(t *testing.T) {
	pol := new(DatasetSimilarityPopulationPolicy)
	pol.Parameters = map[string]float64{"count": 10.0, "foo": 100.0}
	pol.PolicyType = PopulationPolicyAprx
	bytes := pol.Serialize()
	newPol := new(DatasetSimilarityPopulationPolicy)
	newPol.Deserialize(bytes)
	if newPol.PolicyType != pol.PolicyType {
		t.Log("Policy types different")
		t.Fail()
	}
	for k, v := range pol.Parameters {
		val, ok := newPol.Parameters[k]
		if !ok || val != v {
			t.Log("Missing key or different value", ok, val, v)
			t.Fail()
		}
	}
}

func TestDatasetSerialize(t *testing.T) {
	datasets := createPoolBasedDatasets(5000, 10, 3)
	e := NewDatasetSimilarityEstimator(SimilarityTypeBhattacharyya, datasets)
	e.Configure(nil)
	e.Compute()
	s := e.SimilarityMatrix()
	b := s.Serialize()

	sim := new(DatasetSimilarityMatrix)
	sim.Deserialize(b)

	for i := range datasets {
		if sim.closestIndex.similarity[i] != s.closestIndex.similarity[i] ||
			sim.closestIndex.closestIdx[i] != s.closestIndex.closestIdx[i] {
			log.Println("closestIndex contains different elements at dataset", i)
			log.Println(sim.closestIndex, s.closestIndex)
			t.FailNow()
		}
	}

	for i := range datasets {
		for j := range datasets {
			if sim.Get(i, j) != s.Get(j, i) {
				log.Printf("Sim and s provide different values: (%.5f vs %.5f)\n",
					sim.Get(i, j),
					s.Get(i, j))
				t.FailNow()
			}
		}
	}

	cleanDatasets(datasets)
}

func TestDatasetSimilarityClosestIndex(t *testing.T) {
	datasets := createPoolBasedDatasets(5000, 100, 3)
	sim := NewDatasetSimilarities(len(datasets))
	chosenIdx := rand.Int() % len(datasets)
	for i := range datasets {
		val := rand.Float64()
		if i == chosenIdx {
			val = 1
		}
		sim.Set(chosenIdx, i, val)
	}

	for i := range datasets {
		closest, _ := sim.closestIndex.Get(i)
		if closest != chosenIdx {
			t.Log("Closest idx is not accurate")
			t.FailNow()
		}
	}

	cleanDatasets(datasets)
}

func TestDatasetSimilaritiesDisabledIndex(t *testing.T) {
	datasets := createPoolBasedDatasets(5000, 100, 3)
	sim := NewDatasetSimilarities(len(datasets))
	sim.IndexDisabled(true)
	chosenIdx := rand.Int() % len(datasets)
	for i := range datasets {
		val := rand.Float64()
		if i == chosenIdx {
			val = 1
		}
		sim.Set(chosenIdx, i, val)
	}

	for i := range datasets {
		closest, valI := sim.closestIndex.Get(i)
		valA := sim.Get(chosenIdx, i)
		if closest != -1 || valI != -1 || valA <= 0.0 {
			t.Log("Closest idx is wrong")
			t.FailNow()
		}
	}
	cleanDatasets(datasets)
}

func TestDeserializeSimilarityEstimator(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 10, 2)
	// jaccard
	est := NewDatasetSimilarityEstimator(SimilarityTypeJaccard, datasets)
	est.Configure(map[string]string{"concurrency": "10"})
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		buf := est.Serialize()
		newEst := DeserializeSimilarityEstimator(buf)
		idxA, idxB := rand.Intn(len(datasets)), rand.Intn(len(datasets))
		a, b := datasets[idxA], datasets[idxB]
		if newEst.Similarity(a, b) != est.Similarity(a, b) ||
			newEst.Similarity(a, b) != est.SimilarityMatrix().Get(idxA, idxB) {
			t.Log("Deserialization error")
			t.Fail()
		}
	}
	// bhat
	est = NewDatasetSimilarityEstimator(SimilarityTypeBhattacharyya, datasets)
	est.Configure(map[string]string{"concurrency": "10"})
	err = est.Compute()
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		buf := est.Serialize()
		newEst := DeserializeSimilarityEstimator(buf)
		idxA, idxB := rand.Intn(len(datasets)), rand.Intn(len(datasets))
		a, b := datasets[idxA], datasets[idxB]
		if newEst.Similarity(a, b) != est.Similarity(a, b) ||
			newEst.Similarity(a, b) != est.SimilarityMatrix().Get(idxA, idxB) {
			t.Log("Deserialization error")
			t.Fail()
		}
	}

	// script
	est = NewDatasetSimilarityEstimator(SimilarityTypeScript, datasets)
	est.Configure(map[string]string{"concurrency": "10", "script": analysisScript})
	err = est.Compute()
	if err != nil {
		t.Log(err)
		t.Fail()
	} else {
		buf := est.Serialize()
		newEst := DeserializeSimilarityEstimator(buf)
		idxA, idxB := rand.Intn(len(datasets)), rand.Intn(len(datasets))
		a, b := datasets[idxA], datasets[idxB]
		val1, val2 := newEst.Similarity(a, b), est.Similarity(a, b)
		if val1 != val2 && !math.IsNaN(val1) && !math.IsNaN(val2) {
			t.Log("Deserialization error")
			t.Log(est.Similarity(a, b))
			t.Log(newEst.SimilarityMatrix())
			t.Log(est.SimilarityMatrix())
			t.Fail()
		}
	}
	cleanDatasets(datasets)
}

func smGetCountOfOnesAndZeros(sm *DatasetSimilarityMatrix) (int, int) {
	zeroElem, oneElem := 0, 0
	for i := 0; i < sm.Capacity(); i++ {
		for j := 0; j < sm.Capacity(); j++ {
			val := sm.Get(i, j)
			if val == 0.0 {
				zeroElem++
			}
			if val == 1.0 {
				oneElem++
			}
		}
	}
	return zeroElem, oneElem
}

func smCheckSymmetry(sm *DatasetSimilarityMatrix) bool {
	for i := 0; i < sm.Capacity(); i++ {
		for j := 0; j < sm.Capacity(); j++ {
			if sm.Get(i, j) != sm.Get(j, i) {
				return false
			}
		}
	}
	return true
}

func smCheckDiagonal(sm *DatasetSimilarityMatrix) bool {
	for i := 0; i < sm.Capacity(); i++ {
		if sm.Get(i, i) != 1.0 {
			return false
		}
	}
	return true

}

func smElementsInRange(sm *DatasetSimilarityMatrix) bool {
	for i := 0; i < sm.Capacity(); i++ {
		for j := 0; j < sm.Capacity(); j++ {
			if sm.Get(i, j) < 0 || sm.Get(i, j) > 1.0 {
				return false
			}
		}
	}
	return true
}

func smSanityCheck(s *DatasetSimilarityMatrix, t *testing.T) {
	if s == nil {
		t.Log("Nil similarities returned")
		t.FailNow()
	}
	if !smCheckSymmetry(s) {
		t.Log("SM not symmetrical")
		t.Fail()
	}

	if !smCheckDiagonal(s) {
		t.Log("The diagonal elements are not 1s!")
		t.Fail()
	}
	if !smElementsInRange(s) {
		t.Log("Elements out of range found!")
		t.Fail()
	}
	zeros, ones := smGetCountOfOnesAndZeros(s)
	if zeros > 3*s.Capacity()*s.Capacity()/4 {
		t.Logf("Too many zero elements found (%d)\n", zeros)
		t.Fail()
	}

	if ones < s.Capacity() || ones > 4*s.Capacity()*s.Capacity()/5 {
		t.Logf("Wrong number of 1s found in the SM (%d)\n", ones)
		t.Fail()
	}

}

func estimatorsCheck(a, b AbstractDatasetSimilarityEstimator, t *testing.T) {
	if a.concurrency != b.concurrency {
		t.Logf("Concurrency: expected %d, found %d\n", a.concurrency, b.concurrency)
		t.Fail()
	}

	if a.duration != b.duration {
		t.Logf("Duration: expected %.5f, found %.5f\n", a.duration, b.duration)
		t.Fail()
	}

	for i := range a.datasets {
		if a.datasets[i].Path() != b.datasets[i].Path() {
			t.Logf("Dataset paths differ")
			t.Fail()
		}
	}
	if a.popPolicy.PolicyType != b.popPolicy.PolicyType {
		t.Logf("Dataset paths differ")
		t.Fail()
	}
	if len(a.popPolicy.Parameters) != len(b.popPolicy.Parameters) {
		t.Logf("Poppolicy parameters differ")
		t.Fail()
	}
	for k, v := range a.popPolicy.Parameters {
		if val, ok := b.popPolicy.Parameters[k]; ok {
			if val != v {
				t.Log("Parameter is different")
				t.Fail()
			}
		} else {
			t.Log("Parameter not found")
			t.Fail()
		}
	}
	count := a.similarities.Capacity()
	if count != b.similarities.Capacity() {
		t.Log("SM have different size")
		t.Fail()
	}
	for i := 0; i < count; i++ {
		for j := 0; j < count; j++ {
			if a.similarities.Get(i, j) != b.similarities.Get(i, j) {
				t.Log("SM have different elements")
				t.Fail()
			}
		}
	}
}
func smPrint(sm *DatasetSimilarityMatrix, t *testing.T) {
	for i := 0; i < sm.Capacity(); i++ {
		for j := 0; j < sm.Capacity(); j++ {
			t.Logf("%.10f\t", sm.Get(i, j))
		}
		t.Logf("\n")
	}
}
