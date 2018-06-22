// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"encoding/json"
	"io"
	"os"

	"github.com/DevMine/repotool/model"
	"github.com/DevMine/srcanlzr/src/ast"
)

// Supported VCS (Version Control System)
const (
	Git = "git"
	Hg  = "mercurial"
	SVN = "subversion"
	Bzr = "bazaar"
	CVS = "cvs"
)

var suppVCS = []string{
	Git,
	Hg,
	SVN,
	Bzr,
	CVS,
}

// Supported programming languages
const (
	Go     = "go"
	Ruby   = "ruby"
	Python = "python"
	C      = "c"
	Java   = "java"
	Scala  = "scala"
)

var suppLang = []string{
	Go,
	Ruby,
	Python,
	C,
	Java,
	Scala,
}

// Supported paradigms
const (
	Structured     = "structured"
	Imperative     = "imperative"
	Procedural     = "procedural"
	Compiled       = "compiled"
	Concurrent     = "concurrent"
	Functional     = "functional"
	ObjectOriented = "object oriented"
	Generic        = "generic"
	Reflective     = "reflective"
)

var suppParadigms = []string{
	Structured,
	Imperative,
	Procedural,
	Compiled,
	Concurrent,
	Functional,
	ObjectOriented,
	Generic,
	Reflective,
}

// A Language represents a programming language.
type Language struct {
	// The programming language name (e.g. go, ruby, java, etc.)
	//
	// The name must match one of the supported programming languages defined in
	// the constants.
	Lang string `json:"language"` // TODO rename into name

	// The paradigms of the programming language (e.g. structured, imperative,
	// object oriented, etc.)
	//
	// The name must match one of the supported paradigms defined in the
	// constants.
	Paradigms []string `json:"paradigms"`
}

// Package holds information about a package, which is, basically, just a
// folder that contains at least one source file.
type Package struct {
	// The package documentation, or nil.
	// TODO support docucmentation for multiple languages.
	Doc []string `json:"doc,omitempty"`

	// The package name. This should be the name of the parent folder.
	Name string `json:"name"`

	// The full path of the package. The path must be relative to the root of
	// the project and never be an absolute path.
	Path string `json:"path"`

	// The list of all source files contained in the package.
	SrcFiles []*SrcFile `json:"source_files"`

	// The total number of lines of code of the package.
	LoC int64 `json:"loc"`
}

// SrcFile holds information about a source file.
type SrcFile struct {
	// The path of the source file, relative to the root of the project.
	Path string `json:"path"`

	// Programming language used.
	Lang *Language `json:"language"`

	// List of the imports used by the srouce file.
	Imports []string `json:"imports,omitempty"`

	// Types definition
	TypeSpecs []*ast.TypeSpec `json:"type_specifiers,omitempty"`

	// Structures definition
	// TODO rename JSON key into structures
	Structs []*ast.StructType `json:"structs,omitempty"`

	// List of constants defined at the file level (e.g. global constants)
	Constants []*ast.GlobalDecl `json:"constants,omitempty"`

	// List of variables defined at the file level (e.g. global variables)
	Vars []*ast.GlobalDecl `json:"variables,omitempty"`

	// List of functions
	Funcs []*ast.FuncDecl `json:"functions,omitempty"`

	// List of interfaces
	Interfaces []*ast.Interface `json:"interfaces,omitempty"`

	// List of classes
	Classes []*ast.ClassDecl `json:"classes,omitempty"`

	// List of enums
	Enums []*ast.EnumDecl `json:"enums,omitempty"`

	// List of traits
	// See http://en.wikipedia.org/wiki/Trait_%28computer_programming%29
	Traits []*ast.Trait `json:"traits,omitempty"`

	// The total number of lines of code.
	LoC int64 `json:"loc"`
}

// Project is the root of the src API and therefore it must be at the root of
// the JSON.
//
// It contains the metadata of a project and the list of all packages.
type Project struct {
	// The name of the project. Since it may be something really difficult to
	// guess, it should generally be the name of the folder containing the
	// project.
	Name string `json:"name"`

	// The repository in which the project is hosted, or nil. This field is not
	// meant to be filled by one of the language parsers. Only repotool should
	// take care of it. For more details, see:
	//    https://github.com/DevMine/repotool
	//
	// Since this field uses an external type, it is not unmarshalled by
	// src.Unmarshal itself but by the standard json.Unmarshal function.
	// To do so, its unmarshalling is defered using json.RawMessage.
	// See the RepoRaw field.
	Repo *model.Repository `json:"repository,omitempty"`

	// The list of all programming languages used by the project. Each language
	// must be added by the corresponding language parsers if and only if the
	// project contains at least one line of code written in this language.
	Langs []*Language `json:"languages"`

	// List of all packages of the project. We call "package" every folder that
	// contains at least one source file.
	Packages []*Package `json:"packages"`

	// The total number of lines of code in the whole project, independently of
	// the language.
	LoC int64 `json:"loc"`
}

// Encode writes JSON representation of the project into w.
//
// For now, encoding still make use of the json package of the standard libary.
func (p *Project) Encode(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(p)
}

// EncodeToFile writes JSON representation of the project into a file located at path.
func (p *Project) EncodeToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return p.Encode(f)
}
