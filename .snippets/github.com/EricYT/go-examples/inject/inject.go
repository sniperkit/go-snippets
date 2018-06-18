package main

import (
	"fmt"

	"github.com/codegangsta/inject"
)

type Voider interface {
	Void()
}

type SpecialString interface{}

type MyStruct struct {
	Key   string        `inject:"key"`
	Value SpecialString `inject:"value"`
	Index int64
}

func (m MyStruct) Void() {}

type NoMyStruct struct{}

func main() {
	injector := inject.New()
	key := "class"
	injector.Map(key)
	value := "MyStruct"
	injector.MapTo(value, (*SpecialString)(nil))

	fmt.Printf("injector is %+v\n", injector)

	var void interface{} = MyStruct{}
	//var void interface{} = NoMyStruct{}

	if _, ok := void.(Voider); !ok {
		panic("void can not convert to Voider")
	} else {
		err := injector.Apply(&void)
		if err != nil {
			fmt.Println("injector apply error:", err)
			return
		}
		fmt.Printf("MyStruct void is %+v\n", void)
		return
	}
}
