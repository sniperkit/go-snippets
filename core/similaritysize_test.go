package core

import "testing"

func TestSizeCompute(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeSize, datasets)
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

func TestSizeComputeAppxCnt(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeSize, datasets)
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

func TestSizeComputeAppxThres(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeSize, datasets)
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

func TestSizeSerialization(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 20, 4)
	//	est := NewDatasetSimilarityEstimator(SIMILARITY_TYPE_JACCARD, datasets)
	est := *new(SizeEstimator)
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

	newEst := *new(SizeEstimator)
	newEst.Deserialize(bytes)

	estimatorsCheck(est.AbstractDatasetSimilarityEstimator, newEst.AbstractDatasetSimilarityEstimator, t)
	cleanDatasets(datasets)
}
