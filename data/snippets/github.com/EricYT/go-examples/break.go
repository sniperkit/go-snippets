package main

import "fmt"

func main() {
	for i := 0; i < 10; i++ {
		// The L just a lable,break to layer L.Not a target
	L:
		for {
			//L:
			for {
				break L
				//        break
			}
			fmt.Println("GO out 1 ", i)
		}
		fmt.Println("GO out ", i)
	}

	// goto test

	var retry int
L1:
	for {
		fmt.Println("retry: ", retry)
		retry++
		if retry > 5 {
			break L1
		}
	}
	fmt.Println("break for retry: ", retry)
}
