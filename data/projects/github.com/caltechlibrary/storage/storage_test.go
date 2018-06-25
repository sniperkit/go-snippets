//
// storage package wraps both local disc and S3 storage with CRUD operations and common os.*, ioutil.* functions.
//
// @author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2017, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package storage

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"
)

var (
	testS3 bool
	testGS bool
)

func TestFS(t *testing.T) {
	store, err := Init(FS, nil)
	if err != nil {
		t.Errorf("Init() failed, %s", err)
		t.FailNow()
	}

	// Create a directories  if needed
	err = store.Mkdir("testdata", 0775)
	if err != nil {
		t.Errorf("Can't create testdata directory, %s", err)
		t.FailNow()
	}
	err = store.MkdirAll("testdata/subdir1/subdir2", 0775)
	if err != nil {
		t.Errorf("Can't create testdata/subdir1/subdir2 directory, %s", err)
		t.FailNow()
	}

	fname := path.Join("testdata", "helloworld.txt")
	helloworld := []byte(`Hello World!!!!`)
	err = store.Create(fname, bytes.NewReader(helloworld))
	if err != nil {
		t.Errorf("Create error for %s, %s", fname, err)
		t.FailNow()
	}

	// Stat for FS
	fInfo, err := store.Stat(fname)
	if err != nil {
		t.Errorf("Stat error for %s, %s", fname, err)
		t.FailNow()
	}
	if fInfo == nil {
		t.Errorf("Stat missing fInfo object %ss", fname)
		t.FailNow()
	}
	if fInfo.Name() != path.Base(fname) {
		t.Errorf("Expected %s, got %s", path.Base(fname), fInfo.Name())
	}

	buf, err := store.Read(fname)
	if err != nil {
		t.Errorf("Read error for %s, %s", fname, err)
		t.FailNow()
	}
	if bytes.Compare(buf, helloworld) != 0 {
		t.Errorf("Expected %q, got %q", helloworld, buf)
		t.FailNow()
	}
	helloworld = []byte("Hello World.")
	err = store.Update(fname, bytes.NewReader(helloworld))
	if err != nil {
		t.Errorf("Update error for %s, %s", fname, err)
		t.FailNow()
	}
	// Re-read the data we just wrote out
	buf, err = store.Read(fname)
	if err != nil {
		t.Errorf("Read error for %s, %s after update", fname, err)
		t.FailNow()
	}
	if bytes.Compare(buf, helloworld) != 0 {
		t.Errorf("Expected %q, got %q after update", helloworld, buf)
		t.FailNow()
	}
	// Re-read the data we just wrote out
	buf, err = store.ReadFile(fname)
	if err != nil {
		t.Errorf("ReadFile error for %s, %s after update", fname, err)
		t.FailNow()
	}
	if bytes.Compare(buf, helloworld) != 0 {
		t.Errorf("Expected %q, got %q after update", helloworld, buf)
		t.FailNow()
	}

	// Write file out again
	data := []byte("Hi There")
	err = store.WriteFile(fname, data, 0664)
	if err != nil {
		t.Errorf("WriteFile error for %s, %s", fname, err)
		t.FailNow()
	}
	buf, err = store.ReadFile(fname)
	if err != nil {
		t.Errorf("ReadFile error for %s, %s after update", fname, err)
		t.FailNow()
	}
	if bytes.Compare(data, buf) != 0 {
		t.Errorf("Expected %q, got %q after update", data, buf)
		t.FailNow()
	}

	err = store.Delete(fname)
	if err != nil {
		t.Errorf("Delete error for %s, %s", fname, err)
		t.FailNow()
	}
	// Cleanup if successful so far
	err = store.Remove("testdata/subdir1/subdir2")
	if err != nil {
		t.Errorf("Could not remove testdata/subdir1/subdir2s, %s", err)
	}
	err = store.RemoveAll("testdata")
	if err != nil {
		t.Errorf("Could not remove testdata and it's children, %s", err)
	}
}

