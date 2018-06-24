package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ModelerType uint8

const (
	ScriptBasedModelerType ModelerType = iota
	KNNModelerType         ModelerType = iota + 1
)

func NewModelerType(t string) ModelerType {
	if strings.ToLower(t) == "script" {
		return ScriptBasedModelerType
	} else {
		return KNNModelerType
	}
}

// NewModeler is the factory method for the modeler object
func NewModeler(
	modelerType ModelerType,
	datasets []*Dataset,
	sr float64,
	evaluator DatasetEvaluator) Modeler {
	if modelerType == ScriptBasedModelerType {
		modeler := new(ScriptBasedModeler)
		modeler.datasets = datasets
		modeler.samplingRate = sr
		modeler.evaluator = evaluator
		return modeler
	} else if modelerType == KNNModelerType {
		modeler := new(KNNModeler)
		modeler.datasets = datasets
		modeler.samplingRate = sr
		modeler.evaluator = evaluator
		return modeler
	}
	return nil
}

// Modeler is the interface for the objects that model the dataset space.
type Modeler interface {
	//  Configure is responsible to provide the necessary configuration
	// options to the Modeler struct. Call it before Run.
	Configure(map[string]string) error
	// Run initiates the modeling process.
	Run() error

	// Datasets returns the datasets slice
	Datasets() []*Dataset
	// Samples returns the indices of the chosen datasets.
	Samples() map[int]float64
	// AppxValues returns a slice of the approximated values
	AppxValues() []float64

	// ErrorMetrics returns a list of error metrics for the specified modeler
	ErrorMetrics() map[string]float64

	// ExecTime returns the total execution time of the Modeler
	ExecTime() float64
	// EvalTime returns the evaluation time of the Modeler
	EvalTime() float64
}

// AbstractModeler implements the common methods of the Modeler structs
type AbstractModeler struct {
	datasets  []*Dataset       // the datasets the modeler refers to
	evaluator DatasetEvaluator // the evaluator struct that gets the values

	samplingRate float64         // the portion of the datasets to examine
	samples      map[int]float64 // the dataset indices chosen for samples
	appxValues   []float64       // the appx values of ALL the datasets

	execTime float64 // the total time in seconds
	evalTime float64 // the time needed to evaluate the datasets in seconds
}

// Datasets returns the datasets slice
func (a *AbstractModeler) Datasets() []*Dataset {
	return a.datasets
}

// Samples return the indices of the chosen datasets
func (a *AbstractModeler) Samples() map[int]float64 {
	return a.samples
}

// AppxValues returns the values of all the datasets
func (a *AbstractModeler) AppxValues() []float64 {
	return a.appxValues
}

// ErrorMetrics returns a list of error metrics for the specified model
func (a *AbstractModeler) ErrorMetrics() map[string]float64 {
	if a.appxValues == nil || len(a.appxValues) == 0 {
		return nil
	}
	errors := make(map[string]float64)

	var actual []float64
	for _, d := range a.datasets {
		val, err := a.evaluator.Evaluate(d.Path())
		if err != nil {
			log.Println(err)
			actual = append(actual, math.NaN())
		} else {
			actual = append(actual, val)
		}
	}

	// evaluation for the entire dataset
	allIndices := make([]int, len(actual))
	for i := range actual {
		allIndices[i] = i
	}
	for k, v := range a.getMetrics(allIndices, actual, "all") {
		errors[k] = v
	}

	// evaluation for the unknown datasets
	var unknownIndices []int
	for i := range actual {
		if _, ok := a.samples[i]; !ok {
			unknownIndices = append(unknownIndices, i)
		}
	}
	for k, v := range a.getMetrics(unknownIndices, actual, "unknown") {
		errors[k] = v
	}

	tempErrors := make(map[string][]float64)
	for r := 0; r < 10; r++ {
		perm := rand.Perm(len(unknownIndices))
		var idcs []int
		for _, i := range perm {
			idcs = append(idcs, unknownIndices[i])
		}
		for _, s := range []float64{0.05, 0.10, 0.20, 0.30, 0.40} {
			newlength := int(float64(len(actual)) * s)
			if newlength < len(idcs) && newlength > 0 {
				for k, v := range a.getMetrics(idcs[1:newlength], actual, fmt.Sprintf("%02.0f%%", 100*s)) {
					//errors[k] = v
					if _, ok := tempErrors[k]; !ok {
						tempErrors[k] = make([]float64, 0)
					}
					tempErrors[k] = append(tempErrors[k], v)
				}
			}
		}
	}
	for k, v := range tempErrors {
		errors[k] = Mean(v)
	}

	return errors
}

