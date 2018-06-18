package main

import (
	"fmt"
	"encoding/json"
)

// 31. Каким образом напечатать содержимое сложной структуры данных?

// Ответы:
//

type params struct {
	brand string
	model string
	version string
	engine string
	publication string
}

func main() {
	p := &params{"brand", "model", "version", "engine", "publication"}
	fmt.Printf("%#v\n", p)

	b, _ := json.MarshalIndent(p, " ", "")
	fmt.Println( string(b))

	fmt.Println(p)

	fmt.Printf("%t\n", p)
}

