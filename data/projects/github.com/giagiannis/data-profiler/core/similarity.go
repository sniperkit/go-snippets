package core

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
)

// DatasetSimilarityEstimator is the interface that each Similarity estimator
// obeys.
type DatasetSimilarityEstimator interface {
	// computes the similarity matrix
	Compute() error
	// returns the datasets slice
	Datasets() []*Dataset
	// returns the similarity for 2 datasets
	Similarity(a, b *Dataset) float64
	// returns the similarity struct
	SimilarityMatrix() *DatasetSimilarityMatrix
	// provides configuration options
	Configure(map[string]string)
	// list of options for the estimator
	Options() map[string]string
	// sets the population policy for the estimator
	SetPopulationPolicy(DatasetSimilarityPopulationPolicy)
	// returns the population policy
	PopulationPolicy() DatasetSimilarityPopulationPolicy
	// returns a serialized esimator object
	Serialize() []byte
	// instantiates an estimator from a serialized object
	Deserialize([]byte)
	// returns the seconds needed to execute the computation
	Duration() float64
	// returns the max number of threads to be used
	Concurrency() int

	// sets the duration
	setDuration(float64)
	// sets the similarity matrix
	setSimilarityMatrix(*DatasetSimilarityMatrix)
}

// AbstractDatasetSimilarityEstimator is the base struct for the similarity
// estimator objects
type AbstractDatasetSimilarityEstimator struct {
	datasets     []*Dataset
	popPolicy    DatasetSimilarityPopulationPolicy
	similarities *DatasetSimilarityMatrix
	duration     float64
	concurrency  int
}

// Datasets returns the datasets of the estimator
func (a *AbstractDatasetSimilarityEstimator) Datasets() []*Dataset {
	return a.datasets
}

// SimilarityMatrix returns the similarity matrix of the estimator
func (a *AbstractDatasetSimilarityEstimator) SimilarityMatrix() *DatasetSimilarityMatrix {
	return a.similarities
}

func (a *AbstractDatasetSimilarityEstimator) setSimilarityMatrix(sm *DatasetSimilarityMatrix) {
	a.similarities = sm
}

// SetPopulationPolicy sets the population policy to be used
func (a *AbstractDatasetSimilarityEstimator) SetPopulationPolicy(pol DatasetSimilarityPopulationPolicy) {
	a.popPolicy = pol
}

// PopulationPolicy gets the population policy to be used
func (a *AbstractDatasetSimilarityEstimator) PopulationPolicy() DatasetSimilarityPopulationPolicy {
	return a.popPolicy
}

// Duration returns the duration of the compution
func (a *AbstractDatasetSimilarityEstimator) Duration() float64 {
	return a.duration
}

// setDuration sets the duration of the compution
func (a *AbstractDatasetSimilarityEstimator) setDuration(d float64) {
	a.duration = d
}

// Concurrency returns the max number of threads to be used for the computation
func (a *AbstractDatasetSimilarityEstimator) Concurrency() int {
	return a.concurrency
}

// datasetSimilarityEstimatorSerialize is used to generate an array of bytes of
// the abstract object
func datasetSimilarityEstimatorSerialize(e AbstractDatasetSimilarityEstimator) []byte {
	buffer := new(bytes.Buffer)

	buffer.Write(getBytesInt(len(e.Datasets())))
	for _, d := range e.Datasets() {
		buffer.WriteString(d.Path() + "\n")
	}
	pop := e.PopulationPolicy()
	popSer := pop.Serialize()
	buffer.Write(getBytesInt(len(popSer)))
	buffer.Write(popSer)

	sim := e.SimilarityMatrix().Serialize()
	buffer.Write(getBytesInt(len(sim)))
	buffer.Write(sim)
	buffer.Write(getBytesInt(e.Concurrency()))
	buffer.Write(getBytesFloat(e.Duration()))
	cnt := buffer.Bytes()
	bufLen := getBytesInt(len(cnt))
	return append(bufLen, cnt...)
}

