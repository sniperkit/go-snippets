package main

import (
    "fmt"
    "strings"
)

func main() {
    var s string = "Hello world"
    fmt.Println(strings.TrimSuffix(s, "world"))
    fmt.Println(s)
}
