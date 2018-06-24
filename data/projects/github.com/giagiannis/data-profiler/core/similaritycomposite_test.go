package core

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestCompositeCompute(t *testing.T) {
	datasets := createPoolBasedDatasets(10000, 10, 3)

	x := NewDatasetSimilarityEstimator(SimilarityTypeBhattacharyya, datasets)
	x.Configure(nil)
	x.Compute()
	f1, _ := ioutil.TempFile("/tmp", "xest")
	f1.Write(x.Serialize())
	f1.Close()
	defer os.Remove(f1.Name())

	x = NewDatasetSimilarityEstimator(SimilarityTypeCorrelation, datasets)
	x.Configure(nil)
	x.Compute()
	f2, _ := ioutil.TempFile("/tmp", "xest")
	f2.Write(x.Serialize())
	f2.Close()
	defer os.Remove(f2.Name())

	var configuration = map[string]string{
		"concurrency": "10",
		"expression":  "0.8*x + 0.2*y",
		"x":           f1.Name(),
		"y":           f2.Name(),
	}
	e := new(CompositeEstimator)
	e.datasets = datasets
	e.Configure(configuration)

	err := e.Compute()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	sm := e.SimilarityMatrix()
	smSanityCheck(sm, t)
	cleanDatasets(datasets)
}

func TestCompositeConfiguration(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 100, 3)
	e := createCompositeEstimator(datasets)
	if len(e.estimators) != 2 {
		t.Logf("Expected %d estimators and found %d\n", 2, len(e.estimators))
		t.Fail()
	}
	if _, ok := (e.estimators["x"]).(*BhattacharyyaEstimator); !ok {
		t.Logf("Expected type %s was not identified\n",
			"*core.BhattacharyyaEstimator")
		t.Fail()
	}
	if _, ok := (e.estimators["y"]).(*CorrelationEstimator); !ok {
		t.Logf("Expected type %s was not identified\n",
			"*core.CorrelationEstimator")
		t.Fail()
	}

	cleanDatasets(datasets)
}

func TestCompositeSerialization(t *testing.T) {
	datasets := createPoolBasedDatasets(1000, 10, 3)
	e := createCompositeEstimator(datasets)
	e.Compute()
	ser := e.Serialize()
	if ser == nil {
		t.Log("Nil serialization object")
		t.Fail()
	}
	newEst := new(CompositeEstimator)
	newEst.Deserialize(ser)
	if newEst.expression != e.expression {
		t.Log("Difference in expressions")
		t.Fail()
	}
	for k, est := range newEst.estimators {
		if val, ok := e.estimators[k]; ok {
			if reflect.TypeOf(val) != reflect.TypeOf(est) {
				t.Log("Estimator types are different")
				t.Fail()
			}
		} else {
			t.Log("Estimator not found")
			t.Fail()
		}
	}
	estimatorsCheck(newEst.AbstractDatasetSimilarityEstimator,
		e.AbstractDatasetSimilarityEstimator, t)
	cleanDatasets(datasets)
}

func createCompositeEstimator(datasets []*Dataset) *CompositeEstimator {
	x := NewDatasetSimilarityEstimator(SimilarityTypeBhattacharyya, datasets)
	x.Configure(nil)
	x.Compute()
	f1, _ := ioutil.TempFile("/tmp", "xest")
	f1.Write(x.Serialize())
	f1.Close()
	defer os.Remove(f1.Name())

	x = NewDatasetSimilarityEstimator(SimilarityTypeCorrelation, datasets)
	x.Configure(nil)
	x.Compute()
	f2, _ := ioutil.TempFile("/tmp", "yest")
	f2.Write(x.Serialize())
	f2.Close()
	defer os.Remove(f2.Name())

	var configuration = map[string]string{
		"concurrency": "10",
		"expression":  "0.8*x + 0.2*y",
		"x":           f1.Name(),
		"y":           f2.Name(),
	}
	e := new(CompositeEstimator)
	e.datasets = datasets
	e.Configure(configuration)
	return e
}
