package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rd := rand.NewSource(time.Now().Unix())
	fmt.Println(rd.Int63())

	rand.Seed(time.Now().Unix())
	fmt.Println(rand.Perm(13))
}
