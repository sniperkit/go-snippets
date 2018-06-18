package main

import (
	"flag"
	"fmt"

	"github.com/go-hep/croot"
)

type Det struct {
	E float64
	T float64
}

type Event struct {
	I int64
	A Det
	B Det
}

var evtmax *int64 = flag.Int64("evtmax", 10000, "number of events to read")
var fname *string = flag.String("fname", "event.root", "file to read back")

func tree0(f croot.File) {
	t := f.Get("tree").(croot.Tree)
	e := &Event{}

	t.SetBranchAddress("evt", e)

	// fill some events with random numbers
	nevents := int64(*evtmax)
	if nevents < 0 || nevents > int64(t.GetEntries()) {
		nevents = int64(t.GetEntries())
	}
	for iev := int64(0); iev != nevents; iev++ {
		if iev%1000 == 0 {
			fmt.Printf(":: processing event %d...\n", iev)
		}
		if t.GetEntry(iev, 1) <= 0 {
			panic("error")
		}
		if iev%1000 == 0 {
			fmt.Printf("ievt: %d\n", iev)
			fmt.Printf("evt.a.e= %8.3f\n", e.A.E)
			fmt.Printf("evt.a.t= %8.3f\n", e.A.T)
			fmt.Printf("evt.b.e= %8.3f\n", e.B.E)
			fmt.Printf("evt.b.t= %8.3f\n", e.B.T)
		}
	}
}

func main() {
	flag.Parse()

	croot.RegisterType(&Event{})

	fmt.Printf(":: opening [%s]...\n", *fname)
	f, err := croot.OpenFile(*fname, "read", "my event file", 1, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close("")
	tree0(f)

}

// EOF
