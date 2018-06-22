package main

import (
  "fmt"
  "mystrops"
  "os"
  "strconv"
)

func init() {
    fmt.Println("invoke init before main()")
}

func main() {

  args := len(os.Args)
  if args != 2 {
    fmt.Println("exiting..." + strconv.Itoa(args))
    os.Exit(1)
  }

  output := mystrops.Reverse(os.Args[1])

  fmt.Println("output : " + output)
}
