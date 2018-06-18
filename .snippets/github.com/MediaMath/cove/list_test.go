package cove

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"reflect"
	"testing"

	"github.com/MediaMath/cove/cmd"
)

//These tests verify the behavior of go list
func TestSinglePathResult(t *testing.T) {
	cmd, err := list("os/exec")
	if !eq(cmd, sl("os/exec")) {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestMultiPathResult(t *testing.T) {
	cmd, err := list("os/...")
	if !eq(cmd, sl("os", "os/exec", "os/signal", "os/user")) {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResult(t *testing.T) {
	cmd, err := list("os/foodog")
	if len(cmd) > 0 || err == nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithPassingToo(t *testing.T) {
	cmd, err := list("os/foodog", "os/signal")
	if !eq(cmd, sl("os/signal")) || err == nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithErrorsOff(t *testing.T) {
	cmd, err := list("-e", "os/foodog")
	if !eq(cmd, sl("os/foodog")) || err != nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithPassingTooWithErrorsOff(t *testing.T) {
	cmd, err := list("-e", "os/foodog", "os/signal")
	if !eq(cmd, sl("os/foodog", "os/signal")) || err != nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func list(args ...string) ([]string, error) {
	return cmd.Output(GoCmd("list", args...))
}

func eq(slice1 []string, slice2 []string) bool {
	return reflect.DeepEqual(slice1, slice2)
}

func sl(items ...string) []string {
	return items
}
