package core

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

// DatasetEvaluatorType represents the type of the dataset evaluator
type DatasetEvaluatorType uint8

const (
	// OnlineEval dynamically parses the dataset values
	OnlineEval DatasetEvaluatorType = iota + 1
	// FileBasedEval returns the pre-computed values of an operator
	FileBasedEval
)

// DatasetEvaluator reflects the interface of an evaluator object.
type DatasetEvaluator interface {
	Evaluate(string) (float64, error)
}

// NewDatasetEvaluator returns a new DatasetEvaluator object
func NewDatasetEvaluator(evalType DatasetEvaluatorType,
	params map[string]string) (DatasetEvaluator, error) {
	if evalType == OnlineEval {
		eval := new(OnlineDatasetEvaluator)
		if _, ok := params["script"]; !ok {
			return nil, errors.New("Online evaluator needs script param")
		}
		eval.script = params["script"]
		if _, ok := params["testset"]; !ok {
			return nil, errors.New("Online evaluator needs testset param")
		}
		eval.testset = params["testset"]
		return eval, nil
	} else if evalType == FileBasedEval {
		eval := new(FileBasedEvaluator)
		if _, ok := params["scores"]; !ok {
			return nil, errors.New("File based evaluator needs scores file")
		}
		f, err := os.Open(params["scores"])
		if err != nil {
			return nil, err
		}
		scores := NewDatasetScores()
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		err = scores.Deserialize(buf)
		if err != nil {
			return nil, err
		}
		eval.scores = *scores
		return eval, nil
	}
	return nil, errors.New("Not suitable evaluator found")
}

// OnlineDatasetEvaluator is responsible to execute the training script and fetch
// the model accuracy
type OnlineDatasetEvaluator struct {
	script  string
	testset string
}

// Evaluate evaluates a new dataset, based on its path
func (e *OnlineDatasetEvaluator) Evaluate(dataset string) (float64, error) {
	cmd := exec.Command(e.script, dataset, e.testset)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err, string(out))
		return -1, err
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		log.Println(err, "for dataset", dataset)
		return -1, err
	}
	return val, nil
}

// FileBasedEvaluator returns the scores of an operator based on a scores file.
type FileBasedEvaluator struct {
	scores DatasetScores
}

// Evaluate returns the score for a given dataset
func (e *FileBasedEvaluator) Evaluate(dataset string) (float64, error) {
	val, ok := e.scores.Scores[dataset]
	if !ok { // try without the path
		val, ok = e.scores.Scores[path.Base(dataset)]
	}
	if !ok {
		return math.NaN(), errors.New("Dataset not found")
	}
	return val, nil
}

// DatasetScores is used to store the scores of a set of datasets
type DatasetScores struct {
	Scores map[string]float64
}

// NewDatasetScores initializes a new DatasetScores struct
func NewDatasetScores() *DatasetScores {
	o := new(DatasetScores)
	o.Scores = make(map[string]float64)
	return o
}

// Serialize returns a stream containing a DatasetScores object
func (s *DatasetScores) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	e := gob.NewEncoder(buf)
	err := e.Encode(s.Scores)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize constructs a DatasetScores strucy based on a byte array
func (s *DatasetScores) Deserialize(buf []byte) error {
	content := bytes.NewBuffer(buf)
	d := gob.NewDecoder(content)
	err := d.Decode(&s.Scores)
	return err
}
