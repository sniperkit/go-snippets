package main

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"os"
)

func main() {
	out := "darnit"
	if len(os.Args) == 2 {
		out = os.Args[1]
	}

	fmt.Printf("gosh %v\n", out)
}
