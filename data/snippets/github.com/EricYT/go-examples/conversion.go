package main

import "fmt"

func main() {
  var i int = 1230
  var f float32 = float32(i)
  // Convert i to a character string interpreting the integer as a unicode value.
  var s string = string(i)

  fmt.Println("i is ", i)
  fmt.Println("f is ", f)
  fmt.Println("s is ", s)
}

