package main

import (
    "fmt"
)

func UniqueInts(na []int) []int {
    res := make([]int, 0)
    for index, val := range na {
        exists := false
        for i := 0; i < index; i++ {
            if na[i] == val {
                exists = true
                break
            }
        }
        if !exists {
            res = append(res, val)
        }
    }
    return res
}

func main() {
    slice := []int{9, 1, 9, 5, 4, 4, 2, 1, 5, 4, 8, 8, 4, 3, 6, 9, 5, 7, 5}
    unique := UniqueInts(slice)
    fmt.Println(unique)
}
