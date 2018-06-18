// ported from binutils cxxfilt
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/qeedquan/go-binutils/iberty/demangle"
)

const (
	stdsym = "_$."
	hpsym  = "_$.<>#,*&[]:(){}"
)

var (
	nflag = flag.Bool("n", true, "don't strip underscore")
	pflag = flag.Bool("p", false, "don't demangle parameters")
	tflag = flag.Bool("t", false, "demangle types")
	iflag = flag.Bool("i", false, "don't be verbose")
	uflag = flag.Bool("_", false, "strip underscore")
	sflag = flag.String("s", "auto", "set demangling style")
)

var (
	flags = demangle.PARAMS | demangle.ANSI | demangle.VERBOSE
	strip = false
	style = "auto"
	valid string
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("cxxfilt: ")

	parseFlags()
	if flag.NArg() > 0 {
		for _, name := range flag.Args() {
			demangler(name)
			fmt.Println()
		}
	} else {
		readStdin()
	}
}

func isAlnum(ch rune) bool {
	switch {
	case 'a' <= ch && ch <= 'z',
		'A' <= ch && ch <= 'Z',
		'0' <= ch && ch <= '9':
		return true
	}
	return false
}

const (
	eof = -1
)

func getch(r *bufio.Reader) rune {
	ch, err := r.ReadByte()
	if err != nil {
		return eof
	}
	return rune(ch)
}

func readStdin() {
	switch demangle.CurrentStyle() {
	case demangle.Gnu,
		demangle.Lucid,
		demangle.Arm,
		demangle.Java,
		demangle.Edg,
		demangle.Gnat,
		demangle.GnuV3,
		demangle.Dlang,
		demangle.Auto:
		valid = stdsym

	case demangle.HP:
		valid = hpsym

	default:
		log.Fatal("internal error: no alphabet style for current style")
	}

	var buf [32767]byte
	r := bufio.NewReader(os.Stdin)
	for {
		i := 0
		ch := getch(r)
		for ch != eof && isAlnum(ch) || strings.IndexRune(valid, ch) >= 0 {
			if i >= len(buf) {
				break
			}
			buf[i], i = byte(ch), i+1
			ch = getch(r)
		}

		if i > 0 {
			demangler(string(buf[:i]))
		}

		if ch == eof {
			break
		}

		fmt.Printf("%c", ch)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: cxxfilt [options] file ...")
	flag.PrintDefaults()

	fmt.Println("\nsupported demanglers: ")
	engines := demangle.Engines()
	for i := 0; i < len(engines); i++ {
		fmt.Printf(engines[i].Name)
		if i+1 < len(engines) {
			fmt.Printf(", ")
		}
	}
	fmt.Println()
	os.Exit(2)
}

func parseFlags() {
	flag.Usage = usage
	flag.Parse()
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "n":
			strip = false
		case "p":
			flags &^= demangle.PARAMS
		case "t":
			flags |= demangle.TYPES
		case "i":
			flags &^= demangle.VERBOSE
		case "_":
			strip = true
		case "s":
			style = f.Value.String()
		}
	})

	value := demangle.CplusNameToStyle(style)
	if value == demangle.Unknown || demangle.CplusSetStyle(value) == demangle.Unknown {
		log.Fatalf("unknown demangling style %q", style)
	}
}

func demangler(name string) {
	// _ and $ are sometimes found at start of function names
	// in assembler syntax in order to distinguish them from other
	// names (eg register names), so skip them here
	skip := 0
	if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "$") {
		skip++
	}
	if strip && strings.HasPrefix(name, "_") {
		skip++
	}

	result := demangle.Cplus(name[skip:], flags)
	if result == "" {
		fmt.Print(name)
	} else {
		if strings.HasPrefix(name, ".") {
			fmt.Print(".")
		}
		fmt.Print(result)
	}
}
