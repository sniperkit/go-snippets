// ported from gnu addr2line
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/qeedquan/go-binutils/bfd"
	"github.com/qeedquan/go-binutils/iberty/demangle"
)

var (
	filename      = flag.String("e", "a.out", "executable file")
	prettyPrint   = flag.Bool("p", false, "pretty print")
	demangler     = flag.String("C", "", "use demangling style")
	target        = flag.String("b", "", "set target")
	sectionName   = flag.String("j", "", "set section name")
	unwindInlines = flag.Bool("i", false, "unwind inline")
	withAddresses = flag.Bool("a", false, "show addresses")
	withFunctions = flag.Bool("f", false, "show functions")
	baseName      = flag.Bool("s", false, "strip directory names")

	pc            bfd.VMA
	syms          *bfd.SymbolTable
	found         bool
	xfilename     string
	function      string
	line          int64
	discriminator int64
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("addr2line: ")
	flag.Usage = usage
	flag.Parse()
	if *demangler != "" {
		if demangle.CplusNameToStyle(*demangler) == demangle.Unknown {
			log.Fatalf("unknown demangling style %q", *demangler)
		}
	}

	process(*filename, *sectionName, *target)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: [options] addr ...")
	flag.PrintDefaults()
	os.Exit(2)
}

func ck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func process(file, sect, target string) {
	abfd, err := bfd.Openr(file, target)
	ck(err)
	defer bfd.Close(abfd)

	abfd.SetFlags(abfd.Flags() | bfd.DECOMPRESS)

	// bug in addr2line.c (?) seems to check against success and then erroring out
	// we just ignore the return code here
	bfd.CheckFormat(abfd, bfd.Archive)

	_, err = bfd.CheckFormatMatches(abfd, bfd.Object)
	ck(err)

	var section *bfd.Section
	if sect != "" {
		if section = bfd.GetSectionByName(abfd, *sectionName); section == nil {
			log.Fatalf("%s: cannot find section %s", file, sect)
		}
	}

	slurp(abfd)
	translate(abfd, section)
}

func slurp(abfd *bfd.File) {
	if abfd.Flags()&bfd.HAS_SYMS == 0 {
		return
	}

	var dynamic bool
	storage := bfd.GetSymtabUpperBound(abfd)
	if storage == 0 {
		storage = bfd.GetDynamicSymtabUpperBound(abfd)
		dynamic = true
	}
	if storage < 0 {
		ck(bfd.GetError())
	}

	var symcount int64
	var err error
	syms = bfd.AllocSymbolTable(storage)
	if dynamic {
		symcount, err = bfd.CanonicalizeDynamicSymtab(abfd, syms)
	} else {
		symcount, err = bfd.CanonicalizeSymtab(abfd, syms)
	}
	ck(err)

	if symcount == 0 && !dynamic {
		storage = bfd.GetDynamicSymtabUpperBound(abfd)
		if storage > 0 {
			syms.Free()
			syms = bfd.AllocSymbolTable(storage)
			symcount, _ = bfd.CanonicalizeDynamicSymtab(abfd, syms)
		}
	}

	if symcount <= 0 {
		syms.Free()
		syms = nil
	}
}

func translate(abfd *bfd.File, section *bfd.Section) {
	addr := flag.Args()
	readStdin := len(addr) == 0
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if readStdin {
			if !scanner.Scan() {
				break
			}
			pc, _ = bfd.ScanVMA(scanner.Text(), 16)
		} else {
			if len(addr) == 0 {
				break
			}
			pc, _ = bfd.ScanVMA(addr[0], 16)
			addr = addr[1:]
		}

		if bfd.GetFlavor(abfd) == bfd.TargetElfFlavor {
			sign := bfd.VMA(1 << uint64(bfd.GetArchSize(abfd)-1))
			pc &= (sign << 1) - 1
			if bfd.GetSignExtendVMA(abfd) != 0 {
				pc = (pc ^ sign) - sign
			}
		}

		if *withAddresses {
			fmt.Printf("0x")
			bfd.PrintfVMA(abfd, pc)

			if *prettyPrint {
				fmt.Printf(": ")
			} else {
				fmt.Printf("\n")
			}
		}

		if section != nil {
			findOffsetInSection(abfd, section)
		} else {
			bfd.MapOverSections(abfd, findAddressInSection)
		}

		if !found {
			if *withFunctions {
				if *prettyPrint {
					fmt.Printf("?? ")
				} else {
					fmt.Printf("??\n")
				}
			}
			fmt.Printf("??:0\n")
		} else {
			for {
				if *withFunctions {
					name := function
					if name == "" {
						name = "??"
					} else if *demangler != "" {
						alloc := bfd.Demangle(abfd, name, demangle.ANSI|demangle.PARAMS)
						if alloc != "" {
							name = alloc
						}
					}

					fmt.Printf("%s", name)
					if *prettyPrint {
						fmt.Printf(" at ")
					} else {
						fmt.Printf("\n")
					}
				}

				if *baseName && xfilename != "" {
					xfilename = filepath.Base(xfilename)
				}

				if xfilename == "" {
					xfilename = "??"
				}
				fmt.Printf("%s:", xfilename)
				if line != 0 {
					if discriminator != 0 {
						fmt.Printf("%d (discriminator %d)\n", line, discriminator)
					} else {
						fmt.Printf("%d\n", line)
					}
				} else {
					fmt.Printf("?\n")
				}

				if !*unwindInlines {
					found = false
				} else {
					found, xfilename, function, line = bfd.FindInlinerInfo(abfd)
				}

				if !found {
					break
				}

				if *prettyPrint {
					fmt.Printf(" (inlined by) ")
				}
			}
		}
	}
}

func findOffsetInSection(abfd *bfd.File, section *bfd.Section) {
	if found {
		return
	}

	if section.Flags()&bfd.SEC_ALLOC == 0 {
		return
	}

	size := bfd.GetSectionSize(section)
	if bfd.Size(pc) >= size {
		return
	}

	found, xfilename, function, line, discriminator = bfd.FindNearestLineDiscriminator(abfd, section, syms, pc)
}

func findAddressInSection(abfd *bfd.File, section *bfd.Section) {
	if found {
		return
	}

	if section.Flags()&bfd.SEC_ALLOC == 0 {
		return
	}

	vma := bfd.GetSectionVMA(abfd, section)
	if pc < vma {
		return
	}

	size := bfd.GetSectionSize(section)
	if bfd.Size(pc) >= bfd.Size(vma)+size {
		return
	}

	found, xfilename, function, line, discriminator = bfd.FindNearestLineDiscriminator(abfd, section, syms, pc-vma)
}
