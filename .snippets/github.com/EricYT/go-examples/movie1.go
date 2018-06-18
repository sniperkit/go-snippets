package main

import "fmt"
import "math"

func main() {
  leap := 1
  total:= 0

  fmt.Println("")

  for i := 100; i < 200; i++ {
    k := int(math.Sqrt(float64(i+1)))
    for j := 2; j <= k; j++ {
      if i % j == 0 {
        leap = 0
        break
      }
    }
    if 1 == leap {
      fmt.Printf("%-4d", i)
      if total++; total % 10 == 0 {
        fmt.Println("")
      }
    }
    leap = 1
  }

  fmt.Println("")
  fmt.Println("100 - 200 total ", total)
}
