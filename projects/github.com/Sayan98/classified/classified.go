package classified

import (
    "strings"
    "unicode"
    "fmt"
)

/*
Given a string, redact it

Example - "Go Lang" becomes "[REDACTED] [REDACTED]"

Link - https://en.wikipedia.org/wiki/Redaction
*/
func Redact(s string) string {
    redaction := "[REDACTED]"
    str := strings.Fields(s)
    for idx, _ := range str {
        str[idx] = redaction
    }

    return strings.Join(str, " ")
} 

/*
Given a string, sanitize it

Example - "Go Lang" becomes "██ ████"

Link - https://en.wikipedia.org/wiki/Sanitization_(classified_information)
*/
func Sanitize(s string) string {
    sanitization := []rune("█")[0]
    str := []rune(s)
    for idx, runeVal := range str {
        if unicode.IsLetter(runeVal) || unicode.IsNumber(runeVal) {
            str[idx] = sanitization
        }
    }
    return string(str)
}

func main() {
    test := "Go is a programming language created at Google in 2009 by Robert Griesemer, Rob Pike, and Ken Thompson."

    fmt.Println("original  : ", test)
    fmt.Println("readacted : ", Redact(test))
    fmt.Println("sanitized : ", Sanitize(test))
}
