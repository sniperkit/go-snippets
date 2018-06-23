// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
)

func fatal(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if l := len(flag.Args()); l != 1 {
		fatalf("invalid # of arguments, expected 1, found %d", l)
	}

	bs, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fatal(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", string(bs), 0)
	if err != nil {
		panic(err)
	}

	ast.Print(fset, f)
}
