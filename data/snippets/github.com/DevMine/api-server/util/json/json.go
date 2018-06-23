// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package json provides json marshaling functions that do not return any
// error but may panic if something bad occurs.
package json

import "encoding/json"

// MarshalPanic marshalizes the interface v into JSON.
// If an error occurs in the process, it panics.
func MarshalPanic(v interface{}) []byte {
	bs, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return bs
}

// MarshalIndentPanic marshaliues the interface v into JSON and indents it.
// If an error occurs in the process, it panics.
func MarshalIndentPanic(v interface{}) []byte {
	bs, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		panic(err)
	}
	return bs
}
