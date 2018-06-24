package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// OnlineIndexer is used to execute online indexing. The user can supply a map containing
// distances from original datasets and the indexer returns the coordinates of the
// specified dataset.
type OnlineIndexer struct {
	coordinates []DatasetCoordinates       // the coordinates of the datasets
	estimator   DatasetSimilarityEstimator // estimator object to calculate distances
	script      string                     // script to evaluate the acceptable solutions

	dimensionality    int // the number of dimensions of the coordinates
	datasetsToCompare int // number of datasets to compare when calculate is called
}

// NewOnlineIndexer is a constructor function used to initialize an OnlineIndexer
// object.
func NewOnlineIndexer(estimator DatasetSimilarityEstimator,
	coordinates []DatasetCoordinates,
	script string) *OnlineIndexer {
	indexer := new(OnlineIndexer)
	indexer.coordinates = coordinates
	indexer.estimator = estimator
	indexer.script = script

	indexer.dimensionality = len(coordinates[0])
	indexer.datasetsToCompare = len(coordinates)
	return indexer
}

// DatasetsToCompare is a setter method to determine the number of datasets
// that will be utilized for the assigment of coordinates
func (o *OnlineIndexer) DatasetsToCompare(datasets int) {
	if datasets > 0 {
		o.datasetsToCompare = datasets
	}
}

// Calculate method is responsible to calculate the coordinates of the specified
// dataset. In case that such a dataset cannot be represented by the specified
// coordinates system, an error is returned.
func (o *OnlineIndexer) Calculate(dataset *Dataset) (DatasetCoordinates, float64, error) {
	// calculate the distances for the new dataset
	log.Println("Creating dataset permutation")
	perm := rand.Perm(len(o.estimator.Datasets()))

	writer, err := ioutil.TempFile("/tmp", "coordinates.csv-")
	if err != nil {
		log.Println(err)
		return nil, -1.0, err
	}
	defer os.Remove(writer.Name())

	for i := 0; i < o.dimensionality; i++ {
		fmt.Fprintf(writer, "x%d,", i+1)
	}
	fmt.Fprintf(writer, "d\n")

	actualDistances := make([]float64, o.datasetsToCompare)
	for c, i := range perm {
		if c >= o.datasetsToCompare {
			break
		}
		dat := o.estimator.Datasets()[i]
		sim := o.estimator.Similarity(dataset, dat)

		for j := 0; j < o.dimensionality; j++ {
			fmt.Fprintf(writer, "%.5f,", o.coordinates[i][j])
		}
		actualDistances[c] = SimilarityToDistance(sim)
		fmt.Fprintf(writer, "%.5f\n", actualDistances[c])
	}
	writer.Close()
	coords, err := o.executeScript(writer.Name())
	if err != nil {
		return nil, -1.0, err
	}
	calculatedDistances := make([]float64, o.datasetsToCompare)
	for c, i := range perm {
		if c >= o.datasetsToCompare {
			break
		}
		for j := 0; j < o.dimensionality; j++ {
			calculatedDistances[c] += (coords[j] - o.coordinates[i][j]) * (coords[j] - o.coordinates[i][j])
		}
		calculatedDistances[c] = math.Sqrt(calculatedDistances[c])
	}

	return coords, o.stress(actualDistances, calculatedDistances), err
}

// returns the stress factor (the difference between the actual and the calculated
// distances
func (o *OnlineIndexer) stress(actual, calculated []float64) float64 {
	sum1 := 0.0
	for i := range actual {
		sum1 += (actual[i] - calculated[i]) * (actual[i] - calculated[i])
	}
	return math.Sqrt(sum1)
}

// solveQuadSystem solves the quadratic polynomial system in order to identify
// the coordinates of the new point.
func (o *OnlineIndexer) executeScript(fileName string) (DatasetCoordinates, error) {
	log.Println("Executing", o.script, "with argument", fileName)
	// execute script
	cmd := exec.Command(o.script, fileName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
		return nil, err
	}
	buffer := bytes.NewBuffer(output)
	l, err := buffer.ReadString('\n')
	if err != nil {
		return nil, err
	}
	ar := strings.Split(l, " ")
	solution := make(DatasetCoordinates, o.dimensionality)
	for j := range solution {
		solution[j], err = strconv.ParseFloat(strings.TrimSpace(ar[j]), 64)
		if err != nil {
			return nil, err
		}
	}
	return solution, nil
}
