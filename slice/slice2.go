package main

import (
  "fmt"
)

func main() {
  slice := []int{10, 20, 30, 40, 50}

  fmt.Println(cap(slice))

  slice1 := slice[1:3]

  fmt.Println(cap(slice1))

  //fmt.Println(slice1[4])
}
