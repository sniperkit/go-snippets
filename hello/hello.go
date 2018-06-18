package main

import (
	"flag"
	"fmt"
)

func main() {

	var (
		name = flag.String("name", "Datalayer", "The given name.")
	)
	flag.Parse()

	fmt.Println("Hello", *name)

	i := 1
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}

}
