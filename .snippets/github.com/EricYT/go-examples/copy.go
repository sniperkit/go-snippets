package main

import (
  "fmt"
  "math/rand"
  "sort"
)

func main() {
  source := make([]int, 10)

  for i := 0; i < 10 ; i++ {
    source[i] = int(rand.Int31n(1000))
  }

  sort.Ints(source)

  worst := make([]int, 5)

  copy(worst, source)
  fmt.Println(worst)
  clean(worst)

  copy(worst[2:3], source)
  fmt.Println(worst)
  clean(worst)

  copy(worst[1:3], source)
  fmt.Println(worst)
}

func clean(source []int) []int {
  length := len(source)
  for i := 0; i < length ; i++ {
    source[i] = 0
  }
  return source
}
