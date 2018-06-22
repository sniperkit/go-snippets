// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"testing"
)

var p1 = &Project{
	Name: "foobar",
	Repo: nil,
	Langs: []*Language{
		&Language{
			Lang:      Go,
			Paradigms: []string{Compiled, Concurrent, Imperative, Structured},
		},
	},
	Packages: []*Package{
		&Package{
			Doc:  []string{"foo"},
			Name: "foo",
			Path: "foo",
			SrcFiles: []*SrcFile{
				&SrcFile{Path: "foo/foo.go"},
			},
			LoC: 42,
		},
	},
	LoC: 42,
}

var p2 = &Project{
	Name: "foobar",
	Repo: nil,
	Langs: []*Language{
		&Language{
			Lang:      Java,
			Paradigms: []string{ObjectOriented, Structured, Imperative, Generic, Reflective, Concurrent},
		},
	},
	Packages: []*Package{
		&Package{
			Doc:  []string{"bar"},
			Name: "bar",
			Path: "bar",
			SrcFiles: []*SrcFile{
				&SrcFile{Path: "bar/bar.java"},
			},
			LoC: 42,
		},
	},
	LoC: 42,
}

var p1p2 = &Project{
	Name: "foobar",
	Repo: nil,
	Langs: []*Language{
		&Language{
			Lang:      Go,
			Paradigms: []string{Compiled, Concurrent, Imperative, Structured},
		},
		&Language{
			Lang:      Java,
			Paradigms: []string{ObjectOriented, Structured, Imperative, Generic, Reflective, Concurrent},
		},
	},
	Packages: []*Package{
		&Package{
			Doc:  []string{"foo"},
			Name: "foo",
			Path: "foo",
			SrcFiles: []*SrcFile{
				&SrcFile{Path: "foo/foo.go"},
			},
			LoC: 42,
		},
		&Package{
			Doc:  []string{"bar"},
			Name: "bar",
			Path: "bar",
			SrcFiles: []*SrcFile{
				&SrcFile{Path: "bar/bar.java"},
			},
			LoC: 42,
		},
	},
	LoC: 84,
}

func TestMergeAll(t *testing.T) {
	// Error: no project supplied as argument.
	expectedErr := errors.New("MergeAll: at least one project must be supplied")
	if _, err := mergeAll(); err == nil {
		t.Errorf("mergeAll: no error found, expected \"%v\"", expectedErr)
	} else if err.Error() != expectedErr.Error() {
		t.Errorf("mergeAll: found error \"%v\", expected \"%v\"", err, expectedErr)
	}

	// Only 1 project to merge should return the very same project.
	mergedP1, err := mergeAll(p1)
	if err != nil {
		t.Fatal(err)
	}
	if !equals(mergedP1, p1) {
		t.Errorf("mergeAll:")
		fmt.Println("found:")
		prettyPrint(mergedP1)
		fmt.Println("\nexpected:")
		prettyPrint(p1)
	}

	// 2 Projects.
	mergedP1P2, err := mergeAll(p1, p2)
	if err != nil {
		t.Fatal(err)
	}
	if !equals(mergedP1P2, p1p2) {
		t.Errorf("mergeAll:")
		fmt.Println("found:")
		prettyPrint(mergedP1P2)
		fmt.Println("\nexpected:")
		prettyPrint(p1p2)
	}
}

func equals(p1, p2 *Project) bool {
	if p1 == nil || p2 == nil {
		return false
	}

	if p1.Name != p2.Name {
		return false
	}

	if len(p1.Langs) != len(p2.Langs) {
		return false
	}
	for i := 0; i < len(p1.Langs); i++ {
		lang1 := p1.Langs[i]
		lang2 := p2.Langs[i]

		if lang1.Lang != lang2.Lang {
			return false
		}

		if len(lang1.Paradigms) != len(lang2.Paradigms) {
			return false
		}
		for j := 0; j < len(lang1.Paradigms); j++ {
			if lang1.Paradigms[j] != lang2.Paradigms[j] {
				return false
			}
		}
	}

	if len(p1.Packages) != len(p2.Packages) {
		return false
	}
	for i := 0; i < len(p1.Packages); i++ {
		pkg1 := p1.Packages[i]
		pkg2 := p2.Packages[i]

		if pkg1.Name != pkg2.Name {
			return false
		}
		if pkg1.LoC != pkg2.LoC {
			return false
		}

		if len(pkg1.Doc) != len(pkg2.Doc) {
			return false
		}
		for j := 0; j < len(pkg1.Doc); j++ {
			if pkg1.Doc[j] != pkg2.Doc[j] {
				return false
			}
		}

		if len(pkg1.SrcFiles) != len(pkg2.SrcFiles) {
			return false
		}
		for j := 0; j < len(pkg1.SrcFiles); j++ {
			sf1 := pkg1.SrcFiles[j]
			sf2 := pkg2.SrcFiles[j]
			if sf1.Path != sf2.Path {
				return false
			}
		}
	}

	if p1.LoC != p2.LoC {
		return false
	}

	return true
}

func prettyPrint(p *Project) {
	buf := new(bytes.Buffer)

	buf.WriteString("Project:\n")
	buf.WriteString("   Name: " + p.Name + "\n")

	buf.WriteString("   Langs:\n")
	for _, lang := range p.Langs {
		buf.WriteString("      Lang: " + lang.Lang + "\n")
		buf.WriteString("      Paradigms: " + fmt.Sprint(lang.Paradigms) + "\n")
		buf.WriteString("      ===\n")
	}

	buf.WriteString("   Packages:\n")
	for _, pkg := range p.Packages {
		buf.WriteString("      Doc: " + fmt.Sprint(pkg.Doc) + "\n")
		buf.WriteString("      Name: " + pkg.Name + "\n")
		buf.WriteString("      Path: " + pkg.Path + "\n")

		buf.WriteString("      SrcFiles:\n")
		for _, srcFile := range pkg.SrcFiles {
			buf.WriteString("         Path: " + srcFile.Path + "\n")
			buf.WriteString("         ===\n")
		}

		buf.WriteString("      LoC: : " + strconv.FormatInt(pkg.LoC, 10) + "\n")
		buf.WriteString("      ===\n")
	}
	buf.WriteString("   LoC: : " + strconv.FormatInt(p.LoC, 10) + "\n")

	fmt.Println(buf.String())
}