func (a *AbstractModeler) getMetrics(testIdx []int, actual []float64, label string) map[string]float64 {
	actualUnknown, appxUnknown, residualsUnknown :=
		make([]float64, len(testIdx)),
		make([]float64, len(testIdx)),
		make([]float64, len(testIdx))
	maxValue := math.NaN()
	for i, v := range testIdx {
		if math.IsNaN(maxValue) || maxValue < actual[v] {
			maxValue = actual[v]
		}
		actualUnknown[i] = actual[v]
		appxUnknown[i] = a.appxValues[v]
		residualsUnknown[i] = math.Abs(actualUnknown[i] - appxUnknown[i])
	}
	errors := make(map[string]float64)
	errors["RMSE-"+label] = RootMeanSquaredError(actualUnknown, appxUnknown)
	errors["NRMSE-"+label] = RootMeanSquaredError(actualUnknown, appxUnknown) / maxValue
	errors["RMSLE-"+label] = RootMeanSquaredLogError(actualUnknown, appxUnknown)
	errors["MAPE-"+label] = MeanAbsolutePercentageError(actualUnknown, appxUnknown)
	errors["MdAPE-"+label] = MedianAbsolutePercentageError(actualUnknown, appxUnknown)
	errors["MAE-"+label] = MeanAbsoluteError(actualUnknown, appxUnknown)
	errors["R^2-"+label] = RSquared(actualUnknown, appxUnknown)
	errors["Res000-"+label] = Percentile(residualsUnknown, 0)
	errors["Res025-"+label] = Percentile(residualsUnknown, 25)
	errors["Res050-"+label] = Percentile(residualsUnknown, 50)
	errors["Res075-"+label] = Percentile(residualsUnknown, 75)
	errors["Res100-"+label] = Percentile(residualsUnknown, 100)
	errors["Kendall-"+label] = Kendall(actualUnknown, appxUnknown)
	return errors
}

// ExecTime returns the total exection time of the Modeler
func (a *AbstractModeler) ExecTime() float64 {
	return a.execTime
}

// EvalTime returns the dataset evaluation time of the Modeler
func (a *AbstractModeler) EvalTime() float64 {
	return a.evalTime
}

func (m *AbstractModeler) deploySamples() {
	s := int(math.Floor(m.samplingRate * float64(len(m.datasets))))
	// sample the datasets
	permutation := rand.Perm(len(m.datasets))
	m.samples = make(map[int]float64)
	// deploy samples
	for i := 0; i < len(permutation) && (len(m.samples) < s); i++ {
		idx := permutation[i]
		start2 := time.Now()
		val, err := m.evaluator.Evaluate(m.datasets[idx].Path())
		m.evalTime += (time.Since(start2).Seconds())
		if err != nil {
			log.Printf("%s: %s\n", m.datasets[idx].Path(), err.Error())
		} else {
			m.samples[idx] = val
		}
	}
}

// ScriptBasedModeler utilizes a script to train an ML model and obtain is values
type ScriptBasedModeler struct {
	AbstractModeler
	script      string               // the script to use for modeling
	coordinates []DatasetCoordinates // the dataset coordinates
}

// Configure expects the necessary conf options for the specified struct.
// Specifically, the following parameters are necessary:
// - script: the path of the script to use
func (m *ScriptBasedModeler) Configure(conf map[string]string) error {
	if val, ok := conf["script"]; ok {
		m.script = val
	} else {
		log.Println("script parameter is missing")
		return errors.New("script parameter is missing")
	}

	if val, ok := conf["coordinates"]; ok {
		buf, err := ioutil.ReadFile(val)
		if err != nil {
			log.Println(err)
		}
		m.coordinates = DeserializeCoordinates(buf)
	} else {
		log.Println("coordinates parameter is missing")
		return errors.New("coordinates parameter is missing")
	}
	return nil
}

