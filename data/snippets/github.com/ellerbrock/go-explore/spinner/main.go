package main

import (
	"strconv"
	"time"

	"github.com/briandowns/spinner"
)

func main() {

	for i := 0; i <= 43; i++ {
		s := spinner.New(spinner.CharSets[i], 100*time.Millisecond)
		s.Suffix = "  id: " + strconv.Itoa(i)
		s.Color("green")
		s.Start()
		time.Sleep(3 * time.Second)
		s.Stop()
	}
}
