package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/giagiannis/data-profiler/core"
)

type printUtilsParams struct {
	similarities *core.DatasetSimilarityMatrix
	scores       *core.DatasetScores
}

func printUtilsParseParams() *printUtilsParams {
	params := new(printUtilsParams)

	scoresFile :=
		flag.String("s", "", "scores file")
	similaritiesPath :=
		flag.String("sim", "", "the path of the similarity matrix file")

	flag.Parse()
	if *scoresFile == "" && *similaritiesPath == "" {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *scoresFile != "" {
		f, err := os.Open(*scoresFile)
		if err != nil {
			log.Fatalln(err)
		}
		buf, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatalln(err)
		}
		f.Close()
		params.scores = core.NewDatasetScores()
		params.scores.Deserialize(buf)
	}

	if *similaritiesPath != "" {
		params.similarities = core.NewDatasetSimilarities(0)
		log.Println("Reading", *similaritiesPath)
		f, err := os.Open(*similaritiesPath)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatalln(err)
		}
		err = params.similarities.Deserialize(b)

		if err != nil {
			log.Fatalln(err)
		}
	}
	return params

}

func printUtilsRun() {
	params := printUtilsParseParams()
	if params.scores != nil {
		for k, v := range params.scores.Scores {
			fmt.Printf("%s: %.5f\n", k, v)
		}
	}

	if params.similarities != nil {
		for i := 0; i < params.similarities.Capacity(); i++ {
			for j := 0; j < params.similarities.Capacity(); j++ {
				fmt.Printf("%.5f\t", params.similarities.Get(i, j))
			}
			fmt.Println()
		}
	}
}
