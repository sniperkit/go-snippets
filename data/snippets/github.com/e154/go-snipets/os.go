package main

import (
	"os"
	"fmt"
)

func main()  {

	name, err := os.Hostname()
	if err == nil {
		fmt.Println(name)
	}
}