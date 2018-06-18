package lib

import (
	// "errors"
	"fmt"
	"go/scanner"
	"go/token"
)

type Interface struct {
	Name    string
	Methods Funcs
}

type Interfaces []Interface

func (i *Interface) String() string {
	return fmt.Sprintf("interface:\n\tName:%s\n\tMethods:%v\n", i.Name, i.Methods.String())
}

func (i *Interfaces) String() string {
	var str string
	for _, v := range *i {
		str += v.String()
	}
	return str
}

type Struct struct {
	Name string
}

type Structs []Struct

func (s *Struct) String() string {
	return fmt.Sprintf("struct:\n\tName:%s\n", s.Name)
}

func (s *Structs) String() string {
	var str string
	for _, v := range *s {
		str += v.String()
	}

	return str
}

type Type struct {
	Interfaces Interfaces
	Structs    Structs
}

func (t *Type) GetStructs() Structs {
	return t.Structs
}

func (t *Type) GetInterfaces() Interfaces {
	return t.Interfaces
}

func (t *Type) Scan(tok token.Token, lit string, s *scanner.Scanner) {
	if tok != token.TYPE {
		fmt.Println("is not the type")
		return
	}

	_, tok, lit = s.Scan()
	if tok != token.IDENT {
		fmt.Println("error")
		return
	}
	name := lit

	_, tok, lit = s.Scan()
	switch tok {
	case token.INTERFACE:
		fmt.Println("scan interfaces...")
		t.Interfaces.scan(tok, name, s)
	case token.STRUCT:
		fmt.Println("scan structs...")
		st := Struct{name}
		t.Structs = append(t.Structs, st)
	default:
		fmt.Println("error")
		return
	}
}

const (
	INTERFACE_START = iota
	INTERFACE_METHOD
	INTERFACE_END
)

func (i *Interfaces) scan(tok token.Token, name string, s *scanner.Scanner) {
	inter := Interface{Name: name}
	var lit string
	state := INTERFACE_START
	for {
		switch state {
		case INTERFACE_START:
			if tok == token.LBRACE || tok == token.SEMICOLON {
				state = INTERFACE_METHOD
			}
			if tok == token.RBRACE {
				state = INTERFACE_END
			}
			_, tok, lit = s.Scan()
		case INTERFACE_METHOD:
			// fmt.Println("Interface.scan:token", tok.String(), lit)
			inter.Methods.Scan(tok, lit, s)
			state = INTERFACE_START
		}

		if state == INTERFACE_END {
			break
		}
	}

	*i = append(*i, inter)
	// fmt.Println("get interface", inter.String())
}
