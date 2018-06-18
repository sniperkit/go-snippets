package main

import (
	"os"
	"fmt"
	"encoding/json"
	"strings"
)

var jsonStream = `
{
    "model": [],
    "engine": [],
    "version": [
        {"id": "", "value": "Выберите версию"},
        {"id": "ActiveHybrid 5 (F10)", "value": "ActiveHybrid 5 (F10) 2009->0", "image": "bmw_serie-5_118.png"},
        {"id": "Hatch GT (F07)","value": "Hatch GT (F07) 2009->0", "image": "bmw_serie-5_1925.png"},
        {"id": "Седан (E34)", "value": "Седан (E34) 1987->1996", "image": "bmw_serie-5_128.png"},
        {"id": "Седан (E39)", "value": "Седан (E39) 1995->2003","image": "bmw_serie-5_119.png"},
        {"id": "Седан (E60)", "value": "Седан (E60) 2003->2008", "image": "bmw_serie-5_113.png"},
        {"id": "Седан (F10)", "value": "Седан (F10) 2009->0", "image": "bmw_serie-5_118.png"},
        {"id": "Туринг (E34)", "value": "Туринг (E34) 1991->1996"},
        {"id": "Туринг (E39)", "value": "Туринг (E39) 1997->2004", "image": "bmw_serie-5_112.png"},
        {"id": "Туринг (E61)", "value": "Туринг (E61) 2004->2008", "image": "bmw_serie-5_124.png"},
        {"id": "Туринг (F11)", "value": "Туринг (F11) 2010->0", "image": "bmw_serie-5_127.png"}
    ],
    "tires": [],
    "brand": []
}`

type ModelRecord struct {
	Model 	[]Model		`json:"model"`
	Engine	[]Engine	`json:"engine"`
	Version	[]Version	`json:"version"`
	Tires	[]Tires		`json:"tires"`
	Brand	[]Brand		`json:"brand"`
}
type Engine struct {
	Id		string	`json:"id"`
	Value	string	`json:"value"`
}
type Version struct {
	Id		string	`json:"id"`
	Value	string	`json:"value"`
	Image	string	`json:"image"`
}
type Tires struct {
	Id		string	`json:"id"`
	Value	string	`json:"value"`
}
type Model struct {
	Id		string	`json:"id"`
	Value	string	`json:"value"`
}
type Brand struct {
	Id		string	`json:"id"`
	Value	string	`json:"value"`
}

func main() {

	var m = new(ModelRecord)
	err := json.NewDecoder(strings.NewReader(jsonStream)).Decode(&m)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

//	fmt.Println(m.Version)
	for _, val := range m.Version {
		if(val.Id == "" || val.Value == "val.Value") {
			continue
		}
		var image string
		if(len(val.Image) >0) {
			image = val.Image
		} else {
			image = "none"
		}

		fmt.Printf("%s, %s, %s\n", val.Id, val.Value, image)
	}

}
