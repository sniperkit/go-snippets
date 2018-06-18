package main

import (
	"fmt"
	"net/url"
	"strings"
)

var (
//	uri string = "http://www.goodyear.eu/tirefinder/api/by_car/?brand=ABARTH&model=500&version=500 Abarth 2008->0&publication=2151"
//	uri string = "http://www.goodyear.eu/tirefinder/api/by_car/?brand=ABARTH&model=500&version=500+Abarth&publication=2151"
)

func main() {

	params := url.Values{}
	params.Add("brand", "ABARTH")
	params.Add("model", "500")
	params.Add("version", "500 Abarth 2008->0")
	params.Add("publication", "2151")
	baseUrl, _ := url.Parse("http://www.goodyear.eu/tirefinder/api/by_car/image.jpg")
	baseUrl.RawQuery = params.Encode()
	fmt.Println(baseUrl.String())

	// last element
	p := strings.Split(baseUrl.Path, "/")
	fmt.Println(p[len(p)-1])
}
