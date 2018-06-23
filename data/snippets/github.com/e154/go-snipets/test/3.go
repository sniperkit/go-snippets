package main

import (
	"fmt"
	"github.com/e154/go-snipets/test/tools"
)

// 03. Есть следующий пакет, будет ли работать следующий код в пакете main


// Ответы:
// Да

func main() {
	x := tools.New()
	x.N += 1

	fmt.Println(x.N)
}
