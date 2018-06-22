package main


import (
    "github.com/Sayan98/classified"
    "fmt"
)


func main() {
    myStr := "Hi there!"

    fmt.Println("original  : ", myStr)
    fmt.Println("redacted  : ", classified.Redact(myStr))
    fmt.Println("sanitized : ", classified.Sanitize(myStr))
}