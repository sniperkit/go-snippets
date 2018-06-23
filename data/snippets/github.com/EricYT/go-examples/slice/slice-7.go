
package main

import "fmt"
import "time"

func main() {
  s := make([]int, 10000)

  sT1 := time.Now()
  for i := 0; i < len(s); i++ {
    fmt.Println(s[i])
  }
  sT2 := time.Now()
  fmt.Println("Timeused(ms)1: ", sT2.Sub(sT1))

  s1 := s[:]
  sT3 := time.Now()
  for _, value := range s1 {
    fmt.Println(value)
  }
  sT4 := time.Now()
  fmt.Println("Timeused(ms) 2:", sT4.Sub(sT3))

  sSub2 := sT4.Sub(sT3)
  sSub1 := sT2.Sub(sT1)

  fmt.Println("Timeused(ms) 3:", sSub2-sSub1)

}



