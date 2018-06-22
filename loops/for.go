package main

import (
  "fmt"
  "strconv"
)

func main() {
  var sum int
  for i := 0; i < 5; i++ {
    sum += i
      fmt.Println(i)
  }

  fmt.Println("sum: "+ strconv.Itoa(sum))
}
