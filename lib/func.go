//scan func from the source

package lib

import (
	"fmt"
	"go/scanner"
	"go/token"
)

type Func struct {
	Name     string
	Receiver Value
	Values   Values
	Retruns  Values
}

type Funcs []Func

func (f *Funcs) Scan(tok token.Token, lit string, s *scanner.Scanner) bool {
	funcation := Func{}

	if funcation.scan(tok, lit, s) {
		// fmt.Println("Func.Scan:got", funcation.String())
		*f = append(*f, funcation)
		return true
	}

	return false
}

func (f *Funcs) String() string {
	var str string
	for _, funcation := range *f {
		str += funcation.String()
	}
	return str
}

func (f *Func) String() string {
	str := fmt.Sprintf("func:\n\trececiver:%v\n\tname:%s\n\tvalues:%v\n\treturn:%v", f.Receiver, f.Name, f.Values, f.Retruns)
	return str
}

func (f *Func) scan(tok token.Token, lit string, s *scanner.Scanner) bool {
	// if tok != token.FUNC {
	// 	fmt.Println("Func.scan:not a func")
	// 	return false
	// }

	// tok, lit = scan(s)
	if tok == token.LPAREN {
		f.scanReceiver(tok, lit, s)
		tok, lit = scan(s)
	}
	if tok != token.IDENT {
		fmt.Println("Func.scan:not a name")
		return false
	}
	f.Name = lit

	tok, lit = scan(s)
	if tok != token.LPAREN {
		fmt.Println("Func.scan:not a '('")
		return false
	}
	f.Values = append(f.Values, scanValues(tok, lit, s)...)

	tok, lit = scan(s)
	f.scanReturn(tok, lit, s)
	return true
}

const (
	VALUE_START = iota
	VALUE_NAME
	VALUE_TYPE
	VALUE_END
)

func (f *Func) scanReceiver(tok token.Token, lit string, s *scanner.Scanner) {
	value := Value{}

	if tok != token.LPAREN {
		fmt.Println("Func.scanReceiver:not the '('")
		return
	}

	tok, lit = scan(s)
	if tok != token.IDENT {
		fmt.Println("Func.scanReceiver:is not the name")
		return
	}
	value.Name = lit

	tok, lit = scan(s)
	if tok != token.MUL {
		fmt.Println("Func.scanReceiver:not the '*'")
		value.Type = lit
		return
	}
	value.Type = tok.String()

	tok, lit = scan(s)
	value.Type += lit

	f.Receiver = value
	tok, lit = scan(s)
}

func scanValues(tok token.Token, lit string, s *scanner.Scanner) (values Values) {
	// fmt.Println("--->scan values")
	var value Value
	state := VALUE_START
	for {
		// fmt.Printf("\n%v \t%s \t%q\n", state, tok, lit)

		switch state {
		case VALUE_START:
			if tok == token.LPAREN {
				state = VALUE_NAME
			} else {
				state = VALUE_END
			}
		case VALUE_NAME:
			if tok == token.RPAREN {
				state = VALUE_END
				break
			}
			value.Name = lit
			// fmt.Println(">>>get name:", value.Name)
			state = VALUE_TYPE
		case VALUE_TYPE:
			if tok == token.RPAREN {
				state = VALUE_END
				values = append(values, value)
				break
			}
			if tok == token.COMMA {
				state = VALUE_NAME
				values = append(values, value)
				value.Name = ""
				value.Type = ""
				break
			}
			if tok == token.MUL {
				value.Type = tok.String()
			} else {
				value.Type += lit
				// fmt.Println(">>>get type:", value.Type)

				for index, v := range values {
					if v.Type == "" {
						v.Type = value.Type
						values[index] = v
					}
				}
			}
		}

		if state == VALUE_END {
			break
		}
		if state != VALUE_START {
			_, tok, lit = s.Scan()
		}
		if tok == token.EOF {
			break
		}
	}
	// fmt.Println("-------------", tok)
	return
}

func (f *Func) scanReturn(tok token.Token, lit string, s *scanner.Scanner) {
	var values Values
	// fmt.Println("--->scan return")
	if tok == token.LPAREN {
		values = scanValues(tok, lit, s)
	}

	if tok == token.IDENT {
		value := Value{Type: lit}
		values = append(values, value)
		scan(s)
	}

	f.Retruns = append(f.Retruns, values...)
	return
}

func scan(s *scanner.Scanner) (token.Token, string) {
	_, tok, lit := s.Scan()
	return tok, lit
}
