package core

import (
	"bytes"
	"log"
	"math"
	"strconv"
	"strings"
)

// CorrelationEstimator estimates the similarity between two datasets based
// on a correlation metric. This metric can only be used for datasets that
// consist of a single column and consist of the same number of tuples.
type CorrelationEstimator struct {
	AbstractDatasetSimilarityEstimator
	estType  CorrelationEstimatorType
	normType CorrelationEstimatorNormalizationType
	column   int
}

// CorrelationEstimatorType represents the type of correlation to be used
// by the CorrelationEstimator
type CorrelationEstimatorType uint8

const (
	// CorrelationSimilarityTypePearson represents the Pearson cor. coeff
	CorrelationSimilarityTypePearson = iota
	// CorrelationSimilarityTypeSpearman represents the Spearman cor. coeff
	CorrelationSimilarityTypeSpearman = iota + 1
	// CorrelationSimilarityTypeKendall represents the Kendall cor. coeff
	CorrelationSimilarityTypeKendall = iota + 2
)

// String returns a string representation of the CorrelationEstimatorType
func (s CorrelationEstimatorType) String() string {
	if s == CorrelationSimilarityTypePearson {
		return "pearson"
	} else if s == CorrelationSimilarityTypeSpearman {
		return "spearman"
	} else if s == CorrelationSimilarityTypeKendall {
		return "kendall"
	}
	return ""
}

// CorrelationEstimatorNormalizationType represents the type of the normalization
// action. Since all correlation metrics can take any valuein [-1,1], this
// type reflects the policy with which [-1,1] will be mapped to a similarity
// metric in [0,1]
type CorrelationEstimatorNormalizationType uint8

const (
	// CorrelationSimilarityNormalizationAbs returns |r|, r being the cor. metric
	CorrelationSimilarityNormalizationAbs = iota
	// CorrelationSimilarityNormalizationScale returns r/2 + 0.5, r being the cor. metric
	CorrelationSimilarityNormalizationScale = iota + 1
	// CorrelationSimilarityNormalizationPos returns r, if r>=0 else 0
	CorrelationSimilarityNormalizationPos = iota + 2
)

// String returns a nstring representation of the CorrelationEstimatorNormalizationType
func (s CorrelationEstimatorNormalizationType) String() string {
	if s == CorrelationSimilarityNormalizationAbs {
		return "abs"
	} else if s == CorrelationSimilarityNormalizationScale {
		return "scale"
	} else if s == CorrelationSimilarityNormalizationPos {
		return "pos"
	}
	return ""
}

// Configure provides a set of configuration options to the CorrelationEstimator
// struct.
func (e *CorrelationEstimator) Configure(conf map[string]string) {
	if val, ok := conf["concurrency"]; ok {
		conv, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			log.Println(err)
		} else {
			e.concurrency = int(conv)
			log.Println("Set concurrency to", e.concurrency)
		}
	} else {
		e.concurrency = 1
	}
	if val, ok := conf["column"]; ok {
		conv, err := strconv.ParseInt(val, 10, 32)
		if err != nil {
			log.Println(err)
		} else {
			e.column = int(conv)
			log.Println("Set column to", e.column)
		}
	} else {
		e.column = 0
	}
	if val, ok := conf["correlation"]; ok {
		if strings.ToLower(val) == "pearson" {
			e.estType = CorrelationSimilarityTypePearson
		} else if strings.ToLower(val) == "spearman" {
			e.estType = CorrelationSimilarityTypeSpearman
		} else if strings.ToLower(val) == "kendall" {
			e.estType = CorrelationSimilarityTypeKendall
		}
		log.Println("Set correlation type to", e.estType)
	} else {
		e.estType = CorrelationSimilarityTypePearson
	}
	if val, ok := conf["normalization"]; ok {
		if strings.ToLower(val) == "abs" {
			e.normType = CorrelationSimilarityNormalizationAbs
		} else if strings.ToLower(val) == "scale" {
			e.normType = CorrelationSimilarityNormalizationScale
		} else if strings.ToLower(val) == "pos" {
			e.normType = CorrelationSimilarityNormalizationPos
		}
		log.Println("Set normalization to", e.normType)
	} else {
		e.normType = CorrelationSimilarityNormalizationPos
	}
}

// Options returns a list of options used internally by the CorrelationEstimator
// struct for its execution.
func (e *CorrelationEstimator) Options() map[string]string {
	return map[string]string{
		"concurrency":   "max number of threads to use",
		"correlation":   "one of [Pearson], Spearman, Kendall",
		"column":        "number of column of the datasets to consider - starting from 0 (default)",
		"normalization": "determines how to scale the correlation metric from [-1,1]-> [0,1]. one of: abs scale [pos]",
	}
}

// Serialize returns a byte array in order to serialize the CorrelationEstimator
// struct.
func (e *CorrelationEstimator) Serialize() []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(getBytesInt(int(SimilarityTypeCorrelation)))
	buffer.Write(
		datasetSimilarityEstimatorSerialize(
			e.AbstractDatasetSimilarityEstimator))
	buffer.Write(getBytesInt(int(e.estType)))
	return buffer.Bytes()
}

