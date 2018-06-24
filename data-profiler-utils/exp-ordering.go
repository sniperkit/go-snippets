package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/giagiannis/data-profiler/core"
)

type expOrderingParams struct {
	mlScript    *string // script used for approximation
	output      *string // output file path
	repetitions *int    // number of times to repeat experiment
	threads     *int    // number of threads to utilize

	coords []core.DatasetCoordinates // coords of datasets
	scores []float64                 // scores of datasets

	samplingRates []float64 // samplings rates to run
}

func expOrderingParseParams() *expOrderingParams {
	params := new(expOrderingParams)
	params.mlScript =
		flag.String("ml", "", "ML script to use for approximation")
	params.output =
		flag.String("o", "", "output path")
	params.repetitions =
		flag.Int("r", 1, "number of repetitions")
	params.threads =
		flag.Int("t", 1, "number of threads")
	loger :=
		flag.String("l", "", "log file")

	coordsFile :=
		flag.String("c", "", "coordinates file")
	scoresFile :=
		flag.String("s", "", "scores file")
	idxFile :=
		flag.String("i", "", "index file")
	srString :=
		flag.String("sr", "", "comma separated sampling rates")

	flag.Parse()
	setLogger(*loger)
	if *params.mlScript == "" || *params.output == "" || *coordsFile == "" ||
		*scoresFile == "" || *idxFile == "" || *srString == "" {
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// sampling rates parsing
	a := strings.Split(*srString, ",")
	params.samplingRates = make([]float64, 0)
	for i := range a {
		v, err := strconv.ParseFloat(a[i], 64)
		if err == nil {
			params.samplingRates = append(params.samplingRates, v)
		}
	}

	// idx file parsing
	f, err := os.Open(*idxFile)
	if err != nil {
		log.Fatalln(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	idx := make([]string, 0)
	for i, line := range strings.Split(string(buf), "\n") {
		a := strings.Split(line, "\t")
		if len(a) == 2 {
			j, err := strconv.ParseInt(a[0], 10, 32)
			if err != nil || int(j) != i {
				log.Fatalln(err)
			}
			idx = append(idx, a[1])
		}
	}
	f.Close()

	// coordinates file parsing
	f, err = os.Open(*coordsFile)
	if err != nil {
		log.Fatalln(err)
	}
	buf, err = ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	params.coords = make([]core.DatasetCoordinates, 0)
	for i, line := range strings.Split(string(buf), "\n") {
		a := strings.Split(line, " ")
		res := make(core.DatasetCoordinates, 0)
		if i > 0 && len(a) > 0 {
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
				params.coords = append(params.coords, res)
			}
		}
	}
	f.Close()

	// scores
	f, err = os.Open(*scoresFile)
	if err != nil {
		log.Fatalln(err)
	}
	buf, err = ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	scores := core.NewDatasetScores()
	scores.Deserialize(buf)
	params.scores = make([]float64, len(scores.Scores))
	for i, path := range idx {
		params.scores[i] = scores.Scores[path]
	}
	f.Close()

	return params
}

type evalResults struct {
	tau, rho                                            float64
	top10, top25, top50                                 float64
	top2Perc, top5Perc, top10Perc, top25Perc, top50Perc float64
}

func expOrderingRun() {
	// inititializing steps
	params := expOrderingParseParams()
	rand.Seed(int64(time.Now().Nanosecond()))
	output := setOutput(*params.output)
	fmt.Fprintln(output,
		//"%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		"sr", "tau-avg", "tau-perc-0", "tau-perc-25", "tau-perc-50", "tau-perc-75", "tau-perc-100",
		"top10-avg", "top10-perc-0", "top10-perc-25", "top10-perc-50", "top10-perc-75", "top10-perc-100",
		"top25-avg", "top25-perc-0", "top25-perc-25", "top25-perc-50", "top25-perc-75", "top25-perc-100",
		"top50-avg", "top50-perc-0", "top50-perc-25", "top50-perc-50", "top50-perc-75", "top50-perc-100",
		"top2Perc-avg", "top2Perc-perc-0", "top2Perc-perc-25", "top2Perc-perc-50", "top2Perc-perc-75", "top2Perc-perc-20",
		"top5Perc-avg", "top5Perc-perc-0", "top5Perc-perc-25", "top5Perc-perc-50", "top5Perc-perc-75", "top5Perc-perc-50",
		"top10Perc-avg", "top10Perc-perc-0", "top10Perc-perc-25", "top10Perc-perc-50", "top10Perc-perc-75", "top10Perc-perc-100",
		"top25Perc-avg", "top25Perc-perc-0", "top25Perc-perc-25", "top25Perc-perc-50", "top25Perc-perc-75", "top25Perc-perc-100",
		"top50Perc-avg", "top50Perc-perc-0", "top50Perc-perc-25", "top50Perc-perc-50", "top50Perc-perc-75", "top50Perc-perc-100",
		"rho-avg", "rho-perc-0", "rho-perc-25", "rho-perc-50", "rho-perc-75", "rho-perc-100",
	)

	slice := make([]int, len(params.coords))
	for i := 0; i < len(slice); i++ {
		slice[i] = i
	}

	testset := generateSet(slice, params.coords, params.scores)

	executeScript := func(script, trainset, testset string) []float64 {
		cmd := exec.Command(script, trainset, testset)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		result := make([]float64, 0)
		for _, line := range strings.Split(string(out), "\n") {
			val, err := strconv.ParseFloat(line, 64)
			if err == nil {
				result = append(result, val)
			}

		}
		return result
	}

	topKCommon := func(x, y []int, k int) float64 {
		set := make(map[int]bool)
		for i, rank := range x {
			if rank < k {
				set[i] = true
			}
		}
		count := 0
		for j, rank := range y {
			if rank < k && set[j] {
				count++
			}
		}
		return float64(count) / float64(k)
	}

	topKPercNeeded := func(x, y []int, k int) float64 {
		maxActualRank := -1
		for i, rank := range x {
			if rank < k {
				if maxActualRank < y[i] {
					maxActualRank = y[i]
				}
			}
		}
		return float64(maxActualRank) / float64(len(x))
	}

	eval := func(sr float64) *evalResults {
		perm := rand.Perm(len(params.coords))
		trainsetIndexes := perm[0:int(float64(len(perm))*sr)]
		trainset := generateSet(trainsetIndexes, params.coords, params.scores)
		appxScores := executeScript(*params.mlScript, trainset, testset)
		ranksAppx, ranksScores := getRanks(appxScores), getRanks(params.scores)
		res := new(evalResults)
		res.tau = getKendalTau(ranksAppx, ranksScores)
		res.rho = getPearsonRho(ranksAppx, ranksScores)
		res.top10 = topKCommon(ranksAppx, ranksScores, int(float64(len(ranksAppx))*0.1))
		res.top25 = topKCommon(ranksAppx, ranksScores, int(float64(len(ranksAppx))*0.25))
		res.top50 = topKCommon(ranksAppx, ranksScores, int(float64(len(ranksAppx))*0.5))
		res.top2Perc = topKPercNeeded(ranksAppx, ranksScores, int(float64(len(ranksAppx))*0.02))
		res.top5Perc = topKPercNeeded(ranksAppx, ranksScores, int(float64(len(ranksAppx))*0.05))
		res.top10Perc = topKPercNeeded(ranksAppx, ranksScores, int(float64(len(ranksAppx))*0.1))
		res.top25Perc = topKPercNeeded(ranksAppx, ranksScores, int(float64(len(ranksAppx))*0.25))
		res.top50Perc = topKPercNeeded(ranksAppx, ranksScores, int(float64(len(ranksAppx))*0.5))
		os.Remove(trainset)
		return res
	}

	// execute
	for _, sr := range params.samplingRates {
		resultsRho, resultsTau, resultsTop10, resultsTop25, resultsTop50 :=
			make([]float64, 0), make([]float64, 0), make([]float64, 0), make([]float64, 0), make([]float64, 0)
		resultsTop2Perc, resultsTop5Perc, resultsTop10Perc, resultsTop25Perc, resultsTop50Perc :=
			make([]float64, 0), make([]float64, 0), make([]float64, 0), make([]float64, 0), make([]float64, 0)
		done := make(chan *evalResults)
		slots := make(chan bool, *params.threads)
		for i := 0; i < *params.threads; i++ {
			slots <- true
		}

		for i := 0; i < *params.repetitions; i++ {
			go func(done chan *evalResults, slots chan bool, repetition int) {
				log.Printf("[thread-%d] Starting calculation for SR %.2f\n", repetition, sr)
				<-slots
				done <- eval(sr)
				slots <- true
				log.Printf("[thread-%d] Done calculation for SR %.2f\n", repetition, sr)
			}(done, slots, i)
		}
		for i := 0; i < *params.repetitions; i++ {
			v := <-done
			resultsTau = append(resultsTau, v.tau)
			resultsRho = append(resultsRho, v.rho)
			resultsTop10 = append(resultsTop10, v.top10)
			resultsTop25 = append(resultsTop25, v.top25)
			resultsTop50 = append(resultsTop50, v.top50)
			resultsTop2Perc = append(resultsTop2Perc, v.top2Perc)
			resultsTop5Perc = append(resultsTop5Perc, v.top5Perc)
			resultsTop10Perc = append(resultsTop10Perc, v.top10Perc)
			resultsTop25Perc = append(resultsTop25Perc, v.top25Perc)
			resultsTop50Perc = append(resultsTop50Perc, v.top50Perc)
		}
		metricFormat := "%.5f\t%.5f\t%.5f\t%.5f\t%.5f\t%.5f"
		format := "%.5f"
		for i := 0; i < 10; i++ {
			format += "\t" + metricFormat
		}
		format += "\n"
		fmt.Fprintf(output,
			format,
			sr,
			getAverage(resultsTau),
			getPercentile(resultsTau, 0),
			getPercentile(resultsTau, 25),
			getPercentile(resultsTau, 50),
			getPercentile(resultsTau, 75),
			getPercentile(resultsTau, 100),
			getAverage(resultsTop10),
			getPercentile(resultsTop10, 0),
			getPercentile(resultsTop10, 25),
			getPercentile(resultsTop10, 50),
			getPercentile(resultsTop10, 75),
			getPercentile(resultsTop10, 100),
			getAverage(resultsTop25),
			getPercentile(resultsTop25, 0),
			getPercentile(resultsTop25, 25),
			getPercentile(resultsTop25, 50),
			getPercentile(resultsTop25, 75),
			getPercentile(resultsTop25, 100),
			getAverage(resultsTop50),
			getPercentile(resultsTop50, 0),
			getPercentile(resultsTop50, 25),
			getPercentile(resultsTop50, 50),
			getPercentile(resultsTop50, 75),
			getPercentile(resultsTop50, 100),
			getAverage(resultsTop2Perc),
			getPercentile(resultsTop2Perc, 0),
			getPercentile(resultsTop2Perc, 25),
			getPercentile(resultsTop2Perc, 50),
			getPercentile(resultsTop2Perc, 75),
			getPercentile(resultsTop2Perc, 100),
			getAverage(resultsTop5Perc),
			getPercentile(resultsTop5Perc, 0),
			getPercentile(resultsTop5Perc, 25),
			getPercentile(resultsTop5Perc, 50),
			getPercentile(resultsTop5Perc, 75),
			getPercentile(resultsTop5Perc, 100),
			getAverage(resultsTop10Perc),
			getPercentile(resultsTop10Perc, 0),
			getPercentile(resultsTop10Perc, 25),
			getPercentile(resultsTop10Perc, 50),
			getPercentile(resultsTop10Perc, 75),
			getPercentile(resultsTop10Perc, 100),
			getAverage(resultsTop25Perc),
			getPercentile(resultsTop25Perc, 0),
			getPercentile(resultsTop25Perc, 25),
			getPercentile(resultsTop25Perc, 50),
			getPercentile(resultsTop25Perc, 75),
			getPercentile(resultsTop25Perc, 100),
			getAverage(resultsTop50Perc),
			getPercentile(resultsTop50Perc, 0),
			getPercentile(resultsTop50Perc, 25),
			getPercentile(resultsTop50Perc, 50),
			getPercentile(resultsTop50Perc, 75),
			getPercentile(resultsTop50Perc, 100),
			getAverage(resultsRho),
			getPercentile(resultsRho, 0),
			getPercentile(resultsRho, 25),
			getPercentile(resultsRho, 50),
			getPercentile(resultsRho, 75),
			getPercentile(resultsRho, 100),
		)
	}
	os.Remove(testset)
}
