package cove

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"os"
	"strings"
)

//The Go PATH variable can be a ":" separated list.
//The default behavior for "most" go tools when needing a canonical
//go path location is to just use the first one in the list.
//this function returns that.
func GetFirstGoPath() string {
	return strings.Split(os.Getenv("GOPATH"), ":")[0]
}
