package main

import "time"

//method 1  Using Json File

// conf.json
{
"Users": ["UserA","UserB"],
"Groups": ["GroupA"]
}

import (
"encoding/json"
"os"
"fmt"
)

type Configuration struct {
	Users    []string
	Groups   []string
}

file, _ := os.Open("conf.json")
defer file.Close()
decoder := json.NewDecoder(file)
configuration := Configuration{}
err := decoder.Decode(&configuration)
if err != nil {
fmt.Println("error:", err)
}
fmt.Println(configuration.Users) // output: [UserA, UserB]


// method 2 Using TOML file

// conf.toml
Age = 198
Cats = [ "Cauchy", "Plato" ]
Pi = 3.14
Perfection = [ 6, 28, 496, 8128 ]
DOB = 1987-07-05T05:45:00Z

type Config struct {
	Age int
	Cats []string
	Pi float64
	Perfection []int
	DOB time.Time
}

var conf Config
if _, err := toml.DecodeFile("something.toml", &conf); err != nil {
// handle error
}


// 可以读取多种类型配置文件的模块
// https://github.com/spf13/viper