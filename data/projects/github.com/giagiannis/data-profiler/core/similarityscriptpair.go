package core

import (
	"bytes"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

// ScriptPairSimilarityEstimator executes a script of the extraction of the
// similarity between each pair of datasets
type ScriptPairSimilarityEstimator struct {
	AbstractDatasetSimilarityEstimator
	analysisScript string // the analysis script to be executed
}

// Compute method constructs the Similarity Matrix
func (e *ScriptPairSimilarityEstimator) Compute() error {
	return datasetSimilarityEstimatorCompute(e)
}

// Similarity returns the similarity between the two datasets
func (e *ScriptPairSimilarityEstimator) Similarity(a, b *Dataset) float64 {
	return e.executeScript(a.Path(), b.Path())
}

// Configure sets a number of configuration parameters to the struct. Use this
// method before the execution of the computation
func (e *ScriptPairSimilarityEstimator) Configure(conf map[string]string) {
	if val, ok := conf["concurrency"]; ok {
		conv, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			log.Println(err)
		} else {
			e.concurrency = int(conv)
		}
	} else {
		e.concurrency = 1
	}
	if val, ok := conf["script"]; ok {
		e.analysisScript = val
	} else {
		log.Println("Analysis script not defined - exiting")
	}
}

// Options returns a list of options that the user can set
func (e *ScriptPairSimilarityEstimator) Options() map[string]string {
	return map[string]string{
		"concurrency": "max number of threads to run in parallel",
		"script":      "path of the analysis script to be executed",
	}
}

// Serialize returns a byte array that represents the struct is a serialized version
func (e *ScriptPairSimilarityEstimator) Serialize() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(getBytesInt(int(SimilarityTypeScriptPair)))

	buffer.Write(
		datasetSimilarityEstimatorSerialize(
			e.AbstractDatasetSimilarityEstimator))

	//	buffer.Write(getBytesInt(e.concurrency))
	buffer.WriteString(e.analysisScript + "\n")
	return buffer.Bytes()
}

// Deserialize parses a byte array and forms a ScriptSimilarityEstimator object
func (e *ScriptPairSimilarityEstimator) Deserialize(b []byte) {
	buffer := bytes.NewBuffer(b)
	tempInt := make([]byte, 4)
	buffer.Read(tempInt) // consume estimator type
	buffer.Read(tempInt)
	absEstBytes := make([]byte, getIntBytes(tempInt))
	buffer.Read(absEstBytes)
	e.AbstractDatasetSimilarityEstimator =
		*datasetSimilarityEstimatorDeserialize(absEstBytes)

	line, _ := buffer.ReadString('\n')
	e.analysisScript = strings.TrimSpace(line)
}

// executeScript executed the analysis script into the specified dataset
func (e *ScriptPairSimilarityEstimator) executeScript(pathA, pathB string) float64 {
	cmd := exec.Command(e.analysisScript, pathA, pathB)
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	result, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		log.Println(err)
		return math.NaN()
	}
	return result
}
