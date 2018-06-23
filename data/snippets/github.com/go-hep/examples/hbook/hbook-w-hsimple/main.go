package main

import (
	"compress/flate"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/rio"
)

func main() {
	var err error
	h1, h2 := genHistos()

	err = writeRio("hsimple.rio", h1, h2)
	if err != nil {
		panic(err)
	}

	err = writeGob("hsimple.gob", h1, h2)
	if err != nil {
		panic(err)
	}
}

func genHistos() (*hbook.H1D, *hbook.H1D) {
	h1 := hbook.NewH1D(100, -10., 10.)
	if h1 == nil {
		panic(fmt.Errorf("nil pointer to H1D"))
	}

	h1.Annotation()["name"] = "h1"

	h2 := hbook.NewH1D(100, -10., 10.)
	if h2 == nil {
		panic(fmt.Errorf("nil pointer to H1D"))
	}

	h2.Annotation()["name"] = "h2"

	// draw random values from a gaussian distribution
	const n = 1000
	fmt.Printf(":: filling histos with %d events...\n", n)
	rand.Seed(0)
	for i := 0; i < n; i++ {
		v1 := rand.NormFloat64()
		h1.Fill(v1, 1.0)

		v2 := rand.NormFloat64()
		h2.Fill(v2+2, 1.0)
	}
	fmt.Printf("%s: entries=%v mean=%+8.3f RMS=%+8.3f\n",
		h1.Name(), h1.Entries(), h1.XMean(), h1.XRMS(),
	)
	fmt.Printf("%s: entries=%v mean=%+8.3f RMS=%+8.3f\n",
		h2.Name(), h2.Entries(), h2.XMean(), h2.XRMS(),
	)

	return h1, h2
}

func writeRio(fname string, h1, h2 *hbook.H1D) error {

	fmt.Printf(":: creating output file [%s]...\n", fname)
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	w, err := rio.NewWriter(f)
	if err != nil {
		return err
	}
	defer w.Close()

	err = w.SetCompressor(rio.CompressZlib, flate.DefaultCompression)
	if err != nil {
		panic(err)
	}

	fmt.Printf(":: saving histos to [%s]...\n", fname)
	err = w.WriteValue("h1", h1)
	if err != nil {
		return err
	}

	err = w.WriteValue("h2", h2)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return err
}

func writeGob(fname string, h1, h2 *hbook.H1D) error {
	fmt.Printf(":: creating output file [%s]...\n", fname)
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	fmt.Printf(":: saving histos to [%s]...\n", fname)
	err = enc.Encode(h1)
	if err != nil {
		return err
	}

	err = enc.Encode(h2)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return err

}
