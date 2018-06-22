// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anlzr_test

import "testing"

//var testdata = os.Getenv("GOPATH") + "/src/github.com/DevMine/srcanlzr/testdata/go.json"
var testdata = "../testdata/go.json"

func TestRunAnalyzers(t *testing.T) {
	/*bs, err := ioutil.ReadFile(testdata)
	if err != nil {
		t.Fatal(err)
	}

	p, err := src.Unmarshal(bs)
	if err != nil {
		t.Fatal(err)
	}

	res, err := anlzr.RunAnalyzers(p, anlzr.LoC{})
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%.3f", res.AverageFuncLen) != "8.429" {
		t.Errorf("average_function_length: expected 8.429, found %.3f",
			res.AverageFuncLen)
	}

	if res.MinFuncLen != 1 {
		t.Errorf("min_function_length: expected 1, found %d",
			res.MinFuncLen)
	}

	if res.MaxFuncLen != 413 {
		t.Errorf("max_function_length: expected 413, found %d",
			res.MaxFuncLen)
	}

	if res.MedianFuncLen != 6 {
		t.Errorf("max_function_length: expected 6, found %d",
			res.MedianFuncLen)
	}

	if res.TotalLoC != 142513 {
		t.Errorf("total_loc: expected 142513, found %d",
			res.TotalLoC)
	}*/
}
