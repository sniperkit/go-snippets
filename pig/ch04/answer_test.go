package main

import (
    "fmt"
)

func change(a []int) []int {
    a[1] = 7
    fmt.Printf("%d\n", a)
    return a
}

func main() {
    a := []int{5, 6, 3, 2, 1}
    fmt.Printf("%d\n", a)
    change(a)
    fmt.Printf("%d\n", a)
}
