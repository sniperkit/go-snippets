package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello"
    r := []rune(s)[0]
    fmt.Printf("The character is: %q \n", r)
    fmt.Println(strings.IndexRune("WWHO", r))
    
    fmt.Println(strings.IndexRune("Hello", 'H'))
}
