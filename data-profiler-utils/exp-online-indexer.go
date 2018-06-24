package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/giagiannis/data-profiler/core"
)

type expOnlineIndexerParams struct {
	indexer  *core.OnlineIndexer
	datasets []*core.Dataset
	est      core.DatasetSimilarityEstimator
	coords   []core.DatasetCoordinates

	outputStress *string
	outputCoords *string
	compDatasets *int
}

func expOnlineIndexerParseParams() *expOnlineIndexerParams {
	params := new(expOnlineIndexerParams)
	coordinatesPath :=
		flag.String("c", "", "file that keeps the original coordinates")
	estimatorPath :=
		flag.String("e", "", "path of the serialized estimator object")
	saScript :=
		flag.String("s", "", "Script that executes Simulated Annealing")
	logger :=
		flag.String("l", "", "log file")
	params.outputStress =
		flag.String("os", "", "output for the stress values")
	params.outputCoords =
		flag.String("oc", "", "output for the coords file")
	inputDatasets :=
		flag.String("i", "", "datasets to index")
	params.compDatasets =
		flag.Int("com", 0, "datasets to be used for the comparison")

	flag.Parse()
	setLogger(*logger)
	if *coordinatesPath == "" || *estimatorPath == "" || *saScript == "" {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// parse estimator
	f, err := os.Open(*estimatorPath)
	if err != nil {
		log.Fatalln(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	params.est = core.DeserializeSimilarityEstimator(buf)
	//params.sm = est.SimilarityMatrix()

	// parse coordinates files
	buf, err = ioutil.ReadFile(*coordinatesPath)
	if err != nil {
		log.Fatalln(err)
	}
	//lines := strings.Split(strings.TrimSpace(string(buf)), "\n")
	//if len(lines) < 1 {
	//	log.Fatalln("Empty coordinates file")
	//}
	//var datasetCoordinates []core.DatasetCoordinates

	//header := strings.Split(strings.TrimSpace(lines[0]), " ")
	//dimensionality := len(header)
	//for i := 1; i < len(lines); i++ {
	//	coords := make([]float64, dimensionality)
	//	lineArray := strings.Split(lines[i], " ")
	//	for j := range coords {
	//		val, err := strconv.ParseFloat(lineArray[j], 64)
	//		if err != nil {
	//			log.Println(err)
	//		} else {
	//			coords[j] = val
	//		}
	//	}
	//	datasetCoordinates = append(datasetCoordinates, coords)
	//}
	params.coords = core.DeserializeCoordinates(buf)

	// discovering datasets
	if *inputDatasets != "" {
		params.datasets = core.DiscoverDatasets(*inputDatasets)
	}

	// construct indexer object
	if len(params.datasets) > 0 {
		params.indexer = core.NewOnlineIndexer(params.est, params.coords, *saScript)
		params.indexer.DatasetsToCompare(*params.compDatasets)
	}
	return params
}

func expOnlineIndexerRun() {
	params := expOnlineIndexerParseParams()
	stressOutput := setOutput(*params.outputStress)
	defer stressOutput.Close()
	sm := params.est.SimilarityMatrix()
	currentStress := sammonStress(sm, params.coords)
	log.Println("Initial total stress:", currentStress)
	fmt.Fprintln(stressOutput, currentStress)

	if len(params.datasets) == 0 {
		return
	}

	targetDatasets := make([]*core.Dataset, len(params.est.Datasets()))
	for i := range targetDatasets {
		targetDatasets[i] = params.est.Datasets()[i]
	}

	for _, d := range params.datasets {
		old := sm
		sm = core.NewDatasetSimilarities(old.Capacity() + 1)
		for i := 0; i < old.Capacity(); i++ {
			for j := 0; j < old.Capacity(); j++ {
				sm.Set(i, j, old.Get(i, j))
				sm.Set(j, i, old.Get(j, i))
			}
		}
		for i := 0; i < sm.Capacity()-1; i++ {
			sim := params.est.Similarity(d, targetDatasets[i])
			sm.Set(i, sm.Capacity()-1, sim)
			sm.Set(sm.Capacity()-1, i, sim)
		}
		targetDatasets = append(targetDatasets, d)

		coord, stress, err := params.indexer.Calculate(d)
		if err == nil {
			params.coords = append(params.coords, coord)
			log.Println("Results:", coord, stress)

			currentStress = sammonStress(sm, params.coords)
			log.Println("Current total stress:", currentStress)
			fmt.Fprintln(stressOutput, currentStress)
		} else {
			log.Println(err)
		}
	}
	f := setOutput(*params.outputCoords)
	defer f.Close()
	for i := 0; i < len(params.coords[0]); i++ {
		fmt.Fprintf(f, "x_%d ", i+1)
	}
	fmt.Fprintf(f, "\n")

	for i := range params.coords {
		for j := range params.coords[i] {
			fmt.Fprintf(f, "%.5f ", params.coords[i][j])
		}
		fmt.Fprintf(f, "\n")
	}
}

func sammonStress(similarityMatrix *core.DatasetSimilarityMatrix, coords []core.DatasetCoordinates) float64 {
	stress := 0.0
	foo := 0.0
	for i := 0; i < len(coords); i++ {
		for j := i + 1; j < len(coords); j++ {
			actual := core.SimilarityToDistance(similarityMatrix.Get(i, j))
			measured := getDistance(coords[i], coords[j])
			if actual != 0 {
				stress += ((actual - measured) * (actual - measured) / actual)
				foo += actual
			} else {
				log.Println(i, j, "zero")
			}
		}
	}
	return stress / foo
}
