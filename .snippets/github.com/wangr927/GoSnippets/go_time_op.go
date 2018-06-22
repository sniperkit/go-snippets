package main

import (
	"time"
	"fmt"
)

// time.Tick与time.After的使用
func for_select_timetick_timeafter(){
  tick := time.Tick(500 * time.Millisecond)
	boom := time.After(2000 * time.Millisecond)

	for {
		select {
		case <-tick:
			fmt.Println("tick.")
			case <-boom:
				fmt.Println("boom!")
				return
				default:
					fmt.Println("      .")
					time.Sleep(250 * time.Millisecond)
		}
	}
}
