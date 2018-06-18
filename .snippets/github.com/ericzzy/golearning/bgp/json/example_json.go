package main

import (
    "fmt"
    "encoding/json"
)

type Message struct {
    Name string `json:"name"`
    Body string `json:"body"`
    Time int64  `json:"time"`
    Extra string `json:"extra"`
}

func main() {
    m := Message{Name: "Alice", Body: "Hello", Time: 1294706395881547000}
    b, err := json.Marshal(m)
    if err != nil {
        fmt.Printf("Error is %v \n", err)
        return
    }
    fmt.Printf("json is %s \n", b)
    if string(b) == `{"Name":"Alice","Body":"Hello","Time":1294706395881547000}` {
        fmt.Println("Json string is expected")
    }
    var m2 Message
    err = json.Unmarshal(b, &m2)
    if err != nil {
        fmt.Printf("Error is %v \n", err)
        return
    }
    fmt.Printf("%s", m2)
}
