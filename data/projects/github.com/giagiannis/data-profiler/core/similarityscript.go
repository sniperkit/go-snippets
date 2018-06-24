package core

import (
	"bytes"
	"errors"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

// ScriptSimilarityEstimator utilizes a script to analyze the data based on some external
// algorithm and utilizes various norms to measure the differences between the
// analysis outputs.
type ScriptSimilarityEstimator struct {
	AbstractDatasetSimilarityEstimator
	analysisScript     string                        // the analysis script to be executed
	simType            ScriptSimilarityEstimatorType // similarity type - cosine, manhattan, euclidean
	inverseIndex       map[string]int                // inverse index that maps datasets to ints
	datasetCoordinates [][]float64                   // holds the dataset coordinates
}

// ScriptSimilarityEstimatorType reflects the type of the ScriptSimilarityEstimator
type ScriptSimilarityEstimatorType uint8

const (
	scriptSimilarityTypeManhattan ScriptSimilarityEstimatorType = iota
	scriptSimilarityTypeEuclidean ScriptSimilarityEstimatorType = iota + 1
	scriptSimilarityTypeCosine    ScriptSimilarityEstimatorType = iota + 2
)

// Compute method constructs the Similarity Matrix
func (e *ScriptSimilarityEstimator) Compute() error {
	return datasetSimilarityEstimatorCompute(e)
}

// Similarity returns the similarity between the two datasets
func (e *ScriptSimilarityEstimator) Similarity(a, b *Dataset) float64 {
	var coordsA, coordsB []float64
	if id, ok := e.inverseIndex[a.Path()]; ok {
		coordsA = e.datasetCoordinates[id]
	} else {
		coordsA = e.analyzeDataset(a.Path())
	}
	if id, ok := e.inverseIndex[b.Path()]; ok {
		coordsB = e.datasetCoordinates[id]
	} else {
		coordsB = e.analyzeDataset(b.Path())
	}
	if e.simType == scriptSimilarityTypeCosine {
		val, err := e.cosine(coordsA, coordsB)
		if err != nil {
			log.Println(err)
		}
		return val
	}
	normDegree := 2 // default is EUCLIDEAN distance
	if e.simType == scriptSimilarityTypeManhattan {
		normDegree = 1
	}
	val, err := e.norm(coordsA, coordsB, normDegree)
	if err != nil {
		log.Println(err)
	}
	return DistanceToSimilarity(val)
}

// Configure sets a number of configuration parameters to the struct. Use this
// method before the execution of the computation
func (e *ScriptSimilarityEstimator) Configure(conf map[string]string) {
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
	if val, ok := conf["type"]; ok {
		if val == "cosine" {
			e.simType = scriptSimilarityTypeCosine
		} else if val == "manhattan" {
			e.simType = scriptSimilarityTypeManhattan
		} else if val == "euclidean" {
			e.simType = scriptSimilarityTypeEuclidean
		} else {
			log.Println("Similarity Type not known, valid values: [cosine manhattan euclidean]")
		}
	} else {
		e.simType = scriptSimilarityTypeEuclidean
	}

	// execute analysis for each dataset
	log.Println("Analyzing datasets")
	e.datasetCoordinates = e.analyzeDatasets()
	e.inverseIndex = make(map[string]int)
	for i, d := range e.datasets {
		e.inverseIndex[d.Path()] = i
	}

}

// Options returns a list of options that the user can set
func (e *ScriptSimilarityEstimator) Options() map[string]string {
	return map[string]string{
		"concurrency": "max number of threads to run in parallel",
		"script":      "path of the analysis script to be executed",
		"type":        "the type of the similarity - one of:  [cosine manhattan euclidean]",
	}
}

// Serialize returns a byte array that represents the struct is a serialized version
func (e *ScriptSimilarityEstimator) Serialize() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(getBytesInt(int(SimilarityTypeScript)))

	buffer.Write(
		datasetSimilarityEstimatorSerialize(
			e.AbstractDatasetSimilarityEstimator))

	//	buffer.Write(getBytesInt(e.concurrency))
	buffer.Write(getBytesInt(int(e.simType)))
	buffer.WriteString(e.analysisScript + "\n")

	// write number of coordinates per dataset
	buffer.Write(getBytesInt(len(e.datasetCoordinates[0])))
	for _, arr := range e.datasetCoordinates {
		for _, v := range arr {
			buffer.Write(getBytesFloat(v))
		}
	}
	return buffer.Bytes()
}

