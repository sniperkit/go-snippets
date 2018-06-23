package main

import (
	"fmt"
	"reflect"

	"encoding/json"

	"github.com/EricYT/go-examples/msgpack/proto"
)

func main() {
	fmt.Println("msgpack test start")

	actions := make(map[string]proto.Runer)

	t1 := proto.Oper1{
		Key:   "foo",
		Value: "hello",
	}
	actions["foo"] = t1

	t2 := proto.Oper2{
		Key:   "bar",
		Value: "world",
	}
	actions["bar"] = t2

	foo := proto.Foo{
		Actions: actions,
	}

	//data, err := foo.MarshalMsg(nil)
	data, err := json.Marshal(&foo)
	if err != nil {
		panic(fmt.Sprintf("MarshalMsg error:%s", err))
	}
	fmt.Printf("MarshalMsg data:%s\n", string(data))

	var acs proto.Foo
	err = json.Unmarshal(data, &acs)
	if err != nil {
		panic(fmt.Sprintf("UnmarshalMsg error:%s", err))
	}

	fmt.Printf("acs:%+v\n", acs.Actions)
	tt := acs.Actions["foo1"].(proto.Oper1)
	fmt.Printf("acs foo:%+v\n", tt)
	for key, value := range acs.Actions {
		fmt.Printf("key:%s value:%+v\n", key, value)
		switch key {
		case "foo":
			t1_, ok := value.(proto.Oper1)
			fmt.Printf("t1_ type:%s\n", reflect.TypeOf(t1_))
			if !ok {
				panic(fmt.Sprintf("key:%s convert error", key))
			}
			fmt.Println("Oper1 key:", t1_.Key)
		case "bar":
			t2_, ok := value.(proto.Oper2)
			if !ok {
				panic(fmt.Sprintf("key:%s convert error", key))
			}
			fmt.Println("Oper2 key:", t2_.Key)
		}
	}

}
