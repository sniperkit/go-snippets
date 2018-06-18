package lib

import (
	"fmt"
)

type Value struct {
	Name string
	Type string
}

func (v *Value) IsEqual(value Value) bool {
	return v.Type == value.Type
}

func (v *Value) String() string {
	return fmt.Sprintf("value:\n\tname:%s\n\ttype:%s", v.Name, v.Type)
}

type Values []Value

func (v *Values) IsEqual(values Values) (equal bool) {
	equal = len(*v) == len(values)
	if !equal {
		return
	}

	for index, val := range *v {
		equal = equal && val.IsEqual(values[index])
	}
	return
}

func (v *Values) String() string {
	var str string
	for _, value := range *v {
		str += value.String()
	}
	return str
}
