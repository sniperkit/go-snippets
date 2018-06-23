package main

import "fmt"
import "strings"

func main() {
  str := "\t Hello \nworld\n\n"
  str1 := strings.Trim(str, " \t\n\r")
  fmt.Println("original:", str1)
  words := strings.Split(str1, " ")
  for i, word := range words {
    fmt.Printf("%d %s ", i, word)
  }
  fmt.Printf("\n")

}

