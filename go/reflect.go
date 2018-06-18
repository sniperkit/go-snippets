package main

import (
	"fmt"
	"reflect"
	"strings"
)

func main() {
	structReflect()
}

func structReflect() {
	person := struct {
		Name    string `json:"name"`
		Age     int
		married bool
	}{
		"John Doe",
		34,
		false,
	}

	rv := reflect.ValueOf(person)
	for i := 0; i < rv.NumField(); i++ {
		fld := rv.Field(i)
		typeFld := rv.Type().Field(i)
		fmt.Printf("Field #%d : %s\n", i, typeFld.Name)
		fmt.Printf("Type     : %s\n", fld.Type())
		fmt.Printf("Tag      : %s\n", typeFld.Tag)
		fmt.Printf("JSON Tag : %s\n", typeFld.Tag.Get("json"))

		fieldUnexported := typeFld.Name[0:1] == strings.ToLower(typeFld.Name[0:1])
		if fieldUnexported {
			fmt.Printf("Value   : not available (unexported)\n")
		} else {
			fmt.Printf("Value   : %v\n", fld.Interface())
		}

		fmt.Println()
	}
}