// Deserialize parses a byte array and forms a ScriptSimilarityEstimator object
func (e *ScriptSimilarityEstimator) Deserialize(b []byte) {
	buffer := bytes.NewBuffer(b)
	tempInt := make([]byte, 4)
	buffer.Read(tempInt) // consume estimator type
	buffer.Read(tempInt)
	absEstBytes := make([]byte, getIntBytes(tempInt))
	buffer.Read(absEstBytes)
	e.AbstractDatasetSimilarityEstimator =
		*datasetSimilarityEstimatorDeserialize(absEstBytes)

	buffer.Read(tempInt)
	e.simType = ScriptSimilarityEstimatorType(getIntBytes(tempInt))
	line, _ := buffer.ReadString('\n')
	e.analysisScript = strings.TrimSpace(line)

	e.inverseIndex = make(map[string]int)
	for i, d := range e.datasets {
		e.inverseIndex[d.Path()] = i
	}

	tempFloat := make([]byte, 8)
	buffer.Read(tempInt)
	count := getIntBytes(tempInt)
	e.datasetCoordinates = make([][]float64, len(e.datasets))
	for i := range e.datasets {
		e.datasetCoordinates[i] = make([]float64, count)
		for j := range e.datasetCoordinates[i] {
			buffer.Read(tempFloat)
			e.datasetCoordinates[i][j] = getFloatBytes(tempFloat)
		}
	}
}

func (e *ScriptSimilarityEstimator) analyzeDatasets() [][]float64 {
	c, done := make(chan bool, e.concurrency), make(chan bool)
	coords := make([][]float64, len(e.datasets))
	for i := 0; i < e.concurrency; i++ {
		c <- true
	}
	for i, d := range e.datasets {
		go func(c, done chan bool, i int, path string) {
			<-c
			coords[i] = e.analyzeDataset(path)
			c <- true
			done <- true
		}(c, done, i, d.Path())
	}

	for i := 0; i < len(e.datasets); i++ {
		<-done
	}
	return coords
}

// analyzeDataset executed the analysis script into the specified dataset
func (e *ScriptSimilarityEstimator) analyzeDataset(path string) []float64 {
	log.Println("Analyzing", path)
	cmd := exec.Command(e.analysisScript, path)
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	results := make([]float64, 0)
	for _, sv := range strings.Split(string(out), "\t") {
		conv, err := strconv.ParseFloat(strings.TrimSpace((sv)), 64)
		if err == nil {
			results = append(results, conv)
		} else {
			log.Println(err)
		}
	}
	log.Println("Tuple read:", results)
	return results
}

// norm function calculates the norm between two float slices
func (e *ScriptSimilarityEstimator) norm(a, b []float64, normDegree int) (float64, error) {
	if len(a) != len(b) {
		return -1, errors.New("arrays have different sizes")
	}
	sum := 0.0
	for i := range a {
		dif := math.Abs(a[i] - b[i])
		sum += math.Pow(dif, float64(normDegree))
	}
	return math.Pow(sum, 1.0/float64(normDegree)), nil
}

// cosine calculates the cosine similarity between two vectors
func (e *ScriptSimilarityEstimator) cosine(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		return -1, errors.New("arrays have different sizes")
	}
	nomin, sumA, sumB := 0.0, 0.0, 0.0
	for i := range a {
		nomin += a[i] * b[i]
		sumA += a[i] * a[i]
		sumB += b[i] * b[i]
	}
	denom := math.Sqrt(sumA) * math.Sqrt(sumB)
	if denom == 0.0 {
		return -1, errors.New("Zero denominator to cosine similarity")
	}
	return nomin / denom, nil
}
