package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/giagiannis/data-profiler/core"
)

type indexingParams struct {
	estimator   core.DatasetSimilarityEstimator // estimator object
	coordinates []core.DatasetCoordinates       // coordinates of datasets
	datasets    []*core.Dataset                 // new datasets to index
	script      *string                         // indexing script needed by the indexer
	output      *string                         // where to store the output
	logfile     *string                         // logfile of execution
	k           *string                         // determines the number of datasets to compare
	concurrency *int                            //number of threads to be utilized for the exec
	repetition  *int                            //number of threads to be utilized for the exec
}

func indexingParseParams() *indexingParams {
	params := new(indexingParams)
	estimatorFile :=
		flag.String("e", "", "the serialized estimator file")
	coordinatesFile :=
		flag.String("c", "", "the dataset coordinates file")
	params.script =
		flag.String("s", "", "the script file to be used by the indexer")
	params.output =
		flag.String("o", "", "the output file")
	params.logfile =
		flag.String("l", "", "the log file")
	datasetsPath :=
		flag.String("i", "", "the new dataset to be indexed - if a dir, only the first dataset is considered")
	params.k =
		flag.String("k", "0", "comma separated list of datasets to be used for comparison - 0 means all")
	params.concurrency =
		flag.Int("t", 1, "the number of threads to spawn for the indexing")
	params.repetition =
		flag.Int("r", 1, "times to repeat the experiments for random experiments")
	flag.Parse()
	setLogger(*params.logfile)

	if *estimatorFile == "" || *coordinatesFile == "" || *params.script == "" ||
		*params.output == "" || *datasetsPath == "" {
		fmt.Fprintln(os.Stderr, "Missing arguments, usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// parse estimator object
	f, err := os.Open(*estimatorFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	params.estimator = core.DeserializeSimilarityEstimator(buffer)

	// parse coordinates file
	f, err = os.Open(*coordinatesFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	coordBuffer, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	params.coordinates = make([]core.DatasetCoordinates, 0)
	for _, l := range strings.Split(string(coordBuffer), "\n") {
		e := false
		tuple := make(core.DatasetCoordinates, 0)
		for _, v := range strings.Split(l, " ") {
			if !e {
				v = strings.TrimSpace(v)
				if v != "" {
					val, err := strconv.ParseFloat(v, 64)
					e = (err != nil)
					tuple = append(tuple, val)
				}
			}
		}
		if !e && len(tuple) > 0 {
			params.coordinates = append(params.coordinates, tuple)
		}
	}

	// parse datasets
	params.datasets = core.DiscoverDatasets(*datasetsPath)

	return params
}

func indexingRun() {
	params := indexingParseParams()
	c := make(chan bool, *params.concurrency)
	done := make(chan resultsStruct)
	for j := 0; j < *params.concurrency; j++ {
		c <- true
	}
	testCases := 0
	for _, k := range strings.Split(*params.k, ",") {
		kint, err := strconv.Atoi(k)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		execution := func(c chan bool, done chan resultsStruct, k int) {
			indexer := core.NewOnlineIndexer(params.estimator, params.coordinates, *params.script)
			indexer.DatasetsToCompare(k)
			<-c
			start := time.Now()
			coo, gof, err := indexer.Calculate(params.datasets[0])
			duration := time.Since(start)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			c <- true
			done <- resultsStruct{k, coo, gof, duration}
		}
		if kint == 0 || kint == len(params.estimator.Datasets()) {
			kint = len(params.estimator.Datasets())
			go execution(c, done, kint)
			testCases++
		} else {
			for i := 0; i < *params.repetition; i++ {
				go execution(c, done, kint)
				testCases++
			}
		}
	}
	outfile, er := os.OpenFile(*params.output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if er != nil {
		fmt.Fprintln(os.Stderr, er)
		os.Exit(1)
	}
	defer outfile.Close()

	results := make(map[int][]resultsStruct)
	for i := 0; i < testCases; i++ {
		res := <-done
		results[res.id] = append(results[res.id], res)
	}
	// find best
	keys := make([]int, 0)
	for k := range results {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	best := results[keys[len(keys)-1]][0]

	fmt.Fprintf(outfile, "k distance stress duration\n")
	for _, k := range keys {
		d, g, du := aggregateResults(results[k], best)
		fmt.Fprintf(outfile, "%d %.5f %.5f %.5f\n", k, d, g, du)
	}
}

type resultsStruct struct {
	id       int
	coords   core.DatasetCoordinates
	stress   float64
	duration time.Duration
}

// receives the structs and returns the distance between the best, the stress and the duration
func aggregateResults(array []resultsStruct, best resultsStruct) (float64, float64, float64) {
	coordsDistance := func(a, b core.DatasetCoordinates) float64 {
		sum := 0.0
		for i := range a {
			sum += (a[i] - b[i]) * (a[i] - b[i])
		}
		return math.Sqrt(sum)
	}
	dst, stress, dur := 0.0, 0.0, 0.0
	for i := range array {
		dst += coordsDistance(array[i].coords, best.coords)
		stress += array[i].stress
		dur += array[i].duration.Seconds()
	}
	l := float64(len(array))
	return dst / l, stress / l, dur / l
}
