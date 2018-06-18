package golangSnippets

import (
	"io/ioutil"
	"strings"
)

//ReadLineFromFile will take filepath as input and return array of strings as output (without '\n')
func ReadLineFromFile(filepath string) ([]string, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		var emptyArr []string
		return emptyArr, err
	}
	return strings.Split(string(b), "\n"), nil
}
