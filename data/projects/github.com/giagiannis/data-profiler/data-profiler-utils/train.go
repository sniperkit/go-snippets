package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/giagiannis/data-profiler/core"
)

type trainParams struct {
	testset     *core.Dataset   // path of the test set
	script      *string         // path of the ML script
	datasets    []*core.Dataset // path of the input datasets
	concurrency *int            // number of threads to be utilized
	output      *string         // path of the outfile
	logfile     *string         // logfile path
	repetition  *int            // times to execute experiments
}

func trainParseParams() *trainParams {
	params := new(trainParams)
	testset :=
		flag.String("t", "", "path of the test set")
	params.script =
		flag.String("s", "", "path of the ML script")
	input :=
		flag.String("i", "", "path of the datasets dir")
	params.output =
		flag.String("o", "", "path of the output file")
	params.logfile =
		flag.String("l", "", "path of the log file")
	params.concurrency =
		flag.Int("p", 1, "number of threads")
	params.repetition =
		flag.Int("r", 1, "times to repeat experiments (prints median if >1)")
	flag.Parse()
	setLogger(*params.logfile)
	if *testset == "" ||
		*params.script == "" ||
		*params.output == "" ||
		*input == "" {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	params.datasets = core.DiscoverDatasets(*input)
	params.testset = core.NewDataset(*testset)
	return params
}

func trainRun() {
	params := trainParseParams()
	confparams := make(map[string]string)
	confparams["script"] = *params.script
	confparams["testset"] = params.testset.Path()
	eval, err := core.NewDatasetEvaluator(core.OnlineEval, confparams)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	results := core.NewDatasetScores()
	done := make(chan bool)
	c := make(chan bool, *params.concurrency)
	for i := 0; i < *params.concurrency; i++ {
		c <- true
	}
	lock := new(sync.Mutex)
	for _, d := range params.datasets {
		go func(d core.Dataset, c, done chan bool) {
			<-c
			log.Println("Training with ", d.Path())
			res := make([]float64, 0)
			for i := 0; i < *params.repetition; i++ {
				val, err := eval.Evaluate(d.Path())
				if err != nil {
					fmt.Fprintf(os.Stderr, err.Error())
					os.Exit(1)
				}
				res = append(res, val)
			}
			sort.Float64s(res)
			lock.Lock()
			results.Scores[d.Path()] = res[len(res)/2]
			lock.Unlock()
			c <- true
			done <- true
		}(*d, c, done)
	}
	for range params.datasets {
		<-done
	}
	log.Println("Serializing output to ", *params.output)
	outfile, er := os.OpenFile(*params.output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer outfile.Close()

	if er != nil {
		fmt.Fprintln(os.Stderr, er)
		os.Exit(1)
	}
	b, err := results.Serialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	outfile.Write(b)
}
