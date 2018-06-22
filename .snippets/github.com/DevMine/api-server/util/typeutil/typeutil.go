// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package typeutil provide useful type conversion functions.
// All numbers are considered base 10 int64 or float64.
package typeutil

import (
	"strconv"
)

// StrToUint converts 's' into its uint64 form. 's' is expected to contain
// digits only.
func StrToUint(s string) (n uint64, err error) {
	return strconv.ParseUint(s, 10, 64)
}

// StrToInt converts 's' into its int64 form. 's' is expected to contain
// digits only.
func StrToInt(s string) (n int64, err error) {
	return strconv.ParseInt(s, 10, 64)
}

// StrToFloat converts 's' into its float64 form. 's' is expected to be a float
// in its string representation.
func StrToFloat(s string) (f float64, err error) {
	return strconv.ParseFloat(s, 64)
}

// IntToStr converts 'n' base 10 number to its string representation.
func IntToStr(n int64) string {
	return strconv.FormatInt(n, 10)
}
