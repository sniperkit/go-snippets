package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/giagiannis/data-profiler/core"
)

type similaritiesParams struct {
	input            *string                                 // datasets directory to be discovered
	output           *string                                 // output file
	simType          *core.DatasetSimilarityEstimatorType    // similarity type
	logfile          *string                                 // logfile
	options          *string                                 // options for the estimators
	populationPolicy *core.DatasetSimilarityPopulationPolicy // defines the population policy
	estimatorPath    *string                                 // place to store estimator object
}

func similaritiesParseParams() *similaritiesParams {
	params := new(similaritiesParams)
	params.input =
		flag.String("i", "", "input datasets path")
	params.output =
		flag.String("o", "", "where to store similarities file")
	params.logfile =
		flag.String("l", "", "logfile (default: stderr)")
	estType :=
		flag.String("t", "BHATTACHARYYA", "similarity type [JACCARD|BHATTACHARYYA|SCRIPT|ORDER]")
	params.options =
		flag.String("opt", "", "options in the form val1=key1,val2=key2 (list for opts list)")
	popPolicy :=
		flag.String("p", "FULL", "population policy [FULL|APRX] along with options in the form POLICY,val1=key1,val2=key2")
	params.estimatorPath =
		flag.String("e", "", "if set, serializes the estimator to the specified path")
	flag.Parse()
	setLogger(*params.logfile)

	// population policy parsing
	popPolicyType := strings.Split(*popPolicy, ",")[0]
	params.populationPolicy = new(core.DatasetSimilarityPopulationPolicy)
	params.populationPolicy.Parameters = make(map[string]float64)
	if popPolicyType == "FULL" {
		params.populationPolicy.PolicyType = core.PopulationPolicyFull
	} else if popPolicyType == "APRX" {
		params.populationPolicy.PolicyType = core.PopulationPolicyAprx
		idx := strings.Index(*popPolicy, ",")
		if idx > -1 {
			for k, v := range parseOptions((*popPolicy)[idx+1 : len(*popPolicy)]) {
				val, _ := strconv.ParseFloat(v, 64)
				params.populationPolicy.Parameters[k] = val
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Population policy unknown\n")
		os.Exit(1)
	}

	params.simType = new(core.DatasetSimilarityEstimatorType)
	simEstType := core.NewDatasetSimilarityEstimatorType(*estType)
	if simEstType != nil {
		*params.simType = *simEstType
	}

	if *params.options == "list" {
		for i, v := range core.DatasetSimilarityEstimatorAvailableTypes {
			fmt.Println(i+1, v)
			a := core.NewDatasetSimilarityEstimator(v, nil)
			for k, v := range a.Options() {
				fmt.Println("\t", k, ":", v)
			}
		}
		os.Exit(0)
	}

	if *params.input == "" ||
		*params.output == "" ||
		params.simType == nil {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	return params
}

func similaritiesRun() {
	params := similaritiesParseParams()
	datasets := core.DiscoverDatasets(*params.input)
	est := core.NewDatasetSimilarityEstimator(*params.simType, datasets)
	est.Configure(parseOptions(*params.options))
	est.SetPopulationPolicy(*params.populationPolicy)
	est.Compute()
	log.Printf("Similarity Matrix computation took %.5f sec\n", est.Duration())

	outfile, er := os.OpenFile(*params.output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if er != nil {
		log.Fatal(er)
	}
	defer outfile.Close()
	// serializing similarity matrix
	outfile.Write(est.SimilarityMatrix().Serialize())

	idxFile, er := os.OpenFile(*params.output+".idx", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer idxFile.Close()
	for i, d := range datasets {
		fmt.Fprintf(idxFile, "%d\t%s\n", i, d.Path())
	}
	if er != nil {
		log.Fatal(er)
	}
	// serializing estimator
	if *params.estimatorPath != "" {
		outfile, er = os.OpenFile(*params.estimatorPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if er != nil {
			fmt.Fprintln(os.Stderr, er)
			os.Exit(1)
		}
		defer outfile.Close()
		outfile.Write(est.Serialize())
		log.Println("Serialized estimator to file", outfile.Name())
	}
}