// datasetSimilarityEstimatorDeserialize is used to generate an object based on
// the byte stream
func datasetSimilarityEstimatorDeserialize(b []byte) *AbstractDatasetSimilarityEstimator {
	result := new(AbstractDatasetSimilarityEstimator)
	buffer := bytes.NewBuffer(b)
	tempInt := make([]byte, 4)
	tempFloat := make([]byte, 8)

	buffer.Read(tempInt)
	count := getIntBytes(tempInt)
	result.datasets = make([]*Dataset, count)
	for i := range result.datasets {
		line, _ := buffer.ReadString('\n')
		line = strings.TrimSpace(line)
		result.datasets[i] = NewDataset(line)
	}

	buffer.Read(tempInt)
	count = getIntBytes(tempInt)
	polBytes := make([]byte, count)
	buffer.Read(polBytes)
	result.popPolicy = *new(DatasetSimilarityPopulationPolicy)
	result.popPolicy.Deserialize(polBytes)

	buffer.Read(tempInt)
	count = getIntBytes(tempInt)
	similarityBytes := make([]byte, count)
	buffer.Read(similarityBytes)
	result.similarities = new(DatasetSimilarityMatrix)
	result.similarities.Deserialize(similarityBytes)

	buffer.Read(tempInt)
	result.concurrency = getIntBytes(tempInt)

	buffer.Read(tempFloat)
	result.duration = getFloatBytes(tempFloat)
	return result
}

// datasetSimilarityEstimatorCompute is responsible to execute the computation code of the estimators.
// The provided object must respect the DatasetSimilarityEstimator interface
// and (optionally) extends the AbstractDatasetSimilarityEstimator struct
func datasetSimilarityEstimatorCompute(e DatasetSimilarityEstimator) error {
	err := readDatasets(e)
	if err != nil {
		return err
	}
	start := time.Now()
	if e.PopulationPolicy().PolicyType == PopulationPolicyFull {
		e.SimilarityMatrix().IndexDisabled(true) // I don't need the index
		log.Println("Computing the similarities using", e.Concurrency(), "threads")
		c := make(chan bool, e.Concurrency())
		done := make(chan bool)
		for j := 0; j < e.Concurrency(); j++ {
			c <- true
		}
		for i := 0; i < len(e.Datasets())-1; i++ {
			go func(c, done chan bool, i int) {
				<-c
				for j := i; j < len(e.Datasets()); j++ {
					d1, d2 := e.Datasets()[i], e.Datasets()[j]
					e.SimilarityMatrix().Set(i, j, e.Similarity(d1, d2))
				}
				c <- true
				done <- true
			}(c, done, i)
		}
		for j := 0; j < len(e.Datasets())-1; j++ {
			<-done
		}
		log.Println("Done")
	} else if e.PopulationPolicy().PolicyType == PopulationPolicyAprx {
		e.SimilarityMatrix().IndexDisabled(false) // I need the index
		if count, ok := e.PopulationPolicy().Parameters["count"]; ok {
			log.Printf("Fixed number of points execution (count: %.0f)\n", count)
			chosenIdxs := make(map[int]bool)
			for i := 0.0; i < count; i++ {
				var idx int
				if _, ok2 := e.PopulationPolicy().Parameters["random"]; ok2 {
					for len(chosenIdxs) < len(e.Datasets()) {
						idx = rand.Intn(len(e.Datasets()))
						if _, ok := chosenIdxs[idx]; !ok {
							chosenIdxs[idx] = true
							break
						}
					}
				} else {
					idx, _ = e.SimilarityMatrix().LeastSimilar()
				}
				log.Println("Computing the similarities for ", idx)
				for j := 0; j < len(e.Datasets()); j++ {
					d1, d2 := e.Datasets()[idx], e.Datasets()[j]
					e.SimilarityMatrix().Set(idx, j, e.Similarity(d1, d2))
				}

			}
		} else if threshold, ok := e.PopulationPolicy().Parameters["threshold"]; ok {
			log.Printf("Threshold based execution (threshold: %.5f)\n", threshold)
			idx, val := e.SimilarityMatrix().LeastSimilar()
			iterations := 0
			for val < threshold && iterations < len(e.Datasets()) {
				log.Println("Computing the similarities for ", idx, val)
				for j := 0; j < len(e.Datasets()); j++ {
					d1, d2 := e.Datasets()[idx], e.Datasets()[j]
					e.SimilarityMatrix().Set(idx, j, e.Similarity(d1, d2))
				}
				iterations++
				idx, val = e.SimilarityMatrix().LeastSimilar()
			}
		}
	}
	e.setDuration(time.Since(start).Seconds())
	return nil
}

