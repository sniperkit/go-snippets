package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

type Spouse struct {
	Name string
	Age  int
}

type Child struct {
	Name string
	Age  int
}

type Profile struct {
	Name string
	Age  int
	Job  string
	Spouse
	Childs []*Child
}

func main() {
	person, err := template.ParseFiles("./profile.yaml")
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("new_one.yaml", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	bob := Child{Name: "bob", Age: 10}
	monica := Child{Name: "monica", Age: 3}
	var profile = Profile{
		Name: "john",
		Age:  30,
		//Job:    "IT engineer",
		Spouse: Spouse{Name: "lucy", Age: 28},
		Childs: []*Child{&bob, &monica},
	}

	person.Execute(file, profile)

	// read from configure file
	viper.SetConfigType("yaml")
	viper.SetConfigFile("./new_one.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Profile name is:%s\n", viper.GetString("name"))
}
