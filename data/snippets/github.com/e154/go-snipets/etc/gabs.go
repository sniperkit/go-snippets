package main

import (
	"github.com/Jeffail/gabs"
	"fmt"
)

func main() {

	jsonParsed, _ := gabs.ParseJSON([]byte(`{
	"outter":{
		"inner":{
			"value1":10,
			"value2":22
		},
		"alsoInner":{
			"value1":20
		}
	},
	"error": "qweqwe"
}`))

	var value float64
	var ok bool

	value, ok = jsonParsed.Path("outter.inner.value1").Data().(float64)
	// value == 10.0, ok == true
	fmt.Println(value, ok)

	value, ok = jsonParsed.Search("outter", "inner", "value1").Data().(float64)
	// value == 10.0, ok == true
	fmt.Println(value, ok)

	value, ok = jsonParsed.Path("does.not.exist").Data().(float64)
	// value == 0.0, ok == false
	fmt.Println(value, ok)

	fmt.Println(jsonParsed.Path("error").Exists())
	fmt.Println(jsonParsed.Path("error").Data().(string))
}
