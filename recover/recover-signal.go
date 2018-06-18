package main

import "fmt"

func crash() {
	defer func() {
		if err := recover(); err != nil {
			//panic("crash")
			fmt.Println("crash err:", err)
		}
		fmt.Println("crash recover err is nil")
	}()

	//panic("xxxxx")
	fmt.Println("recover crash normal")
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			panic(err)
		}
		fmt.Println("main recover err is nil")
	}()

	crash()
	fmt.Println("main process normal")
}
