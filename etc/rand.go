package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	locName = "Asia/Novosibirsk"
)

func main() {
	for i:=0;i<100;i++ {
		fmt.Print(rand.Intn(100),",")
	}

	rand.Seed(time.Now().UnixNano())

	loc, _ := time.LoadLocation(locName)
	d := time.Date(1970 + rand.Intn(40), time.Month(rand.Intn(11) + 1), rand.Intn(30) + 1, 00, 00, 00, 123456000, loc)
	fmt.Println(int(d.Unix()))

	r := rand.Float64()
	fmt.Println(r * 100)
}