// Deserialize returns a byte array in order to deserialize the CorrelationEstimator
func (e *CorrelationEstimator) Deserialize(b []byte) {
	buffer := bytes.NewBuffer(b)
	tempInt := make([]byte, 4)
	buffer.Read(tempInt) // consume estimator type
	buffer.Read(tempInt)
	absEstBytes := make([]byte, getIntBytes(tempInt))
	buffer.Read(absEstBytes)
	e.AbstractDatasetSimilarityEstimator =
		*datasetSimilarityEstimatorDeserialize(absEstBytes)

	buffer.Read(tempInt)
	e.estType = CorrelationEstimatorType(getIntBytes(tempInt))
}

// Compute method constructs the Similarity Matrix
func (e *CorrelationEstimator) Compute() error {
	return datasetSimilarityEstimatorCompute(e)
}

// Similarity returns the similarity between two datasets. Since all the
// correlation coefficients are between [-1.0,1.0], the output of this function
// is scaled to [0.0,1.0] by returning (x/2.0 + 0.5), where x is one of
// Pearson, Spearman and Kendall coefficients.
func (e *CorrelationEstimator) Similarity(a, b *Dataset) float64 {
	aTrans, bTrans := e.transformDataset(a), e.transformDataset(b)
	var val float64
	if e.estType == CorrelationSimilarityTypePearson {
		val = Pearson(aTrans, bTrans)
	} else if e.estType == CorrelationSimilarityTypeSpearman {
		val = Spearman(aTrans, bTrans)
	} else if e.estType == CorrelationSimilarityTypeKendall {
		val = Kendall(aTrans, bTrans)
	}
	return e.scaleCorrelationValue(val)
}

func (e *CorrelationEstimator) transformDataset(d *Dataset) []float64 {
	var result []float64
	for _, t := range d.Data() {
		if e.column < len(t.Data) {
			result = append(result, t.Data[e.column])
		} else {
			log.Printf("Given column number (%d) exceeds data columns (%d)\n", e.column, len(t.Data))
		}
	}
	return result
}

func (e *CorrelationEstimator) scaleCorrelationValue(val float64) float64 {
	if e.normType == CorrelationSimilarityNormalizationAbs {
		return math.Abs(val)
	} else if e.normType == CorrelationSimilarityNormalizationScale {
		return (val / 2.0) + 0.5
	} else if e.normType == CorrelationSimilarityNormalizationPos {
		if val > 0 {
			return val
		}
		return 0
	}
	log.Println("Unknown correlation scaling method")
	return -1.0
}

// Pearson returns the pearson correlation coefficient between two variables.
// The two arrays must be of the same size, else 0 is returned.
func Pearson(a, b []float64) float64 {
	aMean, bMean := Mean(a), Mean(b)
	if len(a) != len(b) {
		log.Println("Dataset sizes are different")
		return 0.0
	}
	nom, denomA, denomB := 0.0, 0.0, 0.0
	for i := range a {
		nom += (a[i] - aMean) * (b[i] - bMean)
		denomA += (a[i] - aMean) * (a[i] - aMean)
		denomB += (b[i] - bMean) * (b[i] - bMean)
	}
	if (denomA * denomB) != 0 {
		return nom / (math.Sqrt(denomA) * math.Sqrt(denomB))
	}
	log.Println("Denominator is zero")
	return .0
}

// Spearman returns the rank correlation coefficient between two variables.
// The two arrays must be of the same size, else 0 is returned.
func Spearman(a, b []float64) float64 {
	aRanks, bRanks := Rank(a), Rank(b)
	if len(a) != len(b) {
		log.Println("Dataset sizes are different")
		return 0.0
	}
	// stupid casting is required to use Pearson
	aRanksFloat := make([]float64, len(aRanks))
	bRanksFloat := make([]float64, len(bRanks))
	for i := range aRanks {
		aRanksFloat[i] = float64(aRanks[i])
		bRanksFloat[i] = float64(bRanks[i])
	}
	return Pearson(aRanksFloat, bRanksFloat)
}

// Kendall returns yet another rank correlation coefficient between two
// variables. The two arrays must be of the same size, else 0 is returned.
func Kendall(a, b []float64) float64 {
	aRanks, bRanks := Rank(a), Rank(b)
	if len(a) != len(b) {
		log.Println("Dataset sizes are different")
		return 0.0
	}
	concordant, discordant := 0, 0
	size := len(aRanks)
	for i := 0; i < size; i++ {
		for j := i; j < size; j++ {
			if (aRanks[i]-aRanks[j])*(bRanks[i]-bRanks[j]) > 0 {
				concordant++
			} else if (aRanks[i]-aRanks[j])*(bRanks[i]-bRanks[j]) < 0 {
				discordant++
			} // else neither concordant nor discordant
		}
	}
	return float64(concordant-discordant) / float64(size*(size-1)/2.0)
}

// Mean returns the mean value of a float array
func Mean(a []float64) float64 {
	if len(a) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range a {
		sum += v
	}
	return sum / float64(len(a))
}

// Rank returns an array containing the ranks of a slice a
func Rank(a []float64) []int {
	rank := make([]int, len(a))
	for i := range a {
		for j := range a {
			if a[j] < a[i] {
				rank[i]++
			}
		}
	}
	return rank
}

//StdDev returns the standard deviation of a float slice
func StdDev(a []float64) float64 {
	if len(a) == 0 {
		return math.NaN()
	}
	mean, sum := Mean(a), 0.0
	for _, v := range a {
		sum += (v - mean) * (v - mean)
	}
	return math.Sqrt(sum / float64(len(a)))
}
