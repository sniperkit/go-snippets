package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	trainSet             = "../_testdata/datatraining.csv"
	testSet              = "../_testdata/datatest.csv"
	analysisScript       = "../_rscripts/pca.R"
	pairScript           = "../_rscripts/random_number.R"
	mlScript             = "../_rscripts/lm.R"
	mlScriptAppx         = "../_rscripts/svm-appx.R"
	mdsScript            = "../_rscripts/mdscaling.R"
	tmpDir               = "/tmp/"
	onlineIndexingScript = "../_rscripts/sa.R"
	operatorScript       = "../_testdata/sum.py"
)

func init() {
	log.SetOutput(ioutil.Discard)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	rand.Seed(int64(time.Now().Nanosecond()))
}

// Creates a random tuple with floats, containing fields fields.
func randomTupleGenerator(fields int) []float64 {
	res := make([]float64, fields)
	for i := 0; i < fields; i++ {
		res[i] = rand.Float64()
	}
	return res
}

// Creates a pool of tuples (its sizeis determined by poolSize) and a number of
// datasets each containing tuples with the specified number of attributes.
func createPoolBasedDatasets(poolSize, datasets, attributes int) []*Dataset {
	pool := *new([][]float64)
	for i := 0; i < poolSize; i++ {
		pool = append(pool, randomTupleGenerator(attributes))
	}
	fileNames := make([]string, datasets)
	for i := 0; i < datasets; i++ {
		dSize := rand.Int()%(poolSize/2) + poolSize/2
		buffer := make([]byte, 4)
		rand.Read(buffer)
		fileNames[i] = fmt.Sprintf("%s%x.txt", tmpDir, buffer)
		f, _ := os.Create(fileNames[i])
		defer f.Close()

		builder := bytes.Buffer{}
		for i := 0; i < attributes-1; i++ {
			builder.WriteString(fmt.Sprintf("x%d, ", i))
		}
		builder.WriteString("class\n")
		tuplesChosen := make(map[int]bool)
		for j := 0; j < dSize; j++ {
			var idx int
			for idx = rand.Int() % len(pool); tuplesChosen[idx]; idx = rand.Int() % len(pool) {
			}
			tuplesChosen[idx] = true
			tuple := pool[idx]
			tupleLength := len(tuple)
			for i, v := range tuple {
				builder.WriteString(fmt.Sprintf("%.5f", v))
				if i < tupleLength-1 {
					builder.WriteString(",")
				}
			}
			builder.WriteString("\n")
		}
		f.Write(builder.Bytes())
	}
	result := make([]*Dataset, len(fileNames))
	for i, f := range fileNames {
		result[i] = NewDataset(f)
	}

	return result
}

// Creates a pool of tuples (its size is determined by poolSize) and a number of
// datasets each containing tuples with the specified number of attributes.
func createPoolBasedDatasetsStrict(poolSize, datasetSize, datasets, attributes int) []*Dataset {
	pool := *new([][]float64)
	for i := 0; i < poolSize; i++ {
		pool = append(pool, randomTupleGenerator(attributes))
	}
	fileNames := make([]string, datasets)
	for i := 0; i < datasets; i++ {
		buffer := make([]byte, 4)
		rand.Read(buffer)
		fileNames[i] = fmt.Sprintf("%s%x.txt", tmpDir, buffer)
		f, _ := os.Create(fileNames[i])
		defer f.Close()

		builder := bytes.Buffer{}
		for i := 0; i < attributes-1; i++ {
			builder.WriteString(fmt.Sprintf("x%d, ", i))
		}
		builder.WriteString("class\n")
		tuplesChosen := make(map[int]bool)
		for j := 0; j < datasetSize; j++ {
			var idx int
			for idx = rand.Int() % len(pool); tuplesChosen[idx]; idx = rand.Int() % len(pool) {
			}
			tuplesChosen[idx] = true
			tuple := pool[idx]
			tupleLength := len(tuple)
			for i, v := range tuple {
				builder.WriteString(fmt.Sprintf("%.5f", v))
				if i < tupleLength-1 {
					builder.WriteString(",")
				}
			}
			builder.WriteString("\n")
		}
		f.Write(builder.Bytes())
	}
	result := make([]*Dataset, len(fileNames))
	for i, f := range fileNames {
		result[i] = NewDataset(f)
	}

	return result
}

func cleanDatasets(datasets []*Dataset) {
	for _, f := range datasets {
		os.Remove(f.Path())
	}
}

func createLinearDatasets(datasets, datasetSize, attributes int, noise float64) []*Dataset {
	var result []*Dataset
	for i := 0; i < datasets; i++ {
		coeff := make([]float64, attributes-1)
		for i := range coeff {
			coeff[i] = rand.Float64()
		}
		builder := new(bytes.Buffer)
		for i := range coeff {
			builder.WriteString(fmt.Sprintf("x%d, ", i))
		}
		builder.WriteString("y\n")
		//		for i := 0; i < datasetSize; i++ {
		for step := 0.0; step < 1.0; step += 1.0 / float64(datasetSize) {
			sum := 0.0
			for _, c := range coeff {
				sum += c*step + rand.Float64()*noise
				builder.WriteString(fmt.Sprintf("%.5f, ", step))
			}
			builder.WriteString(fmt.Sprintf("%.5f\n", sum))
		}
		f, _ := ioutil.TempFile("/tmp", "lineardataset")
		f.Write(builder.Bytes())
		f.Close()
		result = append(result, NewDataset(f.Name()))
	}
	return result
}
