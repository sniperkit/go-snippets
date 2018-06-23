package main

import (
	"fmt"
	"os"
)

func main() {
	debug := true
	key := os.Getenv("RANCHER_ACCESS_KEY")
	secret := os.Getenv("RANCHER_SECRET_KEY")

	if len(key) == 0 || len(secret) == 0 {
		defer fmt.Println("RANCHER_ACCESS_KEY & RANCHER_SECRET_KEY are required!")
		os.Exit(1)
	} else {
		if debug == true {
			fmt.Println("RANCHER_ACCESS_KEY:", key)
			fmt.Println("RANCHER_SECRET_KEY:", secret)
		}
	}
}
