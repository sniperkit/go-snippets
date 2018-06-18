package fileutil_test

import (
	"fmt"
	"io/ioutil"

	"github.com/uoregon-libraries/gopkg/fileutil"
)

func ExampleSafeFile() {
	var testOut = []byte("This is a test.\n\nA what?\n\nA test.\n\nA what?\n\nA test.\n\nOh, a test.\n")
	var fname = "/tmp/blah.txt"
	var f = fileutil.NewSafeFile(fname)
	f.Write(testOut)

	// At this point there's no data in the file
	var data, err = ioutil.ReadFile(fname)
	fmt.Printf("data: %q; err: %q\n\n", data, err)

	// Closing twice shouldn't cause problems
	f.Close()
	f.Close()

	// Now that the SafeFile has been closed, we have data
	data, err = ioutil.ReadFile(fname)
	if err != nil {
		fmt.Printf("Error trying to read %q: %s\n\n", fname, err)
	}
	fmt.Println(string(data))

	f = fileutil.NewSafeFile(fname)
	f.Write(testOut)
	f.Cancel()
	data, err = ioutil.ReadFile(fname)
	fmt.Println(err)

	// Output:
	// data: ""; err: "open /tmp/blah.txt: no such file or directory"
	//
	// This is a test.
	//
	// A what?
	//
	// A test.
	//
	// A what?
	//
	// A test.
	//
	// Oh, a test.
	//
	// open /tmp/blah.txt: no such file or directory
}
