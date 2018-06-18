package main

import "fmt"

func main() {
  s := make([]int, 3)

  // append to the slice tail
  s = append(s, 1)
  s = append(s, 2)
  s = append(s, 3)

  fmt.Println(s)
  // [0 0 0 1 2 3]
}
