package main

import "fmt"

func sum(values []int, result chan<- int) {
  total := 0
  for _, value := range values {
    total += value
  }
  result <- total
}

func main() {
  values := []int{1, 2, 3, 4, 5, 6}

  resultChan := make(chan int, 2)

  go sum(values[:3], resultChan)
  go sum(values[3:], resultChan)

  sum1, sum2 := <-resultChan, <-resultChan

  fmt.Println("Sum1 :", sum1, " Sum2 :", sum2)
}

