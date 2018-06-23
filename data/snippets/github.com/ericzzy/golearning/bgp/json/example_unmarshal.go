package main

import (
    "encoding/json"
    "log"
)

func main() {
    var jsonBlob = []byte(`[
        {"Name": "Platypus", "Order": "Monotremata"},
        {"Name": "Quoll", "Order": "Dasyuromorphia"}
    ]`)

    type Animal struct {
        Name, Order string
    }

    var animals []Animal
    err := json.Unmarshal(jsonBlob, &animals)
    if err != nil {
        log.Fatalf("Error is %v", err)
    }
    log.Printf("%+q", animals)
}