func readDatasets(e DatasetSimilarityEstimator) error {
	e.setSimilarityMatrix(NewDatasetSimilarities(len(e.Datasets())))
	log.Println("Fetching datasets in memory")
	if e.Datasets() == nil || len(e.Datasets()) == 0 {
		log.Println("No datasets were given")
		return errors.New("Datasets not set correctly")
	}
	for _, d := range e.Datasets() {
		d.ReadFromFile()
	}
	return nil
}

// DatasetSimilarityEstimatorType represents the type of the Similarity Estimator
type DatasetSimilarityEstimatorType uint

const (
	// SimilarityTypeJaccard estimates the Jaccard coefficient
	SimilarityTypeJaccard DatasetSimilarityEstimatorType = iota
	// SimilarityTypeBhattacharyya estimates the Bhattacharyya coefficient
	SimilarityTypeBhattacharyya DatasetSimilarityEstimatorType = iota + 1
	// SimilarityTypeScript uses a script to transform the data
	SimilarityTypeScript DatasetSimilarityEstimatorType = iota + 2
	// SimilarityTypeComposite utilizes multiple estimators concurrently
	SimilarityTypeComposite DatasetSimilarityEstimatorType = iota + 4
	// SimilarityTypeCorrelation estimates correlation metrics
	SimilarityTypeCorrelation DatasetSimilarityEstimatorType = iota + 5
	// SimilarityTypeSize estimates size metric
	SimilarityTypeSize DatasetSimilarityEstimatorType = iota + 6
	// SimilarityTypeScriptPair estimates the similarity based on a script for each pair
	SimilarityTypeScriptPair DatasetSimilarityEstimatorType = iota + 7
)

// DatasetSimilarityEstimatorAvailableTypes lists the available similarity types
var DatasetSimilarityEstimatorAvailableTypes = []DatasetSimilarityEstimatorType{
	SimilarityTypeBhattacharyya,
	SimilarityTypeJaccard,
	SimilarityTypeCorrelation,
	SimilarityTypeComposite,
	SimilarityTypeScript,
	SimilarityTypeSize,
	SimilarityTypeScriptPair,
}

// NewDatasetSimilarityEstimatorType transforms the similarity type from a
// string to a DatasetSimilarityEstimatorType object
func NewDatasetSimilarityEstimatorType(estimatorType string) *DatasetSimilarityEstimatorType {
	lower := strings.ToLower(estimatorType)
	types := map[string]DatasetSimilarityEstimatorType{
		"bhattacharyya": SimilarityTypeBhattacharyya,
		"jaccard":       SimilarityTypeJaccard,
		"correlation":   SimilarityTypeCorrelation,
		"composite":     SimilarityTypeComposite,
		"script":        SimilarityTypeScript,
		"size":          SimilarityTypeSize,
		"scriptpair":    SimilarityTypeScriptPair,
	}
	if val, ok := types[lower]; ok {
		return &val
	}
	return nil
}
func (t DatasetSimilarityEstimatorType) String() string {
	if t == SimilarityTypeJaccard {
		return "Jaccard"
	} else if t == SimilarityTypeBhattacharyya {
		return "Bhattacharyya"
	} else if t == SimilarityTypeCorrelation {
		return "Correlation"
	} else if t == SimilarityTypeComposite {
		return "Composite"
	} else if t == SimilarityTypeScript {
		return "Script"
	} else if t == SimilarityTypeScriptPair {
		return "ScriptPair"
	} else if t == SimilarityTypeSize {
		return "Size"
	}
	return ""
}

// DatasetSimilarityPopulationPolicy is the struct that hold the Population
// Policy of the Similarity Matrix along with the configuration parameters of it.
type DatasetSimilarityPopulationPolicy struct {
	PolicyType DatasetSimilarityPopulationPolicyType
	Parameters map[string]float64
}

// Serialize method returns a slice of bytes containing the serialized form of
// the Population Policy
func (s *DatasetSimilarityPopulationPolicy) Serialize() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(getBytesInt(int(s.PolicyType)))
	buffer.Write(getBytesInt(len(s.Parameters)))
	for k, v := range s.Parameters {
		buffer.WriteString(k + "\n")
		buffer.Write(getBytesFloat(v))
	}
	return buffer.Bytes()
}

