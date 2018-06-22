package main

import (
	"config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lib"
	"mongodatabase"
)

func main() {
	file, err := ioutil.ReadFile("cities.geojson")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(len(file))
	var cities interface{}
	err = json.Unmarshal(file, &cities)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	cc := cities.(map[string]interface{})
	fmt.Println("File Reading Ends")
	m := mongodatabase.Mongo{}
	err = m.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Connected to mongo !")
	i := 0
	for _, c := range cc["features"].([]interface{}) {
		city := c.(map[string]interface{})
		geo := city["geometry"].(map[string]interface{})
		cords := geo["coordinates"].([]interface{})
		prop := city["properties"].(map[string]interface{})
		m.Insert(config.CITIES_DB_COLLECTION, mongodatabase.CityCollection{
			LocationHash: lib.MD5strings(lib.FloatToString(cords[1].(float64)), lib.FloatToString(cords[0].(float64))),
			Lat:          lib.FloatToString(cords[1].(float64)),
			Lng:          lib.FloatToString(cords[0].(float64)),
			Name:         prop["city"].(string),
		})
		i++
		fmt.Println(i, prop["city"].(string))
	}
	fmt.Println("Done , Inserted:", i)
}
