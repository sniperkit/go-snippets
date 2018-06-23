package main

import (
	"fmt"
	"reflect"
)
import "encoding/json"

type Response1 struct {
	Page   int
	Fruits []string
}

type Response2 struct {
	Page   int        `json:"page"`
	Fruits []string   `json:"fruits"`
	R1     *Response1 `json:"response"`
	//R1     *Response1 `json:"-"`
}

type Foo struct {
	Fb []Bar `json:""`
}

type Bar struct {
	A string `json:"a"`
	B string `json:"b"`
}

type Test struct {
	M map[string]interface{} `json:"m"`
}

func main() {
	slcE := []string{"abc", "def", "xyz"}
	slcD, _ := json.Marshal(slcE)
	fmt.Println(string(slcD))

	mapA := map[string]int{"abc": 1, "def": 2}
	mapC, _ := json.Marshal(mapA)
	fmt.Println("map:", string(mapC))

	res1 := &Response1{
		Page:   1,
		Fruits: []string{"abc", "def"}}
	resJ, _ := json.Marshal(res1)
	fmt.Println(string(resJ))

	res2 := &Response2{
		Page:   2,
		Fruits: []string{"def", "ghi"},
		R1:     res1,
	}
	resJ2, _ := json.Marshal(res2)
	fmt.Println(string(resJ2))

	fb := []Bar{
		Bar{
			A: "a",
			B: "b",
		},
		Bar{
			A: "a1",
			B: "b1",
		},
	}
	_ = &Foo{
		Fb: fb,
	}
	resJ3, _ := json.Marshal(&fb)
	fmt.Println(string(resJ3))

	m := make(map[string]interface{})
	m["a"] = fb
	m["b"] = "yyyy"

	t := Test{
		M: m,
	}

	r, _ := json.Marshal(&t)
	fmt.Println("t:", string(r))

	var t1 = Test{}
	json.Unmarshal(r, &t1)
	fmt.Printf("t Unmarshal:%+v\n", t1)
	fmt.Println("t type:", reflect.TypeOf(t1.M["a"]))
	b, _ := t.M["a"].([]Bar)
	fmt.Printf("b: %+v\n", b[1].A)

	a := 5
	switch {
	case a > 1:
	case a > 2:
	case a > 3:
		fmt.Println("shit2")
	default:
		fmt.Println("fuck3")
	}

}
