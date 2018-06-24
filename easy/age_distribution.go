package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var (
	category []string
	age      []int
)

func init() {
	category = []string{"This program is for humans",
		"Still in Mama's arms",
		"Preschool Maniac",
		"Elementary school",
		"Middle school",
		"High school",
		"College",
		"Working for the man",
		"The Golden Years"}
	age = []int{0, 3, 5, 12, 15, 19, 23, 66, 101}
}

func ageDistribution(m int) string {
	var c int
	for c < len(age) && m >= age[c] {
		c++
	}
	return category[c%len(age)]
}

func main() {
	var m int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &m)
		fmt.Println(ageDistribution(m))
	}
}
