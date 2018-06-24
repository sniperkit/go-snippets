package core

import "testing"

func TestJaccardCompute(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeJaccard, datasets)
	est.Configure(map[string]string{"concurrency": "10"})
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	s := est.SimilarityMatrix()
	smSanityCheck(s, t)
	cleanDatasets(datasets)
}

func TestJaccardComputeAppxCnt(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeJaccard, datasets)
	pol := DatasetSimilarityPopulationPolicy{
		PolicyType: PopulationPolicyAprx,
		Parameters: map[string]float64{
			"count": 20,
		},
	}
	est.SetPopulationPolicy(pol)
	est.Configure(map[string]string{"concurrency": "10"})
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	s := est.SimilarityMatrix()
	smSanityCheck(s, t)
	cleanDatasets(datasets)
}

func TestJaccardComputeAppxThres(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeJaccard, datasets)
	pol := DatasetSimilarityPopulationPolicy{
		PolicyType: PopulationPolicyAprx,
		Parameters: map[string]float64{
			"threshold": 0.95,
		},
	}
	est.SetPopulationPolicy(pol)
	est.Configure(map[string]string{"concurrency": "10"})
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	s := est.SimilarityMatrix()
	smSanityCheck(s, t)
	cleanDatasets(datasets)
}

func TestJaccardSerialization(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 4)
	//	est := NewDatasetSimilarityEstimator(SIMILARITY_TYPE_JACCARD, datasets)
	est := *new(JaccardEstimator)
	est.datasets = datasets
	pol := DatasetSimilarityPopulationPolicy{
		PolicyType: PopulationPolicyAprx,
		Parameters: map[string]float64{
			"count": 5.0,
		},
	}
	est.SetPopulationPolicy(pol)
	est.Configure(map[string]string{"concurrency": "10"})
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	bytes := est.Serialize()

	newEst := *new(JaccardEstimator)
	newEst.Deserialize(bytes)

	estimatorsCheck(est.AbstractDatasetSimilarityEstimator, newEst.AbstractDatasetSimilarityEstimator, t)
	cleanDatasets(datasets)
}
