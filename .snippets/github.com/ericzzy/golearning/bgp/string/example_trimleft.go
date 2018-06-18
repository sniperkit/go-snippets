package main

import (
    "fmt"
    "strings"
)

func main() {
    fmt.Printf("[%q]", strings.TrimLeft(" ! ! ! Hello ! World !! s!", "! "))
}
