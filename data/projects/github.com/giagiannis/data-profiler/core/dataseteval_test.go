package core

import (
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"testing"
)

func TestDatasetEvaluatorOnline(t *testing.T) {
	params := map[string]string{"script": mlScript, "testset": testSet}
	eval, err := NewDatasetEvaluator(OnlineEval, params)
	if err != nil || eval == nil {
		log.Println(err)
		t.FailNow()
	}
	val, err := eval.Evaluate(trainSet)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if val < 0 {
		log.Println("Wrong conversion")
		t.FailNow()
	}
}

func TestDatasetEvaluatorFileBased(t *testing.T) {
	scores := NewDatasetScores()
	datasetStrings := []string{"a", "b", "c", "d"}
	for _, s := range datasetStrings {
		scores.Scores[s] = rand.Float64()
	}
	cnt, err := scores.Serialize()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	f, err := ioutil.TempFile("/tmp", "scores")
	defer os.Remove(f.Name())
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	f.Write(cnt)
	f.Close()

	eval, err := NewDatasetEvaluator(FileBasedEval, map[string]string{"scores": f.Name()})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	for k, v := range scores.Scores {
		sc, err := eval.Evaluate(k)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if sc != v {
			t.Log("Wrong value!")
			t.Fail()
		}
	}
	sc, err := eval.Evaluate("Non existent dataset")
	if !math.IsNaN(sc) || err == nil {
		t.Log("Should have returned an error, but it didn't!")
		t.Fail()
	}
}
