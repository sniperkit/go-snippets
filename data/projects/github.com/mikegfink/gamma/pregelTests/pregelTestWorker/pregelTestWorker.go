package main

import (
	"fmt"
	"os"
	"time"
)

var myId string

var testing bool = true

func main() {
	myId = os.Args[1]
	testPrintf("assigned %d partitions [%d-%d]", 10, 10, 19)

	i := 0
	for i < 10 {
		testPrintf("Superstep %d", i)
		if (i+1)%3 == 0 {
			testPrintf("Saving Checkpoint %d", i/3)
			if i == 2 {
				testPrintf("assigned %d partitions [%d-%d]",
					6, 3, 8)
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
