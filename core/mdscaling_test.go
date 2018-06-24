package core

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMDScalingScript(t *testing.T) {
	datasets := createPoolBasedDatasets(2000, 10, 4)
	est := NewDatasetSimilarityEstimator(SimilarityTypeJaccard, datasets)
	est.Configure(map[string]string{
		"concurrency": "10",
	})
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	md := NewMDScaling(est.SimilarityMatrix(), 2, mdsScript+"-bad")
	err = md.Compute()
	if err == nil {
		t.Log(err)
		t.Fail()
	}
	if md.Coordinates() != nil {
		t.Log("Script should not have succeeded")
		t.FailNow()
	}
	cleanDatasets(datasets)
}

func TestMDScalingCompute(t *testing.T) {
	datasets := createPoolBasedDatasets(100, 50, 2)
	est := NewDatasetSimilarityEstimator(SimilarityTypeJaccard, datasets)
	est.Configure(map[string]string{
		"concurrency": "10",
	})
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	for _, k := range []int{1, 2, 3, 4, 5, 6} {
		md := NewMDScaling(est.SimilarityMatrix(), k, mdsScript)
		err = md.Compute()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if md.Coordinates() == nil {
			t.Log("Coordinates are nil")
			t.FailNow()
		}
		_, err := md.Variances()

		if err != nil {
			t.Log(err)
			t.Fail()
		}
	}
	cleanDatasets(datasets)
}

func TestDatasetCoordinatesSerialization(t *testing.T) {
	datasets := createPoolBasedDatasets(100, 50, 2)
	est := NewDatasetSimilarityEstimator(SimilarityTypeJaccard, datasets)
	est.Configure(map[string]string{
		"concurrency": "10",
	})
	err := est.Compute()
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	for _, k := range []int{1, 2, 3, 4, 5, 6} {
		md := NewMDScaling(est.SimilarityMatrix(), k, mdsScript)
		err = md.Compute()
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if md.Coordinates() == nil {
			t.Log("Coordinates are nil")
			t.FailNow()
		}
		coords := md.Coordinates()
		buf := SerializeCoordinates(coords)
		newCoords := DeserializeCoordinates(buf)
		for i := range coords {
			for j := range coords[i] {
				valStr := fmt.Sprintf("%.5f", coords[i][j])
				val, _ := strconv.ParseFloat(valStr, 64)
				if val != newCoords[i][j] {
					t.Log("Coordinates lost is serialization")
					t.Log(val, newCoords[i][j])
					t.Fail()
				}
			}
		}

	}
	cleanDatasets(datasets)
}
