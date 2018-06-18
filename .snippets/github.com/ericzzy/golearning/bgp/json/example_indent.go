package main

import (
    "log"
    "encoding/json"
    "bytes"
    "os"
)

func main() {
    type Road struct {
        Name string
        Number int
    }

    roads := []Road{
        {"road1", 2},
        {"road2", 3},
    }

    b, err := json.Marshal(roads)
    if err != nil {
        log.Fatal(err)
    }

    var out bytes.Buffer
    json.Indent(&out, b, "", "    ")
    out.WriteTo(os.Stdout)
}
