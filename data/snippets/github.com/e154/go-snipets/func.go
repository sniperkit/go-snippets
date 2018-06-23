package main

import (
	"fmt"
	"math"
	"net/url"
)

func add(x , y int) int {
	return x + y
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

type params struct {
	brand string
	model string
	version string
	engine string
	publication string
}

var startUri string = "http://www.wheel-size.com/auto/get-list/"

func getUri(p params) string {
	params := url.Values{}
	if (len(p.brand) > 0) {
		params.Add("brand", p.brand)
	}
	if (len(p.model) > 0) {
		params.Add("model", p.model)
	}
	if (len(p.version) > 0) {
		params.Add("version", p.version)
	}
	if (len(p.engine) > 0) {
		params.Add("engine", p.engine)
	}
	if (len(p.publication) > 0) {
		params.Add("publication", p.publication)
	}

	baseUrl, _ := url.Parse(startUri)
	baseUrl.RawQuery = params.Encode()
	return baseUrl.String()
}

func sum(nums ...int) {
	fmt.Print(nums, " ")
	total := 0
	for _, num := range nums {
		total += num
	}
	fmt.Println(total)
}

func main() {
	fmt.Printf("Function add return: %d", add(42, 13))

	//////////////////
	//Function values
	//////////////////
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x * x + y * y)
	}
	fmt.Println("Function values")
	fmt.Println(hypot(3,4))

	//////////////////
	//Function closures
	//////////////////
	fmt.Println("Function closures")
	pos, neg := adder(), adder()
	for i := 0; i<10; i++{
		fmt.Println(pos(i), neg(-2*i))
	}

	p := params{}
	p.brand = "brand"
	p.model = "model"
	uri := getUri(p)
	fmt.Printf("uri: %s\n", uri)


	sum(1, 2)
	sum(1, 2, 3)

	nums := []int{1, 2, 3, 4}
	sum(nums...)
}