func TestCloudStorage(t *testing.T) {
	storeType := UNSUPPORTED
	storeTypes := map[string]bool{
		"S3": testS3,
		"GS": testGS,
	}

	for sLabel, ok := range storeTypes {
		options := map[string]interface{}{}
		switch {
		case sLabel == "FS" && ok:
			storeType = FS
		case sLabel == "S3" && ok:
			storeType = S3
			if s := os.Getenv("AWS_BUCKET"); s != "" {
				options["AwsBucket"] = s
			} else {
				t.Errorf("S3 buckets must be defined before running test")
				t.FailNow()
			}
			if s := os.Getenv("AWS_SDK_LOAD_CONFIG"); s == "1" || strings.ToLower(s) == "true" {
				options["AwsSDKLoadConfig"] = true
				options["AwsSharedConfigEnabled"] = true
				if t := os.Getenv("AWS_PROFILE"); t != "" {
					options["AwsProfile"] = t
				} else {
					options["AwsProfile"] = "default"
				}
			}
		case sLabel == "GS" && ok:
			storeType = GS
			if s := os.Getenv("GOOGLE_PROJECT_ID"); s != "" {
				options["GoogleProjectID"] = s
			}
			if s := os.Getenv("GOOGLE_BUCKET"); s != "" {
				options["GoogleBucket"] = s
			} else {
				t.Errorf("Google Bucket not defined, must be defined before running test")
				t.FailNow()
			}
			if s := os.Getenv("GOOGLE_JSON_CONFIG"); s != "" {
				options["GoogleConfigFile"] = s
			}
		default:
			fmt.Printf("Skipping tests for %s\n", sLabel)
			storeType = UNSUPPORTED
		}

		if storeType != UNSUPPORTED {

			// Now run the tests on the configure storage type
			store, err := Init(storeType, options)
			if err != nil {
				t.Errorf("%s for %s", err, sLabel)
				t.FailNow()
			}
			if store == nil {
				t.Errorf("store was nil for %s", sLabel)
				t.FailNow()
			}

			// Create a directories if needed (s3://, gs:// have no concept of directory so these should NEVER fail)
			err = store.Mkdir("testdata", 0775)
			if err != nil {
				t.Errorf("Can't create testdata directory, %s for ", err, sLabel)
				t.FailNow()
			}
			err = store.MkdirAll("testdata/subdir1/subdir2", 0775)
			if err != nil {
				t.Errorf("Can't create testdata/subdir1/subdir2 directory, %s for %s", err, sLabel)
				t.FailNow()
			}

			fname := `testdata/helloworld.txt`
			expected := []byte(`Hello World!!!`)
			err = store.Create(fname, bytes.NewReader(expected))
			if err != nil {
				t.Errorf("%s for %s", err, sLabel)
				t.FailNow()
			}

			// Stat for Storage Type
			fInfo, err := store.Stat(fname)
			if err != nil {
				t.Errorf("Stat error for %s, %s for %s", fname, err, sLabel)
				t.FailNow()
			}
			if fInfo == nil {
				t.Errorf("Stat missing info object %s for %s", fname, sLabel)
				t.FailNow()
			}
			if fInfo.Name() != path.Base(fname) {
				t.Errorf("expected %s, got %s for %s", path.Base(fname), fInfo.Name(), sLabel)
			}
			if fInfo.Size() != int64(len(expected)) {
				t.Errorf("expected %d, got %d for %s", int64(len(expected)), fInfo.Size(), sLabel)
			}

			// Stat for Storage Type  non-object
			fInfo, err = store.Stat(path.Dir(fname))
			if err == nil {
				t.Errorf("Expected err != nil, fInfo: %+v for %s", fInfo, sLabel)
				t.FailNow()
			}

			result, err := store.Read(fname)
			if err != nil {
				t.Errorf("%s for %s", err, sLabel)
				t.FailNow()
			}
			if bytes.Compare(expected, result) != 0 {
				t.Errorf("expected %q, got %q for %s", expected, result, sLabel)
				t.FailNow()
			}
			expected = []byte(`Hello World.`)
			err = store.Update(fname, bytes.NewReader(expected))
			if err != nil {
				t.Errorf("%s for %s", err, sLabel)
				t.FailNow()
			}
			// Now read back the data and make sure it changed
			result, err = store.Read(fname)
			if err != nil {
				t.Errorf("%s for %s", err, sLabel)
				t.FailNow()
			}
			if bytes.Compare(expected, result) != 0 {
				t.Errorf("expected %q, got %q for %s", expected, result, sLabel)
				t.FailNow()
			}

			data := []byte("Hi There")
			err = store.WriteFile(fname, data, 0664)
			if err != nil {
				t.Errorf("Error WriteFile(%q) %s for %s", fname, err, sLabel)
				t.FailNow()
			}
			buf, err := store.ReadFile(fname)
			if err != nil {
				t.Errorf("Error ReadFile(%q) %s for %s", fname, err, sLabel)
			}
			if bytes.Compare(data, buf) != 0 {
				t.Errorf("expected %q, got %q", expected, result, sLabel)
				t.FailNow()
			}

			// Write a stub file in subdir2 since s3:// and gs:// don't actually make sub-directories
			if err := store.WriteFile("testdata/subdir1/subdir2/hello.txt", data, 0664); err != nil {
				t.Errorf("failed to write test data, %s", err)
			}

			// Cleanup if successful so far
			err = store.Remove(fname)
			if err != nil {
				t.Errorf("delete %s, %s", fname, err)
				t.FailNow()
			}
			err = store.Remove("testdata/subdir1/subdir2/hello.txt")
			if err != nil {
				t.Errorf("Could not remove testdata/subdir1/subdir2s/hello.txt, %s", err)
			}
			err = store.RemoveAll("testdata")
			if err != nil {
				t.Errorf("Could not remove testdata and it's children, %s", err)
			}
		}
	}
}

