package toolkit

import (
	"fmt"
	"strings"
)

var _randChars string

/*
func randChars() string {
	if len(_randChars) == 0 {
		alphabets := "abcdefghijklmnopqrstuvwxyz"
		alphabetsCap := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numerics := "0123456789"
		whitespaces := "!@#$^&*-_+"
		_randChars = alphabets + numerics + alphabetsCap + whitespaces
	}
	return _randChars
}

func SetRandChars(chars string) string {
	_randChars = chars
	return _randChars
}
*/

func RandomString(length int) string {
	return GenerateRandomString("", length)
}

func Sprintf(pattern string, parms ...interface{}) string {
	return fmt.Sprintf(pattern, parms...)
}

func Formatf(pattern string, parms ...interface{}) string {
	r := pattern
	for i, v := range parms {
		r = strings.Replace(r, Sprintf("{%d}", i), v.(string), -1)
	}
	return r
}

func Printf(pattern string, parms ...interface{}) {
	fmt.Printf(pattern, parms...)
}

func Printfn(pattern string, parms ...interface{}) {
	fmt.Println(fmt.Sprintf(pattern, parms...))
}

func Println(s ...interface{}) {
	fmt.Println(s...)
}

func Split(txt string, splitterChars []string) (splitValues []string,
	splitters []string) {
	txtLen := len(txt)
	tmp := ""
	for i := 0; i < txtLen; i++ {
		c := string(txt[i])
		if !HasMember(splitterChars, c) {
			tmp += c
		} else {
			splitValues = append(splitValues, tmp)
			splitters = append(splitters, c)
			tmp = ""
		}
	}
	return
}

func TrimByString(text, groupOfChar string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(groupOfChar, r) < 0 {
			return r
		}
		return -1
	}, text)
}
