package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		var s string
		fmt.Sscan(scanner.Text(), &s)
		addr1 := make([]int, 4)
		fmt.Sscanf(s, "%d.%d.%d.%d", &addr1[0], &addr1[1], &addr1[2], &addr1[3])

		scanner.Scan()
		fmt.Sscan(scanner.Text(), &s)
		addr2 := make([]int, 4)
		fmt.Sscanf(s, "%d.%d.%d.%d", &addr2[0], &addr2[1], &addr2[2], &addr2[3])

		h := make([]int, 10)
		for i := 0; i < 5; i++ {
			var hi, lo string
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &hi)
			scanner.Scan()
			fmt.Sscan(scanner.Text(), &lo)
			fmt.Sscanf(hi+lo, "%x", &h[i])
		}
		h[6], h[7] = addr1[0]<<8+addr1[1], addr1[2]<<8+addr1[3]
		h[8], h[9] = addr2[0]<<8+addr2[1], addr2[2]<<8+addr2[3]
		cs := h[0] + h[1] + h[2] + h[3] + h[4] + h[6] + h[7] + h[8] + h[9]
		cs = (cs & (1<<16 - 1)) + (cs >> 16)
		cs ^= 1<<16 - 1
		h[5] = cs

		for ix, i := range h {
			if ix > 0 {
				fmt.Printf(" ")
			}
			fmt.Printf("%02x %02x", i>>8, (i & (1<<8 - 1)))
		}
		fmt.Println()

		scanner.Split(bufio.ScanLines)
		scanner.Scan()
		scanner.Split(bufio.ScanWords)
	}
}
