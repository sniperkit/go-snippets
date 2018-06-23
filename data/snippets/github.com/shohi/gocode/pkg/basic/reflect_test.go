package basic

import (
	"fmt"
	"reflect"
	"testing"
)

type Input struct {
	_ struct{} `type:"structure"`
	A string   `location:"us" type:"string"`
}

func TestReflect(t *testing.T) {
	input := &Input{A: "hello"}

	v := reflect.ValueOf(input)

	fmt.Println(v.Type())
	fmt.Println(v.Elem())
	fmt.Println(v.Interface())

	// print field info
	for k := 0; k < v.Elem().NumField(); k++ {
		fieldInfo := v.Elem().Type().Field(k)

		// Get field metadata
		// fmt.Println(fieldInfo)
		fmt.Println(fieldInfo.Tag.Get("type"))
		fmt.Println(fieldInfo.Name)
		fmt.Println(fieldInfo.Index)
		fmt.Println(fieldInfo.Offset)

		// Get field value
		fmt.Println(v.Elem().Field(k))
	}
}
