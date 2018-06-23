package main

import "fmt"
import "encoding/json"

type NestedData struct {
	Nested string `json:"nested"`
}

type Data struct {
	Value       int         `json:"value"`
	NestedValue *NestedData `json:"nested_value"`
}

type Envelope struct {
	Data interface{} `json:"data"`
}

func main() {
	jsonBlob := []byte(`{"data":{"value":1,"nested_value":{"nested": "test"}}}`)
	env := &Envelope{Data: &Data{}}
	err := json.Unmarshal(jsonBlob, env)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(env.Data.(*Data).NestedValue.Nested)
}