// Deserialize is responsible to instantiate a Population Policy object based on
// its byte representation.
func (s *DatasetSimilarityPopulationPolicy) Deserialize(b []byte) {
	buffer := bytes.NewBuffer(b)
	tempInt := make([]byte, 4)
	tempFloat := make([]byte, 8)
	buffer.Read(tempInt)
	s.PolicyType = DatasetSimilarityPopulationPolicyType(getIntBytes(tempInt))
	buffer.Read(tempInt)
	count := getIntBytes(tempInt)
	if s.Parameters == nil {
		s.Parameters = make(map[string]float64)
	}
	for i := 0; i < count; i++ {
		line, _ := buffer.ReadString('\n')
		line = strings.TrimSpace(line)
		buffer.Read(tempFloat)
		val := getFloatBytes(tempFloat)
		s.Parameters[line] = val
	}

}

// DatasetSimilarityPopulationPolicyType is the type that represents the
// Similarity Matrix population policy
type DatasetSimilarityPopulationPolicyType uint

const (
	// PopulationPolicyFull policy needs no params
	PopulationPolicyFull DatasetSimilarityPopulationPolicyType = iota
	// PopulationPolicyAprx must have defined one of two params:
	// count (how many points) or threshold (percentage in similarity gain)
	PopulationPolicyAprx DatasetSimilarityPopulationPolicyType = iota + 1
)

// NewDatasetSimilarityEstimator is a factory method for the
// DatasetSimilarityEstimator structs, used to initialize the estimator and
// return it to the user.
func NewDatasetSimilarityEstimator(
	estType DatasetSimilarityEstimatorType,
	datasets []*Dataset) DatasetSimilarityEstimator {
	policy := *new(DatasetSimilarityPopulationPolicy)
	policy.PolicyType = PopulationPolicyFull
	if estType == SimilarityTypeJaccard {
		a := new(JaccardEstimator)
		a.SetPopulationPolicy(policy)
		a.datasets = datasets
		return a
	} else if estType == SimilarityTypeBhattacharyya {
		a := new(BhattacharyyaEstimator)
		a.SetPopulationPolicy(policy)
		a.datasets = datasets
		return a
	} else if estType == SimilarityTypeScript {
		a := new(ScriptSimilarityEstimator)
		a.SetPopulationPolicy(policy)
		a.datasets = datasets
		return a
	} else if estType == SimilarityTypeCorrelation {
		a := new(CorrelationEstimator)
		a.SetPopulationPolicy(policy)
		a.datasets = datasets
		return a
	} else if estType == SimilarityTypeComposite {
		a := new(CompositeEstimator)
		a.SetPopulationPolicy(policy)
		a.datasets = datasets
		return a
	} else if estType == SimilarityTypeSize {
		a := new(SizeEstimator)
		a.SetPopulationPolicy(policy)
		a.datasets = datasets
		return a
	} else if estType == SimilarityTypeScriptPair {
		a := new(ScriptPairSimilarityEstimator)
		a.SetPopulationPolicy(policy)
		a.datasets = datasets
		return a
	}

	return nil
}

// DeserializeSimilarityEstimator method is used to deserialize the
// Estimator according to its type
func DeserializeSimilarityEstimator(b []byte) DatasetSimilarityEstimator {
	estimatorType := DatasetSimilarityEstimatorType(getIntBytes(b[0:4]))
	if estimatorType == SimilarityTypeJaccard {
		a := new(JaccardEstimator)
		a.Deserialize(b)
		return a
	} else if estimatorType == SimilarityTypeBhattacharyya {
		a := new(BhattacharyyaEstimator)
		a.Deserialize(b)
		return a
	} else if estimatorType == SimilarityTypeComposite {
		a := new(CompositeEstimator)
		a.Deserialize(b)
		return a
	} else if estimatorType == SimilarityTypeCorrelation {
		a := new(CorrelationEstimator)
		a.Deserialize(b)
		return a
	} else if estimatorType == SimilarityTypeScript {
		a := new(ScriptSimilarityEstimator)
		a.Deserialize(b)
		return a
	} else if estimatorType == SimilarityTypeSize {
		a := new(SizeEstimator)
		a.Deserialize(b)
		return a
	}
	return nil
}

