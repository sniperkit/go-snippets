package main

import "fmt"

func main() {
  str := "Étoilé"

  //Don't do this
  for i := 0; i < len(str); i++ {
    fmt.Printf("%c", str[i])
  }

  fmt.Printf("\n")
  for _, value := range str {
    fmt.Printf("%c", value)
  }
  fmt.Printf("\n")
}
