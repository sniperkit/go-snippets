package main

import (
	"fmt"
	"time"
)

var myId string = "server"
var testing bool = true

func main() {
	i := 0
	testPrintf("assigning %d partitions", 100)
	for i < 10 {
		testPrintf("Superstep %d", i)
		if (i+1)%3 == 0 {
			testPrintf("Saving Checkpoint %d", i/3)
			if i == 2 {
				testPrintf("assigning %d partitions", 100)
			}
		}
		time.Sleep(1 * time.Second)
		i++
	}
}

func testPrintf(format string, a ...interface{}) {
	if testing {
		fmt.Printf("%s: %s\n", myId, fmt.Sprintf(format, a...))
	}

}
