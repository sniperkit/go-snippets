package main

import "fmt"

func main() {
  var s string = "hello"
  fmt.Printf("address of s %p\n", &s)

  s = "world"
  fmt.Printf("address of s %p\n", &s)
}
