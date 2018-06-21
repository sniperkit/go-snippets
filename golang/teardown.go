package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func setup() (*os.File, func(), error) {
	teardown := func() {}
	// make a test file
	f, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		return nil, teardown, err
	}
	teardown = func() {
		// close f
		err := f.Close()
		if err != nil {
			return fmt.Errorf("setup: Close: %s", err)
		}
		// delete the test file
		err = os.RemoveAll(f.Name())
		if err != nil {
			fmt.Errorf("setup: RemoveAll: %s", err)
		}
	}
	return f, teardown, nil
}

func main() {
	_, teardown, _ := setup()
	defer teardown()
}
