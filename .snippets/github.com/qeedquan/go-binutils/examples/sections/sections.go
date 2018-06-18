// example ported from
// https://geekwentfreak-raviteja.rhcloud.com/blog/2009/11/17/manipulate-elf-files-using-libbfd/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/qeedquan/go-binutils/bfd"
)

func main() {
	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}

	abfd, err := bfd.Openr(flag.Arg(0), "")
	ck(err)
	defer bfd.Close(abfd)

	err = bfd.CheckFormat(abfd, bfd.Object)
	ck(err)

	fmt.Printf("exec file format is: %s\n", abfd.Xvec().Name())
	fmt.Printf("entry point is at address: %#x\n", bfd.GetStartAddress(abfd))

	for s := abfd.Sections(); s != nil; s = s.Next() {
		if s.Flags()&bfd.SEC_LOAD != 0 {
			if s.LMA() != s.VMA() {
				fmt.Printf("loadable section %s: lma = %#x (vma = %#x)  size = %#x\n",
					s.Name(), s.LMA(), s.VMA(), s.Size())
			} else {
				fmt.Printf("loadable section %s: addr = %#x size %#x\n",
					s.Name(), s.LMA(), s.Size())
			}
		} else {
			fmt.Printf("non-loadable section %s: addr %#x size %#x\n", s.Name(), s.VMA(), s.Size())
		}
	}
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: file")
	flag.PrintDefaults()
	os.Exit(2)
}
