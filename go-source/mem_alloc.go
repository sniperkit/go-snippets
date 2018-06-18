package main

func test() *int {
	x := new(int)
	*x = 0xAABB
	return x
}

func main() {
	println(*test())
}
