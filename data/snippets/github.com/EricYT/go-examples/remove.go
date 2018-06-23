package main

import (
  "fmt"
)

func removeAtIndex(source []int, index int) []int {
  lastIndex := len(source) - 1
  source[index], source[lastIndex] = source[lastIndex], source[index]
  return source[:lastIndex]
}

func main() {
  source := []int{1, 2, 3, 4, 5, 6,}

  source = removeAtIndex(source, 3)

  fmt.Println(source)

}
