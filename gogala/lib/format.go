package lib

import (
	"io/ioutil"
	"os"

	"golang.org/x/tools/imports"
)

func Format(data []byte) ([]byte, error) {

	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}

	filename := tmp.Name()

	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return nil, err
	}

	out, err := imports.Process(filename, data, nil)
	if err != nil {
		return nil, err
	}

	return out, nil
}
