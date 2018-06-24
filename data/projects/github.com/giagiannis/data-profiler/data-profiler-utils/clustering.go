package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/giagiannis/data-profiler/core"
)

type clusteringParams struct {
	input        *string                       //datasets directory to be discovered
	similarities *core.DatasetSimilarityMatrix // similarities
	scores       *core.DatasetScores           // scores file
	logfile      *string                       // log file
	concurrency  *int                          // max number of threads
	output       *string                       // output file
}

func clusteringParseParams() *clusteringParams {
	params := new(clusteringParams)

	params.input =
		flag.String("i", "", "input datasets path")
	similaritiesFile :=
		flag.String("si", "", "similarities file - required")
	scoresFile :=
		flag.String("sc", "", "scores file - required")
	params.logfile =
		flag.String("l", "", "logfile (default: stdout)")
	params.concurrency =
		flag.Int("t", 1, "number of threads")
	params.output =
		flag.String("o", "", "output file (default: stdout)")
	flag.Parse()
	setLogger(*params.logfile)

	if *params.input == "" || *similaritiesFile == "" || *scoresFile == "" {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// parse similarities file
	f, err := os.Open(*similaritiesFile)
	defer f.Close()
	if err != nil {
		log.Fatalln(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	params.similarities = core.NewDatasetSimilarities(0)
	err = params.similarities.Deserialize(buf)
	if err != nil {
		log.Fatalln(err)
	}

	// parse scores file
	fs, err := os.Open(*scoresFile)
	if err != nil {
		log.Fatalln(err)
	}
	buf, err = ioutil.ReadAll(fs)
	if err != nil {
		log.Fatalln(err)
	}
	params.scores = core.NewDatasetScores()
	err = params.scores.Deserialize(buf)
	if err != nil {
		log.Fatalln(err)
	}

	return params
}

func clusteringRun() {
	params := clusteringParseParams()
	log.Println("Initializing clustering object")
	datasets := core.DiscoverDatasets(*params.input)
	cls := core.NewClustering(params.similarities, datasets)
	cls.SetConcurrency(*params.concurrency)
	log.Println("Executing computation")
	cls.Compute()
	log.Println("Done")
	outF := os.Stdout
	if *params.output != "" {
		var err error
		outF, err = os.OpenFile(*params.output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		defer outF.Close()
	}
	maxTreeHeight, _ := cls.Results().Heights()
	fmt.Fprintf(outF, "level \"cluster\" \"avg error\" \"max error\" \"median error\"\n")
	maxError, _, _ := evaluateClusters(cls.Results().GetClusters(0), params.scores)
	for i := 0; i <= maxTreeHeight; i++ {
		diff, maxDiff, medianDiff := evaluateClusters(cls.Results().GetClusters(i), params.scores)
		fmt.Fprintf(outF, "%d %d %.5f %.5f %.5f\n", i,
			len(cls.Results().GetClusters(i)), diff/maxError, maxDiff/maxError, medianDiff/maxError)
	}
	fmt.Fprintf(outF, "\n\n") // to be used for different gnuplot indices
	if outF != os.Stdout {
		log.Println("Results written to", *params.output)
	}
}

func evaluateClusters(clusters [][]*core.Dataset, scores *core.DatasetScores) (float64, float64, float64) {
	sumDiffs := 0.0
	differences := make([]float64, 0)
	for _, c := range clusters {
		_, d := evaluateCluster(c, scores)
		sumDiffs += d
		differences = append(differences, d)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(differences)))
	return sumDiffs / float64(len(clusters)), differences[0], differences[len(differences)/2]
}

// returns various metrics for the specified clusters
func evaluateCluster(datasets []*core.Dataset, scores *core.DatasetScores) (float64, float64) {
	sum := 0.0
	maxError, minError := -1.0, -1.0
	for _, d := range datasets {
		s := scores.Scores[d.Path()]
		sum += s
		if maxError == -1 || s > maxError {
			maxError = s
		}
		if minError == -1 || s < minError {
			minError = s
		}
	}
	return sum / float64(len(datasets)), maxError - minError
}
