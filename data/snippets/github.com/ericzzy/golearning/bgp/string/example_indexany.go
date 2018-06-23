package main

import (
    "fmt"
    "strings"
)

func main() {
    fmt.Println(strings.IndexAny("chicken", "haenic")) //strings.IndexAny is used to get the the first instance of any code point in the second string.
}