// DatasetSimilarityMatrix represent the struct that holds the results of  a
// dataset similarity estimation. It also provides the necessary
type DatasetSimilarityMatrix struct {
	// the actual similarities holder
	similarities [][]float64
	// indicates whether the closestIndex is disabled or not
	indexDisabled bool
	// index that hold the closest datasets
	closestIndex *closestIndex
	// represents the capacity of the sim matrix
	capacity int
}

// NewDatasetSimilarities is the constructor for the DatasetSimilarities struct,
// expecting the number of datasets that will be held by it. If capacity=0, this
// implies that the Similarity Matrix will be deserialzed.
func NewDatasetSimilarities(capacity int) *DatasetSimilarityMatrix {
	r := new(DatasetSimilarityMatrix)
	r.indexDisabled = false
	r.capacity = capacity
	if capacity != 0 {
		r.allocateStructs()
	}
	return r
}

// IndexDisabled sets whether the closest dataset index should be disabled or not.
// The index is useless if the FULL Estimator strategy is being followed.
func (s *DatasetSimilarityMatrix) IndexDisabled(flag bool) {
	s.indexDisabled = flag
}

// FullyCalculatedNodes returns the number of nodes the similarity of which
// has been calculated for all the nodes. This number can work as a measure
// of how close to the full similarity matrix the current object is.
func (s *DatasetSimilarityMatrix) FullyCalculatedNodes() int {
	if s.indexDisabled {
		return s.capacity
	}
	count := 0
	for i := 0; i < s.capacity; i++ {
		if idx, _ := s.closestIndex.Get(i); idx == i {
			count++
		}
	}
	return count
}

func (s *DatasetSimilarityMatrix) allocateStructs() {
	s.similarities = make([][]float64, s.capacity-1)
	for i := 0; i < s.capacity-1; i++ {
		s.similarities[i] = make([]float64, s.capacity-i-1)
	}
	s.closestIndex = newClosestIndex(s.capacity)
}

// Capacity returns the capacity of the Similarity Matrix
func (s *DatasetSimilarityMatrix) Capacity() int {
	return s.capacity
}

// Set is a setter function for the similarity between two datasets
func (s *DatasetSimilarityMatrix) Set(idxA, idxB int, value float64) {
	if idxA == idxB { // do nothing
		if !s.indexDisabled {
			s.closestIndex.CheckAndSet(idxA, idxB, value)
			s.closestIndex.CheckAndSet(idxB, idxA, value)
		}
		return
	} else if idxA > idxB { //we only want to fill the upper diagonal elems
		t := idxB
		idxB = idxA
		idxA = t
	}
	s.similarities[idxA][idxB-idxA-1] = value
	if !s.indexDisabled {
		s.closestIndex.CheckAndSet(idxA, idxB, value)
		s.closestIndex.CheckAndSet(idxB, idxA, value)
	}
}

// Get returns the similarity between two dataset paths
func (s *DatasetSimilarityMatrix) Get(idxA, idxB int) float64 {
	if !s.indexDisabled {
		idxA, _ = s.closestIndex.Get(idxA)
		idxB, _ = s.closestIndex.Get(idxB)
	}
	if idxA == idxB {
		return 1.0
	} else if idxA > idxB {
		t := idxB
		idxB = idxA
		idxA = t
	}
	return s.similarities[idxA][idxB-idxA-1]
}

// LeastSimilar method returns the dataset that presents the lowest
// similarity among the examined datasets
func (s *DatasetSimilarityMatrix) LeastSimilar() (int, float64) {
	return s.closestIndex.LeastSimilar()
}

