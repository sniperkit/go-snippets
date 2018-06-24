package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/giagiannis/data-profiler/core"
)

type mdsParams struct {
	logfile *string         // the logfile to be used
	output  *string         // the output file
	modules map[string]bool // which modules to activate for the util

	script       *string                       // the mds script to run
	k            *int                          // the number of the coordinates to eval
	similarities *core.DatasetSimilarityMatrix // the similarity matrix
}

var mdsModules = map[string]string{
	"gof":    "prints the gof values from 1 up to the specified k",
	"coords": "prints the datasets coordinates for the specified k",
}

func mdsParseParams() *mdsParams {
	params := new(mdsParams)
	params.k =
		flag.Int("k", 2, "the number of the principal coordinates to use - 0 for autosearch")
	params.script =
		flag.String("sc", "", "the script to be used for the MDS eval")
	similaritiesPath :=
		flag.String("sim", "", "the path of the similarity matrix file")
	params.logfile =
		flag.String("l", "", "the logfile to be used")
	params.output =
		flag.String("o", "", "the output file")
	modulesStr :=
		flag.String("m", "coords", "which modules to run (type list to see the options) ")
	flag.Parse()
	setLogger(*params.logfile)

	if *modulesStr == "list" {
		for k, v := range mdsModules {
			fmt.Printf("%s: %s\n", k, v)
		}
		os.Exit(1)
	} else {
		params.modules = make(map[string]bool)
		for _, mod := range strings.Split(*modulesStr, ",") {
			params.modules[mod] = true
		}
	}

	if *similaritiesPath == "" || *params.output == "" || *params.script == "" {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

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
	return params
}

func mdsRun() {
	params := mdsParseParams()

	if params.modules["coords"] {
		extension := filepath.Ext(*params.output)
		basename := strings.TrimSuffix(*params.output, extension)
		outfile, er := os.OpenFile(basename+"-coords"+extension, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if er != nil {
			fmt.Fprintln(os.Stderr, er)
			os.Exit(1)
		}
		defer outfile.Close()

		log.Println("Executing MDS")
		mds := core.NewMDScaling(params.similarities, *params.k, *params.script)
		err := mds.Compute()
		log.Println("Done")
		if err != nil {
			log.Fatalln(err)
		}

		coords := mds.Coordinates()
		if coords == nil {
			return
		}
		buf := core.SerializeCoordinates(coords)
		outfile.Write(buf)

		//		for i := 0; i < len(coords[0]); i++ {
		//			fmt.Fprintf(outfile, "x_%d ", i+1)
		//		}
		//		fmt.Fprintf(outfile, "\n")
		//		for _, d := range coords {
		//			for _, c := range d {
		//				fmt.Fprintf(outfile, "%.5f ", c)
		//			}
		//			fmt.Fprintf(outfile, "\n")
		//		}
	}

	if params.modules["gof"] {
		extension := filepath.Ext(*params.output)
		basename := strings.TrimSuffix(*params.output, extension)
		outfile, er := os.OpenFile(basename+"-gof"+extension, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if er != nil {
			fmt.Fprintln(os.Stderr, er)
			os.Exit(1)
		}
		defer outfile.Close()

		fmt.Fprintf(outfile, "dimensions gof stress\n")
		for k := 1; k <= *params.k; k++ {
			log.Println("Executing MDS for k =", k)
			mds := core.NewMDScaling(params.similarities, k, *params.script)
			err := mds.Compute()
			log.Println("Done")
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Fprintf(outfile, "%d %.5f %5f\n", k, mds.Gof(), mds.Stress())
		}

	}
}
