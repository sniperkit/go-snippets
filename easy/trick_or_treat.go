package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func trickOrTreat(v, z, w, h uint) uint {
	if v+z+w == 0 {
		return 0
	}
	return (v*3 + z*4 + w*5) * h / (v + z + w)
}

func main() {
	var v, z, w, h uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(),
			"Vampires: %d, Zombies: %d, Witches: %d, Houses: %d",
			&v, &z, &w, &h)
		fmt.Println(trickOrTreat(v, z, w, h))
	}
}
