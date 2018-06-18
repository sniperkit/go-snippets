// 当前程序的包名
package main

// 导入包
import (
	"fmt"
	"math"
)

// 省略调用，不能跟别名一同使用
// import . "fmt"

// 首字母小写是private
// 首字母大写是public

// 常量的定义
const PI = 3.14

// 全局变量的声明与赋值
var name = "gopher"

// 一般类型定义声明
type newType int

// 声明结构
type gopher struct{}

// 声明接口
type goland interface{}

// 自定义类型，但最好不要去重新定义内置类型
type (
	byte int8
	rune int32
	文本   string
)

// main函数是程序入口点
func main() {
	// var a string
	// var a [0]int
	// var a [1]int
	// var a [1]bool
	var a [1]byte
	fmt.Println(a)
	fmt.Println(math.MaxFloat32)
	fmt.Println(math.MinInt8)

	var b 文本
	b = "中文"
	fmt.Println(b)

	// 下面四条一样
	// 全局变量不能使用:=
	// var c int
	// c = 123

	// var c int = 1
	// var c = 1
	// c := 1

	// var a, b, c, d int = 1, 2, 3, 4
	// var a, b, c, d = 1, 2, 3, 4
	// a, b, c, d := 1, 2, 3, 4
	// a, _, c, d := 1, 2, 3, 4 // 空白符号

	// Go 不存在隐式转换

	var str1 int = 65
	str2 := string(str1) // strconv.Itoa(65) => "65"
	// strconv.Atoi("65") => int 65
	fmt.Println(str2)
}
