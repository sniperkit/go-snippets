// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"os"
	"testing"
)

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Open("./testdata/small.json")
		if err != nil {
			b.Fatal(err)
		}
		_, _ = Decode(f)
		f.Close()
	}
}
