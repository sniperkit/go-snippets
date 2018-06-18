package lib

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
)

type Tokens struct {
	Funcations Funcs
	Types      Type
}

func (t *Tokens) Scan(path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	src, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)

	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok == token.FUNC {
			fmt.Println("scan funcation...")
			_, tok, lit = s.Scan()
			t.Funcations.Scan(tok, lit, &s)
		}
		if tok == token.TYPE {
			fmt.Println("scan types...")
			t.Types.Scan(tok, lit, &s)
		}
	}
	return
}
