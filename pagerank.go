package main

import (
	"bytes"
	crand "crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"

	"github.com/ttlaak/webgraph"
)

// this function gets the filenames from the current directory (.), if
// it can't open it or read through it for all the filenames it gives
// an error
func getFilenamesFromWD() ([]string, error) {
	wd, err := os.Open(".")
	if err != nil {
		return nil, err
	}
	defer wd.Close()
	fnames, err2 := wd.Readdirnames(-1)
	if err2 != nil {
		return nil, err2
	}
	return fnames, nil
}

// this function takes the filenames from the current directory and
// prepares a webgraph data structure for calculation of the PageRank
func setupGraph(fnames []string) (wg *webgraph.Graph, err error) {
	wg = webgraph.New(len(fnames))
	// add nodes
	for _, fname := range fnames {
		if node_err := wg.AddNode(fname); node_err != nil {
			return nil, node_err
		}
	}
	if *beVerbose {
		fmt.Println("Webgraph: Nodes added")
	}

	wg.Fixate()
	if *beVerbose {
		fmt.Printf("Webgraph: Fixated\nWebgraph: Adding arrows\n")
	}

	// add arrows
	percentage := 0.0
	for i, fname := range fnames {
		bs, rerr := ioutil.ReadFile(fname)
		if rerr != nil {
			return nil, rerr
		}

		neighbours := bytes.Split(bs, []byte("\n"))
		for _, nb := range neighbours[:len(neighbours)-1] {
			aerr := wg.AddArrow(fname, string(nb))
			if aerr != nil {
				return nil, aerr
			}
		}

		if *beVerbose {
			progress := float64(i+1) / (float64(len(fnames)))
			if progress >= percentage {
				fmt.Printf("%4.0f%%", 100.0*percentage)
				os.Stdout.Sync()
				percentage += 0.1
			}
		}
	}
	if *beVerbose {
		fmt.Printf("  Done\nCalculating default weights\n")
	}

	wg.CalculateDefaultWeights()

	return
}

// Two vector operations coming up.
// First taking the L1Norm
func L1Norm(vec []float64) (sum float64) {
	for _, x := range vec {
		sum += math.Abs(x)
	}
	return
}

// modifies its argument by dividing each element by the L1Norm
func Normalize(vec []float64) {
	n := L1Norm(vec)
	for i := range vec {
		vec[i] /= n
	}
}

// writes to w a vector of %g's, a %g is either a format for a fully printed
// floating point number (e.g. 0.00001) or a format for the scientific
// notation of a floating point number (e.g. 2.3e-5) whichever is
// the shortest.
func WriteVec(w io.Writer, vec []float64) {
	fmt.Fprintf(w, "%g", vec[0])
	for _, x := range vec[1:] {
		fmt.Fprintf(w, " %g", x)
	}
	fmt.Fprintf(w, "\n")
}

// See "PageRank: Standing on the shoulders of giants" for more
// and read the english wikipedia page about PageRank
func PageRank(wg *webgraph.Graph) error {
	const epsilon float64 = 0.0001
	delta := 1.0

	Rold := make([]float64, wg.FixedLength())

	f, e := os.Create(outFileName)
	if e != nil {
		return e
	}
	defer f.Close()

	for i := 0; i < len(Rold); i++ {
		Rold[i] = rand.Float64()
	}
	Normalize(Rold)
	WriteVec(f, Rold)

	iteration := 1
	for delta > epsilon {
		Rnew, merr := wg.Multiply(Rold)
		if merr != nil {
			return merr
		}

		// delta <- L2Norm(Rnew - Rold)
		delta = 0.0
		for k := 0; k < len(Rnew); k++ {
			diff := Rnew[k] - Rold[k]
			delta += diff * diff
		}
		delta = math.Sqrt(delta)

		WriteVec(f, Rnew)

		if len(Rold) != copy(Rold, Rnew) {
			return errors.New("PageRank: copy(Rold, Rnew) didn't copy all elements")
		}

		if *beVerbose {
			if iteration == 1 {
				fmt.Printf("iteration: 1")
			} else {
				fmt.Printf("..%d", iteration)
			}
			os.Stdout.Sync()
		}
		iteration++
	}
	if *beVerbose {
		fmt.Printf("\n")
	}

	return nil
}

// reads from the os his random number generator to seed math/rand his
// generator
func seedRand() error {
	p := make([]byte, 8)
	_, e := crand.Read(p) // use crypto rand here because platform independent
	if e != nil {
		return e
	}

	var u uint64 = 0
	for _, v := range p {
		u <<= 8
		u += uint64(v)
	}
	rand.Seed(int64(u))

	return nil
}

var (
	outFileName   string
	beVerbose     *bool   = flag.Bool("v", false, "be verbose")
	printMatrix   *bool   = flag.Bool("m", false, "print matrix to stdout")
	orderFileName *string = flag.String("o", "", "names matching the output values on outFile lines")
)

func commandLine() error {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [-v] [-m] [-o <orderFileName>] <outFileName>\n"+
			"  -m: print the Google matrix\n"+
			"  -o <orderFileName>: names matching the output values on outFile lines\n"+
			"  -v: be verbose\n", os.Args[0])
		return errors.New("commandLine: not enough arguments")
	}
	outFileName = flag.Arg(0)
	return nil
}

// where it all comes together and some code to optionally print out the
// underlying matrix or make an orderFile
func main() {
	if clerr := commandLine(); clerr != nil {
		fmt.Fprintln(os.Stderr, clerr)
		return
	}

	srerr := seedRand()
	if srerr != nil {
		fmt.Fprintln(os.Stderr, srerr)
		return
	}
	if *beVerbose {
		fmt.Println("Successfully seeded random generator from /dev/urandom")
	}

	fnames, err := getFilenamesFromWD()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if *beVerbose {
		fmt.Printf("Found %d filenames in the current working directory\n", len(fnames))
	}

	wg, wgerr := setupGraph(fnames)
	if wgerr != nil {
		fmt.Fprintln(os.Stderr, wgerr)
		return
	}
	if *beVerbose {
		fmt.Println("Graph setup successfull")
	}

	if *printMatrix {
		fmt.Printf("Google matrix:\n")
		for i := 0; i < wg.FixedLength(); i++ {
			row, IWVerr := wg.IncomingWeightsVector(wg.Idx2key(i))
			if IWVerr != nil {
				fmt.Fprintln(os.Stderr, IWVerr)
				return
			}
			fmt.Printf("%9.7f", row[0])
			for _, x := range row[1:] {
				fmt.Printf(" %9.7f", x)
			}
			fmt.Printf("\n")
		}
	}

	// In the outFile rows have elements corresponding to the filenames we
	// initially found in the current directory. To let the user know in
	// which order, he can get an orderFile, the next part writes that to
	// disk.
	if *orderFileName != "" {
		of, oferr := os.Create(*orderFileName)
		if oferr != nil {
			fmt.Fprintln(os.Stderr, "Could not open orderFile")
			fmt.Fprintln(os.Stderr, oferr)
			return
		}

		for i := 0; i < wg.FixedLength(); i++ {
			fmt.Fprintln(of, wg.Idx2key(i))
		}

		of.Close()
		if *beVerbose {
			fmt.Printf("orderFile written to %s\n", *orderFileName)
		}
	}

	if *beVerbose {
		fmt.Println("Starting PageRank")
	}
	prerr := PageRank(wg)
	if prerr != nil {
		fmt.Fprintln(os.Stderr, prerr)
		return
	}
	if *beVerbose {
		fmt.Println("PageRank ended successfully")
	}
}