// Run executes the modeling process and populates the samples, realValues and
// appxValues slices.
func (m *ScriptBasedModeler) Run() error {
	start := time.Now()
	m.deploySamples()

	var trainingSet, testSet [][]float64
	for idx, val := range m.samples {
		trainingSet = append(trainingSet, append(m.coordinates[idx], val))
	}
	trainFile := createCSVFile(trainingSet, true)
	for _, v := range m.coordinates {
		testSet = append(testSet, v)
	}
	testFile := createCSVFile(testSet, false)
	appx, err := m.executeMLScript(trainFile, testFile)
	if err != nil {
		return err
	}
	m.appxValues = appx
	os.Remove(trainFile)
	os.Remove(testFile)
	m.execTime = time.Since(start).Seconds()
	return nil
}

// executeMLScript executes the ML script, utilizing the selected samples (indices)
// and populates the real and appx values slices
func (m *ScriptBasedModeler) executeMLScript(trainFile, testFile string) ([]float64, error) {
	var result []float64
	cmd := exec.Command(m.script, trainFile, testFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New(err.Error() + string(out))
	}
	outputString := string(out)
	array := strings.Split(outputString, "\n")
	result = make([]float64, len(m.datasets))
	for i := 0; i < len(m.datasets); i++ {
		val, err := strconv.ParseFloat(array[i], 64)
		if err != nil {
			log.Println(err)
		} else {
			result[i] = val
		}
	}
	return result, nil
}

// KNNModeler utilizes a similarity matrix in order to approximate the training set
type KNNModeler struct {
	AbstractModeler
	k  int                      // the number of neighbors to check
	sm *DatasetSimilarityMatrix // the similarity matrix
}

// Configure is the method used to provide the essential paremeters for the conf of the modeler
func (m *KNNModeler) Configure(conf map[string]string) error {
	if val, ok := conf["k"]; ok {
		intVal, err := strconv.ParseInt(val, 10, 32)
		m.k = int(intVal)
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println("No k parameter provided")
		return errors.New("No k parameter provided")
	}

	if val, ok := conf["smatrix"]; ok {
		buf, err := ioutil.ReadFile(val)
		if err != nil {
			log.Println(err)
			return err
		}
		m.sm = new(DatasetSimilarityMatrix)
		m.sm.Deserialize(buf)
	} else {
		log.Println("No smatrix parameter provided")
		return errors.New("No smatrix parameter provided")
	}

	return nil
}

// Run executes the training part and obtains the model
func (k *KNNModeler) Run() error {
	start := time.Now()
	k.deploySamples()
	k.appxValues = make([]float64, len(k.datasets))
	for i := range k.datasets {
		if _, ok := k.samples[i]; ok {
			k.appxValues[i] = k.samples[i]
		} else {
			k.appxValues[i] = k.approximateValue(i)
		}
	}

	k.execTime = time.Since(start).Seconds()
	return nil
}

func (k *KNNModeler) approximateValue(id int) float64 {
	var pList pairList
	for j := range k.samples {
		s := k.sm.Get(id, j)
		pList = append(pList, pair{j, s})
	}
	sort.Sort(sort.Reverse((pList)))
	weights := 0.0
	values := 0.0
	for i := 0; i < len(pList) && i < k.k; i++ {
		p := pList[i]
		values += p.Similarity * k.samples[p.Id]
		weights += p.Similarity
	}
	return values / weights
}

type pair struct {
	Id         int
	Similarity float64
}

type pairList []pair

func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Similarity < p[j].Similarity }
func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// createCSVFile serializes a double float slice to a CSV file and returns
// the filename
func createCSVFile(matrix [][]float64, output bool) string {
	f, err := ioutil.TempFile("/tmp", "csv")
	if err != nil {
		log.Println(err)
	}
	cols := 0
	if len(matrix) > 0 {
		cols = len(matrix[0])
	}
	if output {
		cols--
	}

	for i := 1; i < cols+1; i++ {
		fmt.Fprintf(f, "x%d", i)
		if i < cols {
			fmt.Fprintf(f, ",")
		}
	}
	if output {
		fmt.Fprintf(f, ",class")
	}
	fmt.Fprintf(f, "\n")

	for i := range matrix {
		for j := range matrix[i] {
			fmt.Fprintf(f, "%.5f", matrix[i][j])
			if j < len(matrix[i])-1 {
				fmt.Fprintf(f, ",")
			}
		}
		fmt.Fprintf(f, "\n")
	}
	f.Close()
	return f.Name()
}

