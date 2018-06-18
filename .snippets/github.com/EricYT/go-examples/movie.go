package main


import "fmt"

func main() {
  f1, f2 := 0, 1

  for i := 1; i <= 20; i++ {
    fmt.Printf("%12d %12d", f1, f2)
    if i % 2 == 0 {
      fmt.Printf("\n")
    }
    f1 = f1 + f2
    f2 = f2 + f1
  }
}
