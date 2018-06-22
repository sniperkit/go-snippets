//usr/bin/env go run $0 $@; exit
// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sort-clonedir.go is made to sort the clonedir directory.
// It checks all the folders at the 1st level of the root dir and makes sure
// that they are indeed programming language names (since the schema to store
// in the clone directory is the following: lang/user/repo). Folders which are
// not a programming language's name are moved to clonedir/unknown directory.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	// update lang map with what is in the database before use
	lang := map[string]struct{}{
		"abap":                  struct{}{},
		"actionscript":          struct{}{},
		"ada":                   struct{}{},
		"agda":                  struct{}{},
		"ags script":            struct{}{},
		"alloy":                 struct{}{},
		"antlr":                 struct{}{},
		"apex":                  struct{}{},
		"apl":                   struct{}{},
		"applescript":           struct{}{},
		"arc":                   struct{}{},
		"arduino":               struct{}{},
		"asp":                   struct{}{},
		"aspectj":               struct{}{},
		"assembly":              struct{}{},
		"ats":                   struct{}{},
		"augeas":                struct{}{},
		"autohotkey":            struct{}{},
		"autoit":                struct{}{},
		"awk":                   struct{}{},
		"bison":                 struct{}{},
		"bitbake":               struct{}{},
		"blitzbasic":            struct{}{},
		"blitzmax":              struct{}{},
		"bluespec":              struct{}{},
		"boo":                   struct{}{},
		"brightscript":          struct{}{},
		"bro":                   struct{}{},
		"c":                     struct{}{},
		"c#":                    struct{}{},
		"c++":                   struct{}{},
		"cartocss":              struct{}{},
		"ceylon":                struct{}{},
		"chapel":                struct{}{},
		"cirru":                 struct{}{},
		"clean":                 struct{}{},
		"clips":                 struct{}{},
		"clojure":               struct{}{},
		"cobol":                 struct{}{},
		"coffeescript":          struct{}{},
		"coldfusion":            struct{}{},
		"common lisp":           struct{}{},
		"component pascal":      struct{}{},
		"cool":                  struct{}{},
		"coq":                   struct{}{},
		"crystal":               struct{}{},
		"css":                   struct{}{},
		"cuda":                  struct{}{},
		"cycript":               struct{}{},
		"d":                     struct{}{},
		"dart":                  struct{}{},
		"delphi":                struct{}{},
		"dm":                    struct{}{},
		"dogescript":            struct{}{},
		"dot":                   struct{}{},
		"dylan":                 struct{}{},
		"ec":                    struct{}{},
		"eiffel":                struct{}{},
		"elixir":                struct{}{},
		"elm":                   struct{}{},
		"emacs lisp":            struct{}{},
		"emberscript":           struct{}{},
		"erlang":                struct{}{},
		"f#":                    struct{}{},
		"factor":                struct{}{},
		"fancy":                 struct{}{},
		"flux":                  struct{}{},
		"forth":                 struct{}{},
		"fortran":               struct{}{},
		"frege":                 struct{}{},
		"game maker language":   struct{}{},
		"gams":                  struct{}{},
		"gap":                   struct{}{},
		"gdscript":              struct{}{},
		"glyph":                 struct{}{},
		"gnuplot":               struct{}{},
		"go":                    struct{}{},
		"golo":                  struct{}{},
		"gosu":                  struct{}{},
		"grace":                 struct{}{},
		"grammatical framework": struct{}{},
		"groovy":                struct{}{},
		"hack":                  struct{}{},
		"harbour":               struct{}{},
		"haskell":               struct{}{},
		"haxe":                  struct{}{},
		"hy":                    struct{}{},
		"idl":                   struct{}{},
		"idris":                 struct{}{},
		"igor pro":              struct{}{},
		"inform 7":              struct{}{},
		"io":                    struct{}{},
		"isabelle":              struct{}{},
		"j":                     struct{}{},
		"jasmin":                struct{}{},
		"java":                  struct{}{},
		"javascript":            struct{}{},
		"jsoniq":                struct{}{},
		"julia":                 struct{}{},
		"kotlin":                struct{}{},
		"krl":                   struct{}{},
		"labview":               struct{}{},
		"lasso":                 struct{}{},
		"livescript":            struct{}{},
		"logos":                 struct{}{},
		"lolcode":               struct{}{},
		"lookml":                struct{}{},
		"loomscript":            struct{}{},
		"lsl":                   struct{}{},
		"lua":                   struct{}{},
		"m":                     struct{}{},
		"makefile":              struct{}{},
		"mathematica":           struct{}{},
		"matlab":                struct{}{},
		"max":                   struct{}{},
		"max/msp":               struct{}{},
		"mercury":               struct{}{},
		"mirah":                 struct{}{},
		"monkey":                struct{}{},
		"moonscript":            struct{}{},
		"nemerle":               struct{}{},
		"nesc":                  struct{}{},
		"netlogo":               struct{}{},
		"nimrod":                struct{}{},
		"nit":                   struct{}{},
		"nix":                   struct{}{},
		"nu":                    struct{}{},
		"objective-c":           struct{}{},
		"objective-c++":         struct{}{},
		"objective-j":           struct{}{},
		"ocaml":                 struct{}{},
		"ooc":                   struct{}{},
		"opa":                   struct{}{},
		"openedge abl":          struct{}{},
		"openscad":              struct{}{},
		"ox":                    struct{}{},
		"oxygene":               struct{}{},
		"oz":                    struct{}{},
		"pan":                   struct{}{},
		"papyrus":               struct{}{},
		"parrot":                struct{}{},
		"pascal":                struct{}{},
		"pawn":                  struct{}{},
		"perl":                  struct{}{},
		"perl6":                 struct{}{},
		"php":                   struct{}{},
		"piglatin":              struct{}{},
		"pike":                  struct{}{},
		"pogoscript":            struct{}{},
		"powershell":            struct{}{},
		"processing":            struct{}{},
		"prolog":                struct{}{},
		"propeller spin":        struct{}{},
		"puppet":                struct{}{},
		"purebasic":             struct{}{},
		"pure data":             struct{}{},
		"purescript":            struct{}{},
		"python":                struct{}{},
		"r":                     struct{}{},
		"racket":                struct{}{},
		"ragel in ruby host":    struct{}{},
		"realbasic":             struct{}{},
		"rebol":                 struct{}{},
		"red":                   struct{}{},
		"robotframework":        struct{}{},
		"ruby":                  struct{}{},
		"rust":                  struct{}{},
		"sas":                   struct{}{},
		"scala":                 struct{}{},
		"scheme":                struct{}{},
		"scilab":                struct{}{},
		"self":                  struct{}{},
		"shell":                 struct{}{},
		"shen":                  struct{}{},
		"slash":                 struct{}{},
		"smalltalk":             struct{}{},
		"sourcepawn":            struct{}{},
		"sqf":                   struct{}{},
		"sql":                   struct{}{},
		"squirrel":              struct{}{},
		"standard ml":           struct{}{},
		"stata":                 struct{}{},
		"supercollider":         struct{}{},
		"swift":                 struct{}{},
		"systemverilog":         struct{}{},
		"tcl":                   struct{}{},
		"tex":                   struct{}{},
		"thrift":                struct{}{},
		"turing":                struct{}{},
		"txl":                   struct{}{},
		"typescript":            struct{}{},
		"unrealscript":          struct{}{},
		"vala":                  struct{}{},
		"vcl":                   struct{}{},
		"verilog":               struct{}{},
		"vhdl":                  struct{}{},
		"viml":                  struct{}{},
		"visual basic":          struct{}{},
		"volt":                  struct{}{},
		"webidl":                struct{}{},
		"xbase":                 struct{}{},
		"xc":                    struct{}{},
		"xml":                   struct{}{},
		"xojo":                  struct{}{},
		"xproc":                 struct{}{},
		"xquery":                struct{}{},
		"xslt":                  struct{}{},
		"xtend":                 struct{}{},
		"zephir":                struct{}{},
		"zimpl":                 struct{}{},
		"unknown":               struct{}{},
	}

	flag.Usage = func() {
		fmt.Printf("usage: %s [REPOSITORIES_ROOT_FOLDER]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
	}

	rootDir := flag.Arg(0)
	if rootDir == "" {
		flag.Usage()
	}

	directories, err := ioutil.ReadDir(rootDir)
	if err != nil {
		fatal(err)
	}

	if err := os.MkdirAll(filepath.Join(rootDir, "unknown"), 0755); err != nil {
		fatal(err)
	}

	for _, fi := range directories {
		name := filepath.Join(rootDir, fi.Name())
		if !fi.IsDir() {
			fmt.Println("not a directory (" + name + "), skipping...")
			continue
		}

		if _, ok := lang[fi.Name()]; ok {
			fmt.Println(name + " is a language directory, skipping...")
		} else {
			dest := filepath.Join(rootDir, "unknown", fi.Name())
			fmt.Println(name + " is not a language directory. Moving to " + dest)
			if err := os.Rename(name, dest); err != nil {
				fatal(err)
			}
		}
	}
}

func fatal(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}
