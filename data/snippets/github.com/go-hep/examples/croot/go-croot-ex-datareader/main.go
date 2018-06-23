package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-hep/croot"
)

var (
	fname = flag.String("file", "", "path to a ROOT file to read")
	tname = flag.String("tree", "", "name of the TTree to read")
)

func main() {
	flag.Parse()

	if *fname == "" {
		fmt.Fprintf(os.Stderr, "**error** you need to give a ROOT file to read\n")
		os.Exit(1)
	}

	if *tname == "" {
		fmt.Fprintf(os.Stderr, "**error** you need to give a TTree name to read\n")
		os.Exit(1)
	}

	fmt.Printf(":: opening file [%s]...\n", *fname)
	f, err := croot.OpenFile(*fname, "read", "my event file", 1, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error** %v\n", err)
		os.Exit(1)
	}

	defer f.Close("")

	t := f.Get(*tname).(croot.Tree)
	evt, err := NewDataReader(t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error** %v\n", err)
		os.Exit(1)
	}

	evtmax := evt.Tree.GetEntries()
	fmt.Printf(":: entries: %v\n", evtmax)
	for iev := int64(0); iev != evtmax; iev++ {
		do_display := iev%(evtmax/10) == 0
		if do_display {
			fmt.Printf(":: processing event %d...\n", iev)
		}
		if evt.GetEntry(iev) <= 0 {
			fmt.Fprintf(os.Stderr, "**error** processing event %d\n", iev)
			os.Exit(1)
		}
		if do_display {
			fmt.Printf("evt: %v %v\n", evt.NTVars.RunNumber, evt.NTVars.EventNumber)
		}
	}
}
