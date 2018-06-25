package main

import (
	"crawler"
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [url]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	inputPtr := flag.String("input", "", "path to saved matrix")
	outputPtr := flag.String("output", "", "path to save matrix")
	maxVisitsPtr := flag.Int("max", 0, "maximum number of visits")
	probabilityPtr := flag.Float64("probability", 0.85, "probability of transition by the link")
	tolerancePtr := flag.Float64("tolerance", 0.0001, "tolerance")
	parallelPtr := flag.Bool("parallel", false, "parallel")
	fileTypePtr := flag.String("filetype", "bin", "type of file (bin or graph)")
	flag.Parse()

	args := flag.Args()

	opts := crawler.NewOptions()
	opts.SameHostOnly = true
	opts.MaxVisits = *maxVisitsPtr
	opts.Tolerance = *tolerancePtr
	opts.FollowingProb = *probabilityPtr
	opts.Parallel = *parallelPtr
	opts.FileType = *fileTypePtr
	c := crawler.NewCrawler(opts)

	if len(args) < 1 && *inputPtr == "" {
		fmt.Println("URL or input file required")
		os.Exit(1)
	}

	var outputf *os.File

	if *outputPtr == "" {
		outputf = nil
	} else {
		outputPath := *outputPtr
		var outputerr error
		outputf, outputerr = os.Create(outputPath)
		if outputerr != nil {
			panic(outputerr)
		}
		defer outputf.Close()
	}

	if len(args) >= 1 {
		c.Run(flag.Arg(0), outputf)
	} else {
		inputPath := *inputPtr
		inputf, inputerr := os.Open(inputPath)
		if inputerr != nil {
			panic(inputerr)
		}
		defer outputf.Close()
		c.ParseMatrix(inputf)
	}

}
