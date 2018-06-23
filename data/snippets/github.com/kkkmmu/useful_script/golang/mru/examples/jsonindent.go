package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
)

func main() {
	file, err := ioutil.ReadFile("asset/db/common/COMMON.json")
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer

	json.Indent(&buffer, file, "", "\t")
	buffer.WriteTo(os.Stdout)
}
