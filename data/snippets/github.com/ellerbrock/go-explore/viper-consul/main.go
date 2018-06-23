package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

func main() {
	viper.AddRemoteProvider("consul", "127.0.0.1:8500", "crypt/config.json")
	viper.SetConfigType("json")
	err := viper.ReadRemoteConfig()

	if err != nil {
		log.Fatal(err)
	}

	firstName := viper.GetString("firstName")
	lastName := viper.GetString("lastName")
	fmt.Println("Name: ", firstName, " ", lastName)
}