// RootMeanSquaredError returns the RMSE of the actual vs the predicted values
func RootMeanSquaredError(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) || len(actual) == 0 {
		log.Printf("actual (%d) and predicted (%d) are of different size", len(actual), len(predicted))
		return math.NaN()
	}
	sum := 0.0
	count := 0.0
	for i := range actual {
		if !math.IsNaN(actual[i]) {
			diff := actual[i] - predicted[i]
			sum += diff * diff
			count += 1
		}
	}
	if count > 0 {
		return math.Sqrt(sum / count)
	}
	return math.NaN()
}

// RootMeanSquaredLogError returns the RMSLE of the actual vs the predicted values
func RootMeanSquaredLogError(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) || len(actual) == 0 {
		log.Println("actual and predicted values are of different size!!")
		return math.NaN()
	}
	sum := 0.0
	count := 0.0
	for i := range actual {
		if !math.IsNaN(actual[i]) && actual[i] > -1 && predicted[i] > -1 {
			diff := math.Log(predicted[i]+1) - math.Log(actual[i]+1)
			sum += diff * diff
			count += 1
		}
	}
	if count > 0 {
		return math.Sqrt(sum / count)
	}
	return math.NaN()
}

// MeanAbsoluteError returns the MAE of the actual vs the predicted values
func MeanAbsoluteError(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) || len(actual) == 0 {
		log.Println("actual and predicted values are of different size!!")
		return math.NaN()
	}
	sum := 0.0
	count := 0.0
	for i := range actual {
		if actual[i] != 0.0 && !math.IsNaN(actual[i]) {
			count += 1.0
			sum += math.Abs((actual[i] - predicted[i]))
		}
	}
	if count > 0 {
		return sum / count
	}
	return math.NaN()
}

// MedianAbsolutePercentageError returns the MdAPE of the actual vs the predicted values
func MedianAbsolutePercentageError(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) || len(actual) == 0 {
		log.Println("actual and predicted values are of different size!!")
		return math.NaN()
	}
	apes := make([]float64, 0)
	for i := range actual {
		if actual[i] != 0.0 && !math.IsNaN(actual[i]) {
			val := math.Abs((actual[i] - predicted[i]) / actual[i])
			apes = append(apes, val)
		}
	}
	if len(apes) > 0 {
		return Percentile(apes, 50)
	}
	return math.NaN()
}

// MeanAbsolutePercentageError returns the MAPE of the actual vs the predicted values
func MeanAbsolutePercentageError(actual, predicted []float64) float64 {
	if len(actual) != len(predicted) || len(actual) == 0 {
		log.Println("actual and predicted values are of different size!!")
		return math.NaN()
	}
	sum := 0.0
	count := 0.0
	for i := range actual {
		if actual[i] != 0.0 && !math.IsNaN(actual[i]) {
			count += 1.0
			val := math.Abs((actual[i] - predicted[i]) / actual[i])
			sum += val
		}
	}
	if count > 0 {
		return sum / count
	}
	return math.NaN()
}

// RSquared returns the coeff. of determination of the actual vs the predicted values
func RSquared(actual, predicted []float64) float64 {
	if len(predicted) != len(actual) || len(predicted) == 0 {
		log.Println("actual and predicted values are of different size!!")
		return math.NaN()
	}
	mean := Mean(actual)
	ssRes, ssTot := 0.0, 0.0
	for i := range actual {
		if !math.IsNaN(actual[i]) {
			ssTot += (actual[i] - mean) * (actual[i] - mean)
			ssRes += (actual[i] - predicted[i]) * (actual[i] - predicted[i])
		}
	}
	if ssTot > 0 {
		return 1.0 - (ssRes / ssTot)
	}
	return math.NaN()
}

// Percentile returns the i-th percentile of an array of values
func Percentile(values []float64, percentile int) float64 {
	valuesCopy := make([]float64, len(values))
	copy(valuesCopy, values)
	if !sort.Float64sAreSorted(valuesCopy) {
		sort.Float64s(valuesCopy)
	}
	idx := int(math.Ceil((float64(percentile) / 100.0) * float64(len(valuesCopy))))
	if idx < len(valuesCopy) {
		return valuesCopy[idx]
	}
	if len(valuesCopy) > 0 {
		return valuesCopy[len(valuesCopy)-1]
	}
	return math.NaN()
}
