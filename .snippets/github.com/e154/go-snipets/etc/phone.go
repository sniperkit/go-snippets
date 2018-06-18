package main

import (
	"fmt"
	"github.com/dongri/phonenumber"
)

func main() {

	fmt.Println(phonenumber.Parse("+7 913 922 0645", "RU"))
	fmt.Println(phonenumber.Parse("7 913 922 0645", "RU"))
	fmt.Println(phonenumber.Parse("913 922 0645", "RU"))
	fmt.Println(phonenumber.Parse("8 913 922 0645", "RU"))
	fmt.Println(phonenumber.Parse("(913) 922 0645", "RU"))
	fmt.Println(phonenumber.Parse("(913)-922-0645", "RU"))
	fmt.Println(phonenumber.Parse("(913)-922--0645", "RU"))
	fmt.Println(phonenumber.Parse("+7 888 888 88 88", "RU"))

}
