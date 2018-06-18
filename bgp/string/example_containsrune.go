package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello World"
    r := []rune(s)
    runeToFind := r[0]
    fmt.Printf("The run to find is %#U \n", runeToFind)
    fmt.Println(strings.ContainsRune("HE", runeToFind))
}
