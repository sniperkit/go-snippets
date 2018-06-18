package main

import (
    "fmt"
    "encoding/json"
)

func main() {
    b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Hello","World"]}`)
    var m interface{}
    err := json.Unmarshal(b, &m)
    if err != nil {
        return
    }
    f := m.(map[string]interface{})
    for k, v := range f {
        switch vv := v.(type) {
        case string:
            fmt.Println(k, "is string", vv)
        case float64:
            fmt.Println(k, "is int", vv)
        case []interface{}:
            fmt.Println(k, "is Array:")
            for i, e := range vv {
                fmt.Println(i, e)
            }
        default:
            fmt.Println(k, "is of a type I don't know how to handle") 
        }
    }
}
