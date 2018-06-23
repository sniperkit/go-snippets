package main

import "fmt"

func modify(arr [4]int) {
  arr[2] = 31
  fmt.Println("In modify:", arr)
}

func modify1(arr *[4]int) {
  arr[2] = 32
  fmt.Println("In modify:", *arr)
}

func main() {
  array := [4]int{1, 2, 3, 4}

  modify(array)
  fmt.Println("Modify array : ", array)

  modify1(&array)
  fmt.Println("Modify array1 : ", array)

  arr1 := [...]int{1, 2, 3, 4}
  fmt.Println("Array : ", arr1)
}
