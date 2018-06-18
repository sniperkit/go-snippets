package main

import (
    "fmt"
    "strings"    
)

func main() {
    i := 100
    fmt.Printf("Fields are: %q\n", strings.Fields(" f'oo bar baz "))
    fmt.Printf("Fields are: %q\n", strings.Fields(" 10 100 "))
    fmt.Printf("Fields are: %q", i)
}
