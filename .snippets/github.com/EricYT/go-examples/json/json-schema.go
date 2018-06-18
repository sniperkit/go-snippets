package main

import "fmt"
import "github.com/xeipuuv/gojsonschema"

type Test struct {
  A string `json:"a"`
  B int    `json:"b"`
  C []string `json:"c"`
}

func main() {
  jsonStr := `{"a":"string", "b":234}`
  loaderStr := gojsonschema.NewGoLoader(jsonStr)

  loader := gojsonschema.NewStringLoader(`{"type": "object", "properties":{"a":{"type":"string"}, "b":{"type":"integer"}, "c":{"type":"array", "items":{"type":"string"}}}, "require": ["a", "b"]}`)

  jsonStruct := Test {
    A : "string",
    B : 123,
    C : []string{"123", "234"},
  }

  structLoader := gojsonschema.NewGoLoader(jsonStruct)

  fmt.Println(loaderStr)

  result, err := gojsonschema.Validate(loader, structLoader)
  if err != nil {
    fmt.Println("Validate error:", err)
  }

  if result.Valid() {
    fmt.Println("The document is validate")
  } else {
    fmt.Println("The document is error:")
    for _, err1 := range result.Errors() {
      fmt.Printf("- %s\n", err1)
    }
  }

}
