package main

import (
	"fmt"
	"time"
	"videos"
)

func main() {
	start := time.Now()
	fmt.Println("fetch starting...")
	// time.Sleep(3000)
	videos.Fetch(1000, "20171228-qicaicms.db")
	fmt.Println("fetch done..")
	end := time.Now()
	diff := end.Sub(start)
	fmt.Println(diff)
}
