package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/giagiannis/data-profiler/core"
)

const (
	optsDelimiter = ","
	kvDelimiter   = "="
)

func parseOptions(optString string) map[string]string {
	opts := make(map[string]string)
	if optString == "" {
		return opts
	}
	for _, option := range strings.Split(optString, optsDelimiter) {
		sp := strings.Split(option, kvDelimiter)
		if len(sp) != 2 {
			fmt.Fprintln(os.Stderr, "Malformed options")
			os.Exit(1)
		}
		opts[sp[0]] = sp[1]
	}
	return opts
}

func setLogger(logfile string) {
	if logfile != "" {
		f, er := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if er != nil {
			fmt.Println(er)
			os.Exit(1)
		} else {
			log.SetOutput(f)
		}
	}
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

}

func setOutput(output string) *os.File {
	outF := os.Stdout
	if output != "" {
		var err error
		outF, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return outF

}

// returns the average of a float slice
func getAverage(sl []float64) float64 {
	if sl == nil || len(sl) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range sl {
		sum += v
	}
	return sum / float64(len(sl))
}

// returns the n-th percentile of a float slice (50 for median)
func getPercentile(sl []float64, percentile int) float64 {
	sort.Float64s(sl)
	idx := int((float64(percentile) / 100.0) * float64(len(sl)))
	if idx >= len(sl) {
		idx = len(sl) - 1
	}
	return sl[idx]
}

// generate set creates a CSV file containing dataset coordinates and values in
// a comma separated format. useful for train/test set generation. Returns a
// string corresponding to the path of the file.
func generateSet(ids []int, coords []core.DatasetCoordinates, scores []float64) string {
	f, err := ioutil.TempFile("/tmp", "set")
	if err != nil {
		log.Fatalln(err)
	}
	if len(coords) < 1 || len(scores) < 1 {
		log.Fatalln("Coordinates or scores length less than one")
	}
	for i := range coords[0] {
		fmt.Fprintf(f, "x%d,", i)
	}
	fmt.Fprintf(f, "class\n")
	for _, v := range ids {
		for _, c := range coords[v] {
			fmt.Fprintf(f, "%.5f,", c)
		}
		fmt.Fprintf(f, "%.5f\n", scores[v])
	}
	f.Close()
	return f.Name()
}

// getRanks returns the ranks slice for a float slice
func getRanks(a []float64) []int {
	ranks := make([]int, len(a))
	for i := range a {
		for j := range a {
			if a[j] < a[i] {
				ranks[i]++
			}
		}
	}
	return ranks
}

// getKendalTau returns the kendal correlation coefficient for two slices
// representing the ranks of each id (id -> rank)
func getKendalTau(x, y []int) float64 {
	c, nc, n := 0, 0, len(x)
	for i := 0; i < len(x); i++ {
		for j := i + 1; j < len(y); j++ {
			val := (x[i] - x[j]) * (y[i] - y[j])
			if val > 0 {
				c++
			} else if val < 0 {
				nc++
			}
		}
	}
	return float64(c-nc) / float64(n*(n-1)/2)
}

// getPearsonRho returns the Pearson correlation coefficient for two slices
// representing the ranks of each id (id -> rank)
func getPearsonRho(x, y []int) float64 {
	meanA, meanB := 0.0, 0.0
	for i := range x {
		meanA += float64(x[i])
		meanB += float64(y[i])
	}
	meanA = meanA / float64(len(x))
	meanB = meanB / float64(len(y))
	nomin, denom1, denom2 := 0.0, 0.0, 0.0
	for i := range x {
		xC, yC := float64(x[i])-meanA, float64(y[i])-meanB
		nomin += xC * yC
		denom1 += xC * xC
		denom2 += yC * yC
	}
	denom := math.Sqrt(denom1 * denom2)
	if denom == 0 {
		return 0
	}
	return nomin / denom
}

func getDistance(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += (a[i] - b[i]) * (a[i] - b[i])
	}
	return math.Sqrt(sum)
}
