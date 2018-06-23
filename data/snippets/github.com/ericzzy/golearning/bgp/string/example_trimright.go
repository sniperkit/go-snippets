package main

import (
    "fmt"
    "strings"
)

func main() {
    fmt.Printf("[%q]", strings.TrimRight(" ! ! ! Hello ! World !! s!", "! "))
}
