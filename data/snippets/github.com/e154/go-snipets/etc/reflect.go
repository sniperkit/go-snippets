package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func getType(myvar interface{}) (res string) {
	t := reflect.TypeOf(myvar)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		res += "*"
	}
	return res + t.Name()
}

type Test struct {
	Title string	`json:"title"`
}

func test2(i interface{}) {
	fmt.Println(getType(i))
	json.Unmarshal([]byte(`{"title":"test"}`), i)
}

func main() {

	t := &Test{}
	test2(t)

	fmt.Println(t)
}
