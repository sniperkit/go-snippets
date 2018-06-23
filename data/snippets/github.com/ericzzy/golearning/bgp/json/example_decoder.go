package main

import (
    "encoding/json"
    //"fmt"
    "log"
    "io"
    "strings"
)

func main() {
    const jsonStream string = `
        {"Name": "Ed", "Text": "Knock knock."}
        {"Name": "Sam", "Text": "Who's there?"}
    `
    log.Println(jsonStream)

    type Message struct {
        Name string
        Text string
    }

    dec := json.NewDecoder(strings.NewReader(jsonStream))
    for {
        var m Message
        if err := dec.Decode(&m); err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }
        log.Printf("%s: %s\n", m.Name, m.Text)
    }
}
