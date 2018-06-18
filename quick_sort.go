package main

import (
	"fmt"
)

func qucik(data []int) {
	if len(data) <= 0 {
		return
	}
	mid, i := data[0], 1
	head, tail := 0, len(data)-1
	for head < tail {
		if data[i] > mid {
			data[tail], data[i] = data[i], data[tail]
			tail--
		} else {
			data[head], data[i] = data[i], data[head]
			head++
			i++
		}
	}
	data[head] = mid
	qucik(data[:head])
	qucik(data[head+1:])
}

func main(){
  data := []int{1, 3, 2, 123, 23, 343, 343, 56, 5, 7}
	fmt.Println(data)

	qucik(data)

	fmt.Println(data)

}
