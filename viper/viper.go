package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {
	fmt.Println("viper tesing start")

	viper.SetConfigType("yaml")
	viper.SetConfigFile("./viper.conf")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("read config error:", err)
		return
	}

	name := viper.GetString("name")
	fmt.Println("name is:", name)

	hobbies := viper.GetStringSlice("hobbies")
	fmt.Println("hobbies is:", hobbies)

	book := viper.GetStringMapString("book")
	fmt.Println("Book is:", book)
	fmt.Println("Book name:", book["name"])
	fmt.Println("Book price:", book["price"])

	return
}
