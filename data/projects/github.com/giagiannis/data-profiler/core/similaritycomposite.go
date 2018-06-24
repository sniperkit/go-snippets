package core

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

// CompositeEstimator returns a similarity function based on
// compositions of simpler similarity expressions. The user needs to provide a
// formula containing the expression of the similarity function, e.g.:
// 0.8 x BHATTACHARRYA + 0.2 CORRELATION
// Note that it is the user's responsibility to guarantee that the overall
// expression remains within the limits of the similarity expression [0,1].
type CompositeEstimator struct {
	AbstractDatasetSimilarityEstimator

	datasetIndexes map[string]int

	// the slice of estimators to use for the similarity evaluation
	estimators map[string]DatasetSimilarityEstimator
	// the math expression used for the estimation
	expression string
}

// Compute method constructs the Similarity Matrix
func (e *CompositeEstimator) Compute() error {
	return datasetSimilarityEstimatorCompute(e)
}

// Similarity returns the similarity between two datasets
func (e *CompositeEstimator) Similarity(a, b *Dataset) float64 {
	expression, err := govaluate.NewEvaluableExpression(e.expression)
	if err != nil {
		log.Println(err)
		return 0.0
	}
	params := make(map[string]interface{})
	for k, est := range e.estimators {
		if est.SimilarityMatrix() != nil {
			params[k] = est.SimilarityMatrix().Get(e.datasetIndexes[a.Path()], e.datasetIndexes[b.Path()])
		} else {
			params[k] = est.Similarity(a, b)
		}
	}
	result, err := expression.Evaluate(params)
	if err != nil {
		log.Println(err)
		return 0.0
	}
	if val, ok := result.(float64); ok {
		return val
	}
	return -1.0
}

// Serialize returns an array of bytes representing the Estimator.
func (e *CompositeEstimator) Serialize() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(getBytesInt(int(SimilarityTypeComposite)))

	buffer.Write(datasetSimilarityEstimatorSerialize(
		e.AbstractDatasetSimilarityEstimator))

	// serialize expression
	buffer.WriteString(e.expression + "\n")

	// serialize estimators
	buffer.Write(getBytesInt(len(e.estimators)))
	for k, est := range e.estimators {
		buffer.WriteString(k + "\n")
		temp := est.Serialize()
		buffer.Write(getBytesInt(len(temp)))
		buffer.Write(temp)
	}

	return buffer.Bytes()
}

// Deserialize constructs an Estimator object based on the byte array provided.
func (e *CompositeEstimator) Deserialize(b []byte) {
	buffer := bytes.NewBuffer(b)
	tempInt := make([]byte, 4)
	buffer.Read(tempInt) // consume estimator type
	buffer.Read(tempInt)
	absEstBytes := make([]byte, getIntBytes(tempInt))
	buffer.Read(absEstBytes)
	e.AbstractDatasetSimilarityEstimator =
		*datasetSimilarityEstimatorDeserialize(absEstBytes)

	// parse expression
	line, _ := buffer.ReadString('\n')
	e.expression = strings.TrimSpace(line)

	// parse the estimators
	buffer.Read(tempInt)
	count := getIntBytes(tempInt)
	e.estimators = make(map[string]DatasetSimilarityEstimator)
	for i := 0; i < count; i++ {
		line, _ := buffer.ReadString('\n')
		key := strings.TrimSpace(line)
		buffer.Read(tempInt)
		length := getIntBytes(tempInt)
		tempBuff := make([]byte, length)
		buffer.Read(tempBuff)
		est := DeserializeSimilarityEstimator(tempBuff)
		e.estimators[key] = est
	}
}

// Configure provides the configuration parameters needed by the Estimator
func (e *CompositeEstimator) Configure(conf map[string]string) {
	e.estimators = make(map[string]DatasetSimilarityEstimator)
	e.concurrency = 1
	for k, v := range conf {
		if k == "concurrency" {
			val, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				log.Println(err)
			}
			e.concurrency = int(val)
		} else if k == "expression" {
			e.expression = v
		} else { // one of x1,x2,... etc.
			cnt, err := ioutil.ReadFile(v)
			if err != nil {
				log.Println(err)
			} else {
				est := DeserializeSimilarityEstimator(cnt)
				e.estimators[k] = est
			}
		}
	}
	log.Println("Creating inverse dataset index")
	e.datasetIndexes = make(map[string]int)
	for i, d := range e.datasets {
		e.datasetIndexes[d.Path()] = i
	}

}

// Options returns the applicable parameters needed by the Estimator.
func (e *CompositeEstimator) Options() map[string]string {
	return map[string]string{
		"concurrency": "max number of threads to be run in parallel",
		"expression": "the math expression that combines the estimators " +
			"e.g.: 0.2*x + 0.8*y " +
			"(x and y must be later defined)",
		"x": "the path of a a similarity estimator",
		"y": "the path of another similarity estimator",
	}
}
