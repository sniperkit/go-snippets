package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "https://github.com/ellerbrock"

	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response Type: %T\n", res)

	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	content := string(bytes)
	fmt.Print(content)
}
