// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"errors"
	"fmt"
)

func mergeAll(ps ...*Project) (*Project, error) {
	if len(ps) == 0 {
		return nil, errors.New("MergeAll: at least one project must be supplied")
	}

	newPrj := &Project{
		Name:     ps[0].Name,
		Langs:    ps[0].Langs,
		Packages: ps[0].Packages,
		LoC:      ps[0].LoC,
	}

	if len(ps) == 1 {
		return newPrj, nil
	}

	for i := 1; i < len(ps); i++ {
		curr := ps[i]
		if curr == nil {
			return nil, fmt.Errorf("MergeAll: ps[%d] is nil", i)
		}
		merge(newPrj, curr)
	}

	return newPrj, nil
}

// merge p2 into p1
func merge(p1, p2 *Project) {
	if p1.Name == "" {
		p1.Name = p2.Name
	}
	p1.Langs = append(p1.Langs, p2.Langs...)
	p1.LoC += p2.LoC
	p1.Packages = mergePackages(p1.Packages, p2.Packages)
}

// mergePackages merges pkg2 into pkg1
func mergePackages(pkgs1, pkgs2 []*Package) []*Package {
	for _, pkg2 := range pkgs2 {
		pkg1 := findPackage(pkgs1, pkg2.Path)
		if pkg1 == nil {
			pkgs1 = append(pkgs1, pkg2)
		} else {
			pkg1.SrcFiles = append(pkg1.SrcFiles, pkg2.SrcFiles...)
			pkg1.Doc = append(pkg1.Doc, pkg2.Doc...)
			pkg1.LoC += pkg2.LoC
		}
	}
	return pkgs1
}

// findPackage find a package by path in pkgs and return its reference.
func findPackage(pkgs []*Package, path string) *Package {
	for _, pkg := range pkgs {
		if pkg.Path == path {
			return pkg
		}
	}
	return nil
}