func TestGetDefaultStore(t *testing.T) {
	// Clear the environment for test.
	opts := map[string]string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "AWS_") == true {
			kv := strings.SplitN(env, "=", 2)
			opts[kv[0]] = kv[1]
			os.Unsetenv(kv[0])
		}
	}
	store := GetDefaultStore()
	if store.Type == UNSUPPORTED {
		t.Errorf("Expected FS type, got UNSUPPORTED")
		t.FailNow()
	}
	if store.Type == S3 {
		t.Errorf("Expected FS type, got S3")
		t.FailNow()
	}
	// See if we have S3 defined
	if len(opts) > 0 {
		for k, v := range opts {
			os.Setenv(k, v)
		}
	}
	if testS3 == true {
		store = GetDefaultStore()
		if store.Type == UNSUPPORTED {
			t.Errorf("Expected S3 type, got UNSUPPORTED")
			t.FailNow()
		}
		if store.Type == FS {
			t.Errorf("Expected S3 type, got FS")
			t.FailNow()
		}
	}
}

func TestWriteFilter(t *testing.T) {
	store := GetDefaultStore()
	start := 0
	finish := 9
	if store.Type == FS {
		_ = os.Mkdir("testdata", 0775)
		defer os.RemoveAll("testdata")
	}
	err := store.WriteFilter("testdata/after-test.txt", func(fp *os.File) error {
		for i := start; i <= finish; i++ {
			fmt.Fprintf(fp, "%d\n", i)
		}
		return nil
	})
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	data, err := store.ReadFile("testdata/after-test.txt")
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	for i, line := range strings.Split(string(data), "\n") {
		if i <= finish && fmt.Sprintf("%d", i) != line {
			t.Errorf("mismatch at line %d, expected %d, got %s", i, i, line)
		}
	}
}

/*
func TestCreateOnExistingS3(t *testing.T) {
	data := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
}
*/

func TestMain(m *testing.M) {
	var all bool

	flag.BoolVar(&all, "all", false, "Run All tests include S3 and GS storage")
	flag.BoolVar(&testS3, "s3", false, "Run S3 storageType tests")
	flag.BoolVar(&testGS, "gs", false, "Run GS storageType tests")
	flag.Parse()
	if all == true {
		testS3 = true
		testGS = true
	}
	os.Exit(m.Run())
}
