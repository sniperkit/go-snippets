package basic

import (
	"fmt"
	"reflect"
	"testing"
)

type MyType struct {
	A int `my:"hello"` // <key, value> in tag must keep in the format of [key="value"]
}

func printField() {
	var mt MyType
	vv := reflect.TypeOf(mt)
	// v := vv.Elem()
	v := vv
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := field.Tag
		fmt.Println("field ==> ", field)
		fmt.Println("tag ==> ", tag)
		fmt.Println("[my]tag ==> ", tag.Get("my"))
	}
	fmt.Println()
}

func TestTag(t *testing.T) {
	printField()
}