func (s DatasetSimilarityMatrix) String() string {
	var buf bytes.Buffer
	for i := 0; i < s.capacity; i++ {
		for j := 0; j < s.capacity; j++ {
			buf.WriteString(fmt.Sprintf("%.5f ", s.Get(i, j)))
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

// Serialize method returns a byte slice that represents the similarity matrix
func (s *DatasetSimilarityMatrix) Serialize() []byte {
	buf := new(bytes.Buffer)
	buf.Write(getBytesInt(s.capacity))
	for i := 0; i < s.capacity-1; i++ {
		for j := 0; j < s.capacity-1-i; j++ {
			buf.Write(getBytesFloat(s.similarities[i][j]))
		}
	}

	for i := 0; i < s.capacity; i++ {
		buf.Write(getBytesFloat(float64(s.closestIndex.closestIdx[i])))
		buf.Write(getBytesFloat(s.closestIndex.similarity[i]))
	}
	charBuf := make([]byte, 1)
	if s.indexDisabled {
		charBuf[0] = 1
	} else {
		charBuf[0] = 0
	}
	buf.Write(charBuf)

	// compress before you send
	var compressed bytes.Buffer
	wr, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	//	wr := gzip.NewWrite(&compressed)
	defer wr.Close()
	if err != nil {
		log.Println("Error message from compression: ", err)
	}
	wr.Write(buf.Bytes())
	wr.Flush()
	log.Println("Compressed bytes", len(compressed.Bytes()), len(buf.Bytes()))
	return compressed.Bytes()
	//	return buf.Bytes()
}

// Deserialize instantiates an empty DatasetSimilarities object. In case of
// parse failure, an error is thrown
func (s *DatasetSimilarityMatrix) Deserialize(buff []byte) error {
	// decompress stream
	re, err := gzip.NewReader(bytes.NewBuffer(buff))
	if err != nil {
		log.Println("Error message from compression: ", err)
	}
	defer re.Close()
	buff, err = ioutil.ReadAll(re)
	if err != nil && err != io.ErrUnexpectedEOF {
		log.Println(err, len(buff))
	}

	buf := bytes.NewBuffer(buff)
	tempInt := make([]byte, 4)
	buf.Read(tempInt)
	s.capacity = getIntBytes(tempInt)
	s.allocateStructs()

	tempFloat := make([]byte, 8)
	for i := 0; i < s.capacity-1; i++ {
		for j := 0; j < s.capacity-1-i; j++ {
			buf.Read(tempFloat)
			s.similarities[i][j] = getFloatBytes(tempFloat)
		}
	}

	for i := 0; i < s.capacity; i++ {
		buf.Read(tempFloat)
		s.closestIndex.closestIdx[i] = int(getFloatBytes(tempFloat))
		buf.Read(tempFloat)
		s.closestIndex.similarity[i] = getFloatBytes(tempFloat)
	}
	charBuf := make([]byte, 1)
	buf.Read(charBuf)
	if charBuf[0] == 1 {
		s.indexDisabled = true
	} else {
		s.indexDisabled = false
	}
	return nil
}

// datasets along with their respective similarities
type closestIndex struct {
	closestIdx []int
	similarity []float64
}

func newClosestIndex(datasets int) *closestIndex {
	res := new(closestIndex)
	res.closestIdx = make([]int, datasets)
	res.similarity = make([]float64, datasets)
	for i := range res.closestIdx {
		res.closestIdx[i] = -1   // represents a NUL dataset index
		res.similarity[i] = -1.0 // represents a NUL similarity
	}
	return res
}

// Returns the index and similarity of the most similar dataset
func (s *closestIndex) Get(idx int) (int, float64) {
	if idx < len(s.closestIdx) {
		return s.closestIdx[idx], s.similarity[idx]
	}
	return -1, -1.0
}

// Sets the index and the similarity of the most similar dataset
func (s *closestIndex) Set(srcIdx, dstIdx int, similarity float64) {
	s.closestIdx[srcIdx] = dstIdx
	s.similarity[srcIdx] = similarity
}

// Sets the index and the similarity of the most similar dataset,
// iff the provided similarity is higher than the one previously stored
func (s *closestIndex) CheckAndSet(srcIdx, dstIdx int, similarity float64) {
	if s.similarity[srcIdx] < similarity {
		s.closestIdx[srcIdx] = dstIdx
		s.similarity[srcIdx] = similarity
	}
}

// Returns the dataset index (and its respective value) with the lowest
// similarity to its most close dataset
func (s *closestIndex) LeastSimilar() (int, float64) {
	minIdx, minV := 0, s.similarity[0]
	for i, v := range s.similarity {
		if v < minV {
			minV = v
			minIdx = i
		} else if v == minV && rand.Int()%2 == 0 { // random index
			minIdx = i
		}
	}
	return minIdx, minV
}
