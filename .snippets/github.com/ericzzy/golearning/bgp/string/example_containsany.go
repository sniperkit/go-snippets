package main

import (
    "fmt"
    "strings"
)

func main() {
    fmt.Println(strings.ContainsAny("team", "i"))
    fmt.Println(strings.ContainsAny("failure", "iu"))
    fmt.Println(strings.ContainsAny("seafood", ""))
    fmt.Println(strings.ContainsAny("", ""))
}
