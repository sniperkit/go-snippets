package core

import "testing"

func TestScriptSimilarityDatasetAnalysis(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 4)
	est := new(ScriptSimilarityEstimator)
	est.datasets = datasets
	conf := map[string]string{
		"concurrency": "10",
		"script":      analysisScript,
		"norm":        "1",
	}
	est.Configure(conf)
	results := est.analyzeDatasets()
	if len(results) != len(datasets) || len(results) == 0 {
		t.Log("Not all datasets analyzed")
		t.FailNow()
	}

	size := len(results[0])
	for _, val := range results {
		if len(val) != size {
			t.Log("Incorrect number of elements")
			t.FailNow()
		}
	}
	cleanDatasets(datasets)
}
func TestScriptSimilarityCompute(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 50, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeScript, datasets)
	conf := map[string]string{
		"concurrency": "10",
		"script":      analysisScript,
		"norm":        "1",
	}
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

func TestScriptSimilarityComputeAppxThres(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 50, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeScript, datasets)
	conf := map[string]string{
		"concurrency": "10",
		"script":      analysisScript,
		"norm":        "1",
	}
	est.Configure(conf)
	pol := DatasetSimilarityPopulationPolicy{
		PolicyType: PopulationPolicyAprx,
		Parameters: map[string]float64{
			"threshold": 0.95,
		},
	}
	est.SetPopulationPolicy(pol)
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	s := est.SimilarityMatrix()

	smSanityCheck(s, t)
	cleanDatasets(datasets)
}

func TestScriptSimilarityComputeAppxCnt(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 50, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeScript, datasets)
	conf := map[string]string{
		"concurrency": "10",
		"script":      analysisScript,
		"norm":        "1",
	}
	est.Configure(conf)
	pol := DatasetSimilarityPopulationPolicy{
		PolicyType: PopulationPolicyAprx,
		Parameters: map[string]float64{
			"count": 10,
		},
	}
	est.SetPopulationPolicy(pol)
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	s := est.SimilarityMatrix()

	smSanityCheck(s, t)
	cleanDatasets(datasets)
}

func TestScriptSimilaritySerialization(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 100, 4)
	//	est := NewDatasetSimilarityEstimator(SIMILARITY_TYPE_JACCARD, datasets)
	est := *new(ScriptSimilarityEstimator)
	est.datasets = datasets
	est.simType = scriptSimilarityTypeEuclidean
	conf := map[string]string{
		"concurrency": "10",
		"script":      analysisScript,
		"norm":        "1",
	}
	est.Configure(conf)

	pol := DatasetSimilarityPopulationPolicy{
		PolicyType: PopulationPolicyFull,
		Parameters: map[string]float64{},
	}
	est.SetPopulationPolicy(pol)
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	bytes := est.Serialize()

	newEst := *new(ScriptSimilarityEstimator)
	newEst.Deserialize(bytes)
	estimatorsCheck(est.AbstractDatasetSimilarityEstimator, newEst.AbstractDatasetSimilarityEstimator, t)

	if est.analysisScript != newEst.analysisScript {
		t.Log("Analysis script differs")
		t.Fail()
	}

	if est.simType != newEst.simType {
		t.Log("norm degree differs")
		t.Fail()
	}

	for k, v := range est.inverseIndex {
		if newEst.inverseIndex[k] != v {
			t.Log("Inverse index failed")
			t.Fail()
		}
	}

	for i, arr := range est.datasetCoordinates {
		for j, v := range arr {
			if v != newEst.datasetCoordinates[i][j] {
				t.Log("Coordinates failed")
				t.Fail()
			}
		}
	}
	if newEst.Similarity(datasets[0], datasets[1]) != newEst.SimilarityMatrix().Get(0, 1) {
		t.Log("Something is seriously wrong here", newEst.SimilarityMatrix().Get(0, 1), newEst.Similarity(datasets[0], datasets[1]))
		t.Fail()
	}
	cleanDatasets(datasets)

}

func TestScriptSimilarityCosine(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 3)
	s := NewDatasetSimilarityEstimator(SimilarityTypeScript, datasets)
	s.Configure(map[string]string{
		"script":      analysisScript,
		"concurrency": "1",
		"type":        "cosine",
	})
	err := s.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	mat := s.SimilarityMatrix()
	for i := range datasets {
		for j := range datasets {
			if mat.Get(i, j) < -1.0 || mat.Get(i, j) > 1.0 {
				t.Logf("Cosine Similarity out of bounds between [%d, %d] -> %.3f\n", i, j, mat.Get(i, j))
				//				t.FailNow()
			}
		}
	}
	cleanDatasets(datasets)
}
