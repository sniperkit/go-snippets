package main

import (
	. "github.com/fiam/gounidecode/unidecode"
	"fmt"
)

func main() {
	fmt.Println(Unidecode("áéíóú")) // Will print aeiou
	fmt.Println(Unidecode("\u5317\u4EB0")) // Will print Bei Jing
	fmt.Println(Unidecode("Κνωσός")) // Will print Knosos
	fmt.Println(Unidecode("Кудряшова Анжела Мефодиевна"))
}
