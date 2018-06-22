package parser

import (
	"fmt"

	"github.com/tidwall/gjson"
)

const jsonString = `{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44},
    {"first": "Roger", "last": "Craig", "age": 68},
    {"first": "Jane", "last": "Murphy", "age": 47}
  ]
}`

const jsonLine = `{"name": "Gilbert", "age": 61}
{"name": "Alexa", "age": 34}
{"name": "May", "age": 57}
{"name": "Deloise", "age": 44}`

// GJSONDemo 展示gjson的用法
// 更多用法参考 https://github.com/tidwall/gjson
func GJSONDemo() {
	result := gjson.Parse(jsonString)
	// 基础用法
	fmt.Println(result.Get("name.last"))       // "Anderson"
	fmt.Println(result.Get("age"))             // 37
	fmt.Println(result.Get("children"))        // ["Sara","Alex","Jack"]
	fmt.Println(result.Get("children.#"))      // >> 3
	fmt.Println(result.Get("children.1"))      // "Alex"
	fmt.Println(result.Get("child*.2"))        // "Jack"
	fmt.Println(result.Get("c?ildren.0"))      // "Sara"
	fmt.Println(result.Get("fav\\.movie"))     // "Deer Hunter"
	fmt.Println(result.Get("friends.#.first")) // ["Dale","Roger","Jane"]
	fmt.Println(result.Get("friends.1.last"))  // "Craig"

	// 条件查询，通过 #[...]或者 #[...]#，支持 ==, !=, <, <=, >, >=
	fmt.Println(result.Get(`friends.#[last=="Murphy"].first`))  // "Dale"
	fmt.Println(result.Get(`friends.#[last=="Murphy"]#.first`)) // ["Dale","Jane"]
	fmt.Println(result.Get(`friends.#[age>45]#.last`))          // ["Craig","Murphy"]
	fmt.Println(result.Get(`friends.#[first%"D*"].last`))       // "Murphy"

	// 处理多条json
	result = gjson.Parse(jsonLine)
	fmt.Println(result.Get("..#"))                 // 4
	fmt.Println(result.Get("..1"))                 // {"name": "Alexa", "age": 34}
	fmt.Println(result.Get("..3"))                 // {"name": "Deloise", "age": 44}
	fmt.Println(result.Get("..#.name"))            // ["Gilbert","Alexa","May","Deloise"]
	fmt.Println(result.Get(`..#[name="May"].age`)) // 57

	// result 转数组
	for _, name := range result.Get("..#.name").Array() {
		fmt.Println(name)
	}

	// unmarshal to a map[string]interface{}
	m, ok := result.Get("..1").Value().(map[string]interface{})
	if ok {
		fmt.Println(m["name"].(string), m["age"].(float64))
	}
}
