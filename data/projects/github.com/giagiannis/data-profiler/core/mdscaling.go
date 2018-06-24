package core

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// DatasetCoordinates is a struct for representing the dataset coordinates
type DatasetCoordinates []float64

// DeserializeCoordinates instantiated a new DatasetCoordinates slice, based
// on a CSV serialization form
func DeserializeCoordinates(buffer []byte) []DatasetCoordinates {
	coords := make([]DatasetCoordinates, 0)
	for _, line := range strings.Split(string(buffer), "\n") {
		a := strings.Split(line, ",")
		res := make(DatasetCoordinates, 0)
		if len(a) > 0 {
			for _, s := range a {
				if s != "" {
					v, err := strconv.ParseFloat(s, 64)
					if err != nil {
						log.Fatalln(err)
					}
					res = append(res, v)
				}
			}
			if len(res) > 0 {
				coords = append(coords, res)
			}
		}
	}
	return coords
}

// SerializeCoordinates returns a CSV serilization of a coordinates slice
func SerializeCoordinates(coords []DatasetCoordinates) []byte {
	buffer := new(bytes.Buffer)
	for _, d := range coords {
		for i, c := range d {
			buffer.WriteString(fmt.Sprintf("%.5f", c))
			if i != len(coords[0]) {
				buffer.WriteString(",")
			}
		}
		buffer.WriteString(fmt.Sprintf("\n"))
	}
	return buffer.Bytes()
}

// MDScaling is responsible for the execution of a MultiDimensional Scaling
// algorithm in order to provide coefficients for each dataset, based on a
// a similarity matrix.
type MDScaling struct {
	script string                   // script to be used for the execution
	k      int                      // number of output dimensions
	matrix *DatasetSimilarityMatrix // the similarity matrix

	coordinates []DatasetCoordinates // the coordinates matrix
	gof         float64              // the gof factor
	stress      float64              // the stress factor
}

// NewMDScaling is the default MDScaling constructor; it initializes a new
// MDScaling object, based on the provided DatasetSimilarities struct and the
// k factor that determines the number of target dimensions. If k<1, then
// auto estimation takes place.
func NewMDScaling(matrix *DatasetSimilarityMatrix, k int, script string) *MDScaling {
	mds := new(MDScaling)
	mds.matrix = matrix
	mds.k = k
	mds.script = script
	return mds
}

// Compute functions executes the Multidimensional Scaling computation.
func (md *MDScaling) Compute() error {
	// create a temp file containing similarity matrix as a csv
	writer, err := ioutil.TempFile("/tmp", "similarities")
	if err != nil {
		return err
	}
	log.Println("Writing distances to file", writer.Name())
	for i := 0; i < md.matrix.Capacity(); i++ {
		for j := 0; j < md.matrix.Capacity(); j++ {
			dist := SimilarityToDistance(md.matrix.Get(i, j))
			writer.WriteString(fmt.Sprintf("%.5f", dist))
			if j < md.matrix.Capacity()-1 {
				writer.WriteString(",")
			}
		}
		writer.WriteString("\n")
	}
	writer.Close()
	// execute computation
	if md.k < 1 || md.k > md.matrix.Capacity()-1 { // binary search in the interval [1, n-1]
		return errors.New("K factor must be between [1, n-1], n being the # of datasets")
	}
	// execute solution
	md.coordinates, md.gof, md.stress, err = md.executeScript(writer.Name())
	if err != nil {
		return err
	}
	os.Remove(writer.Name())
	return nil
}

// Coordinates getter returns the dataset coordinates in a nxk slice (n being
// the number of datasets).
func (md *MDScaling) Coordinates() []DatasetCoordinates {
	return md.coordinates
}

// Gof getter returns the gof factor for the calculated solution
func (md *MDScaling) Gof() float64 {
	return md.gof
}

// Stress returns the stress after the execution of the Sammon mapping
func (md *MDScaling) Stress() float64 {
	return md.stress
}

// Variances function returns the variances of the principal coordinates
func (md *MDScaling) Variances() ([]float64, error) {
	if md.coordinates == nil {
		return nil, errors.New("Coordinates not calculated")
	}
	variances := make([]float64, md.k)
	for i := 0; i < md.k; i++ {
		max, min := md.coordinates[0][i], md.coordinates[0][i]
		for j := 0; j < len(md.coordinates); j++ {
			if md.coordinates[j][i] > max {
				max = md.coordinates[j][i]
			}
			if md.coordinates[j][i] < min {
				min = md.coordinates[j][i]
			}
		}
		variances[i] = math.Abs(max - min)
	}
	return variances, nil
}

// executeScript is used to execute the specified MDS script and parse
// its results. Not part of the structs public API. If successful, returns
// the coordinates slice, the gof factor and nil errors; If not successful,
// returns nil results and an error object
//
func (md *MDScaling) executeScript(smPath string) ([]DatasetCoordinates, float64, float64, error) {
	command := exec.Command(md.script, smPath, strconv.Itoa(md.k))
	buf, err := command.CombinedOutput()
	if err != nil {
		log.Println(string(buf))
		return nil, math.NaN(), math.NaN(), err
	}
	lines := strings.Split(string(buf), "\n")

	// parse gof float
	gof, err := strconv.ParseFloat(lines[0], 64)
	if err != nil {
		return nil, math.NaN(), math.NaN(), err
	}
	// parse stress float
	stress, err := strconv.ParseFloat(lines[1], 64)
	if err != nil {
		return nil, math.NaN(), math.NaN(), err
	}
	count := 0
	for i := 2; i < len(lines) && len(lines[i]) > 0; i++ {
		if len(lines[i]) > 0 {
			count++
		}
	}
	// parse coordinates slices
	coordinates := make([]DatasetCoordinates, count)
	for i := 2; i < len(lines) && len(lines[i]) > 0; i++ {
		splitLine := strings.Split(lines[i], " ")
		coordinates[i-2] = make([]float64, md.k)
		for j := 0; j < md.k; j++ {
			coordinates[i-2][j], err = strconv.ParseFloat(splitLine[j], 64)
			if err != nil {
				log.Println(err)
				return nil, math.NaN(), math.NaN(), err
			}
		}
	}
	return coordinates, gof, stress, nil
}

// DistanceToSimilarity returns the similarity based on the distance
func DistanceToSimilarity(distance float64) float64 {
	//	return 1.0 / (1.0 + distance)
	return math.Sqrt(1 - distance)
}

// SimilarityToDistance returns the distance based on the similarity
func SimilarityToDistance(similarity float64) float64 {
	//	return 1.0/similarity - 1.0
	return 1 - similarity*similarity
}
