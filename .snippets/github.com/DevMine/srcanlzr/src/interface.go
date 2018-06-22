// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"io"
	"os"
)

// Decode decodes a JSON encoded src.Project read from r.
func Decode(r io.Reader) (*Project, error) {
	dec := newDecoder(r)
	return dec.decode()
}

// DecodeFile decodes a JSON encoded src.Project read from a given file.
func DecodeFile(path string) (*Project, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
}

// MergeAll merges a list of projects.
//
// There must be at least one project. In this case, it just returns a copy of
// the project. Moreover, the projects must be distinct.
//
// The merge only performs shallow copies, which means that if the field value
// is a pointer it copies the memory address and not the value pointed.
func MergeAll(ps ...*Project) (*Project, error) {
	return mergeAll(ps...)
}

// Merge merges two project. See MergeAll for more details.
func Merge(p1, p2 *Project) *Project {
	// merge() merges p2 into p1, therefore we need to copy p1 before merging.
	newPrj := &Project{
		Name:     p1.Name,
		Langs:    p1.Langs,
		Packages: p1.Packages,
		LoC:      p1.LoC,
	}
	merge(newPrj, p2)
	return newPrj
}
