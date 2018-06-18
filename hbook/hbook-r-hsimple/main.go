package main

import (
	"encoding/gob"
	"fmt"
	"os"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/rio"
)

func main() {

	readRio("hsimple.rio")
	readGob("hsimple.gob")
}

func readRio(fname string) {
	fmt.Printf(":: opening input file [%s]...\n", fname)
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r, err := rio.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	fmt.Printf(":: loading histos from [%s]...\n", fname)
	var (
		h1 hbook.H1D
		h2 hbook.H1D
	)

	err = read(r, "h1", &h1)
	if err != nil {
		panic(err)
	}

	err = read(r, "h2", &h2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: entries=%v mean=%+8.3f RMS=%+8.3f\n",
		h1.Name(), h1.Entries(), h1.XMean(), h1.XRMS(),
	)
	fmt.Printf("%s: entries=%v mean=%+8.3f RMS=%+8.3f\n",
		h2.Name(), h2.Entries(), h2.XMean(), h2.XRMS(),
	)

	err = r.Close()
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func readGob(fname string) {
	fmt.Printf(":: opening input file [%s]...\n", fname)
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)

	fmt.Printf(":: loading histos from [%s]...\n", fname)
	var (
		h1 hbook.H1D
		h2 hbook.H1D
	)

	err = dec.Decode(&h1)
	if err != nil {
		panic(err)
	}

	err = dec.Decode(&h2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: entries=%v mean=%+8.3f RMS=%+8.3f\n",
		h1.Name(), h1.Entries(), h1.XMean(), h1.XRMS(),
	)
	fmt.Printf("%s: entries=%v mean=%+8.3f RMS=%+8.3f\n",
		h2.Name(), h2.Entries(), h2.XMean(), h2.XRMS(),
	)

	err = f.Close()
	if err != nil {
		panic(err)
	}
}
func read(r *rio.Reader, name string, ptr interface{}) error {
	var err error

	rec := r.Record(name)
	err = rec.Connect(name, ptr)
	if err != nil {
		return err
	}

	err = rec.Read()
	if err != nil {
		return err
	}

	blk := rec.Block(name)
	if blk == nil {
		return fmt.Errorf("no block [%s] in record [%s]", name, name)
	}

	err = blk.Read(ptr)
	if err != nil {
		return err
	}

	return err
}
