// Read json configuration without goconfigure.
package goconfigure

import (
	"encoding/json"
	"io/ioutil"
)

// Read json to an object from a filename.
func ReadJsonConfigFromFile(fileName string, o *interface{}) error {
	fileRead, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	return ReadJsonConfig(fileRead, o)
}

// Read json to an object from a stream of bytes.
func ReadJsonConfig(fileRead []byte, o *interface{}) error {
	err := json.Unmarshal(fileRead, o)
	if err != nil {
		return err
	}

	return nil
}
