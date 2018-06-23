package main

import (
    "fmt"
    "strings"
)

func main() {
    fmt.Printf("[%q]", strings.Trim(" ! ! ! Hello ! World !! s!", "! "))
}
