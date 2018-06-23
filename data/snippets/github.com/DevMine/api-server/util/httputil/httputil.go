// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package httputil provides utilities useful when dealing with http responses.
package httputil

import (
	"encoding/json"
)

// ResponseError shall be used when an error response shall be returned.
type ResponseError struct {
	// Message shall be an informative message about the error.
	Message string `json:"message"`
}

// NewResponseError creates a new ResponseError.
func NewResponseError(msg string) *ResponseError {
	return &ResponseError{Message: msg}
}

// JSON formats a ResponseError to JSON.
func (re *ResponseError) JSON() string {
	// We do not want to create a 500 internal server error (by calling
	// util/json which may panic) when trying to display an error message to
	// the user, hence use encoding/json here.
	bs, err := json.MarshalIndent(re, "", "    ")
	if err != nil {
		return "{}"
	}
	return string(bs)
}
