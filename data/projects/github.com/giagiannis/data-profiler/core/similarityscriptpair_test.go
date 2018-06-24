package core

import "testing"

func TestScriptPairSimilarityCompute(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 30, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeScriptPair, datasets)
	conf := map[string]string{
		"concurrency": "10",
		"script":      pairScript,
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

func TestScriptPairSimilaritySerialization(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 30, 4)
	//	est := NewDatasetSimilarityEstimator(SIMILARITY_TYPE_JACCARD, datasets)
	est := *new(ScriptPairSimilarityEstimator)
	est.datasets = datasets
	conf := map[string]string{
		"concurrency": "10",
		"script":      pairScript,
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

	newEst := *new(ScriptPairSimilarityEstimator)
	newEst.Deserialize(bytes)
	estimatorsCheck(est.AbstractDatasetSimilarityEstimator, newEst.AbstractDatasetSimilarityEstimator, t)

	if est.analysisScript != newEst.analysisScript {
		t.Log("Analysis script differs")
		t.Log(est.analysisScript, newEst.analysisScript)
		t.Fail()
	}

	//	if newEst.Similarity(datasets[0], datasets[1]) != newEst.SimilarityMatrix().Get(0, 1) {
	//		t.Log("Something is seriously wrong here", newEst.SimilarityMatrix().Get(0, 1), newEst.Similarity(datasets[0], datasets[1]))
	//		t.Fail()
	//	}
	cleanDatasets(datasets)

}
