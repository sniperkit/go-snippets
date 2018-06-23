package main

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"os"
)

var ps *process.Process

// display memory info
func mem(n int) {
	if ps == nil {
		p, err := process.NewProcess(int32(os.Getpid()))
		if err != nil {
			panic(err)
		}
		ps = p
	}

	mem, _ := ps.MemoryInfoEx()
	fmt.Printf("%d. VMS: %d MB, RSS: %d MB\n", n, mem.VMS>>20, mem.RSS>>20)
}

func main() {
	// 1. the memory info after init
	mem(1)

	// 2. create a 10 * 1024 array
	data := new([10][1024 * 1024]byte)
	mem(2)

	// 3. fill up data
	for i := range data {
		for x, n := 0, len(data[i]); x < n; x++ {
			data[i][n] = 1
		}
		mem(3)
	}
}
