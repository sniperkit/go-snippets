package main

import "fmt"

func printf(Format string, args ...interface{}) (int, error) {
  _, err := fmt.Printf(Format, args...)
  return len(args), err
}

func main() {
  count := 0

  func() {
    printf("Say hello in the closure %d\n", count)
    count++
  }()

  fmt.Println("The count ", count)
}
