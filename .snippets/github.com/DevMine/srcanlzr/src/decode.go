// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run ./gen/gen_ast_decoder.go

package src

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/DevMine/repotool/model"
)

type decoder struct {
	scan *scanner
	buf  []byte
	err  error
}

// newDecoder creates a new JSON decoder that reads from r.
func newDecoder(r io.Reader) *decoder {
	return &decoder{scan: newScanner(r)}
}

// decode decodes JSON input into a src.Project structure.
func (dec *decoder) decode() (*Project, error) {
	prj := dec.decodeProject()
	if dec.err != nil {
		return nil, dec.errorf(dec.err)
	}
	return prj, nil
}

func (dec *decoder) errorf(err error) error {
	return fmt.Errorf("malformed json at %d: %v", dec.scan.globPos, err)
}

// decodeProject decodes a project object.
func (dec *decoder) decodeProject() *Project {
	if !dec.assertNewObject() {
		return nil
	}

	prj := Project{}

	if dec.isEmptyObject() {
		return &prj
	}
	if dec.err != nil {
		return nil
	}

	for {
		key, err := dec.scan.nextKey()
		if err != nil {
			if err == io.EOF {
				break
			}
			dec.err = err
			return nil
		}
		if key == "" {
			dec.err = errors.New("empty key")
			return nil
		}

		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}

		switch key {
		case "packages":
			dec.scan.back()
			prj.Packages = dec.decodePackages()
		case "languages":
			dec.scan.back()
			prj.Langs = dec.decodeLanguages()
		case "repository":
			dec.scan.back()
			prj.Repo = dec.decodeRepository()
		case "loc":
			if tok != scanInt64Lit {
				dec.err = fmt.Errorf("expected integer literal, found %v", tok)
				return nil
			}
			prj.LoC, dec.err = dec.unmarshalInt64(val)
		case "name":
			if tok != scanStringLit {
				dec.err = fmt.Errorf("expected string literal, found %v", tok)
				return nil
			}
			prj.Name, dec.err = dec.unmarshalString(val)
		default:
			dec.err = fmt.Errorf("unexpected key '%s' for project object", key)
		}

		if dec.err != nil {
			return nil
		}

		if dec.isEndObject() {
			break
		}
		if err != nil {
			return nil
		}
	}
	return &prj
}

// decodePackages decodes a list of package objects.
func (dec *decoder) decodePackages() []*Package {
	if !dec.assertNewArray() {
		return nil
	}

	pkgs := []*Package{}

	if dec.isEmptyArray() {
		return pkgs
	}
	if dec.err != nil {
		return nil
	}

	for {
		pkg := dec.decodePackage()
		if dec.err != nil {
			return nil
		}
		pkgs = append(pkgs, pkg)

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}

	return pkgs
}

// decoderPackage decodes a package object.
func (dec *decoder) decodePackage() *Package {
	if !dec.assertNewObject() {
		return nil
	}

	pkg := Package{}

	if dec.isEmptyObject() {
		return &pkg
	}
	if dec.err != nil {
		return nil
	}

	for {
		var key string
		key, dec.err = dec.scan.nextKey()
		if dec.err != nil {
			return nil
		}

		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}

		switch key {
		case "source_files":
			dec.scan.back()
			pkg.SrcFiles = dec.decodeSrcFiles()
		case "doc":
			dec.scan.back()
			pkg.Doc = dec.decodeStrings()
		case "loc":
			if tok != scanInt64Lit {
				dec.err = fmt.Errorf("expected integer literal, found %v", tok)
				return nil
			}
			pkg.LoC, dec.err = dec.unmarshalInt64(val)
		case "name":
			if tok != scanStringLit {
				dec.err = fmt.Errorf("expected string literal, found %v", tok)
				return nil
			}
			pkg.Name, dec.err = dec.unmarshalString(val)
		case "path":
			if tok != scanStringLit {
				dec.err = fmt.Errorf("expected string literal, found %v", tok)
				return nil
			}
			pkg.Path, dec.err = dec.unmarshalString(val)
		default:
			dec.err = fmt.Errorf("unexpected key '%s' for package object", key)
		}

		if dec.err != nil {
			return nil
		}

		if dec.isEndObject() {
			break
		}
		if err != nil {
			return nil
		}
	}

	return &pkg
}

// decodeSrcFiles decodes a list of source file objects.
func (dec *decoder) decodeSrcFiles() []*SrcFile {
	if !dec.assertNewArray() {
		return nil
	}

	sf := []*SrcFile{}

	if dec.isEmptyArray() {
		return sf
	}
	if dec.err != nil {
		return nil
	}

	for {
		srcFile := dec.decodeSrcFile()
		if dec.err != nil {
			return nil
		}

		sf = append(sf, srcFile)

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}

	return sf
}

func (dec *decoder) decodeSrcFile() *SrcFile {
	if !dec.assertNewObject() {
		return nil
	}

	sf := SrcFile{}

	if dec.isEmptyObject() {
		return &sf
	}
	if dec.err != nil {
		return nil
	}

	for {
		key, err := dec.scan.nextKey()
		if err != nil {
			if err == io.EOF {
				break
			}
			dec.err = err
			return nil
		}
		if key == "" {
			dec.err = errors.New("empty key")
			return nil
		}

		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}

		switch key {
		case "path":
			if tok != scanStringLit {
				dec.err = fmt.Errorf("expected 'string literal', found '%v'", tok)
			}
			sf.Path, dec.err = dec.unmarshalString(val)
		case "language":
			dec.scan.back()
			sf.Lang = dec.decodeLanguage()
		case "imports":
			dec.scan.back()
			sf.Imports = dec.decodeStrings()
		case "type_specifiers":
			dec.scan.back()
			sf.TypeSpecs = dec.decodeTypeSpecs()
		case "structs":
			dec.scan.back()
			sf.Structs = dec.decodeStructTypes()
		case "constants":
			dec.scan.back()
			sf.Constants = dec.decodeGlobalDecls()
		case "variables":
			dec.scan.back()
			sf.Vars = dec.decodeGlobalDecls()
		case "functions":
			dec.scan.back()
			sf.Funcs = dec.decodeFuncDecls()
		case "interfaces":
			dec.scan.back()
			sf.Interfaces = dec.decodeInterfaces()
		case "classes":
			dec.scan.back()
			sf.Classes = dec.decodeClassDecls()
		case "enums":
			dec.scan.back()
			sf.Enums = dec.decodeEnumDecls()
		case "traits":
			dec.scan.back()
			sf.Traits = dec.decodeTraits()
		case "loc":
			if tok != scanInt64Lit {
				dec.err = fmt.Errorf("expected integer literal, found %v", tok)
				return nil
			}
			sf.LoC, dec.err = dec.unmarshalInt64(val)
		default:
			dec.err = fmt.Errorf("unexpected key '%s' for source file object", key)
		}

		if dec.err != nil {
			return nil
		}

		if dec.isEndObject() {
			break
		}
		if err != nil {
			return nil
		}

	}

	return &sf
}

// decoderStrings decodes a list of strings.
func (dec *decoder) decodeStrings() []string {
	if dec.isNull() {
		return nil
	}
	if !dec.assertNewArray() {
		return nil
	}

	sl := []string{}

	if dec.isEmptyArray() {
		return sl
	}
	if dec.err != nil {
		return nil
	}

	for {
		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}
		if tok != scanStringLit {
			dec.err = fmt.Errorf("expected string, found %v", tok)
			return nil
		}

		str, err := dec.unmarshalString(val)
		if err != nil {
			dec.err = err
			return nil
		}
		sl = append(sl, str)

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}
	return sl
}

// decoderInt64s decodes a list of int 64.
func (dec *decoder) decodeInt64s() []int64 {
	if !dec.assertNewArray() {
		return nil
	}

	il := []int64{}

	if dec.isEmptyArray() {
		return il
	}
	if dec.err != nil {
		return nil
	}

	for {
		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}
		if tok != scanInt64Lit {
			dec.err = fmt.Errorf("expected integer, found %v", tok)
			return nil
		}
		num, err := dec.unmarshalInt64(val)
		if err != nil {
			dec.err = err
			return nil
		}
		il = append(il, num)

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}
	return il
}

// decodeLanguages decodes a list of languages.
func (dec *decoder) decodeLanguages() []*Language {
	if !dec.assertNewArray() {
		return nil
	}

	ls := []*Language{}

	if dec.isEmptyArray() {
		return ls
	}
	if dec.err != nil {
		return nil
	}

	for {
		lang := dec.decodeLanguage()
		if dec.err != nil {
			return nil
		}

		ls = append(ls, lang)

		if dec.isEndArray() {
			break
		}
		if dec.err != nil {
			return nil
		}
	}

	return ls
}

// decodeLanguage decode a src.Language object.
func (dec *decoder) decodeLanguage() *Language {
	if !dec.assertNewObject() {
		return nil
	}

	lang := Language{}

	if dec.isEmptyObject() {
		return &lang
	}
	if dec.err != nil {
		return nil
	}

	for {
		key, err := dec.scan.nextKey()
		if err != nil {
			if err == io.EOF {
				break
			}
			dec.err = err
			return nil
		}
		if key == "" {
			dec.err = errors.New("empty key")
			return nil
		}

		val, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return nil
		}

		switch key {
		case "paradigms":
			// Since the '[' character has been consumed, we need to step back
			// brefore calling decodeStringsList.
			dec.scan.back()
			lang.Paradigms = dec.decodeStrings()
		case "language":
			if tok != scanStringLit {
				dec.err = fmt.Errorf("expected 'string literal', found '%v'", tok)
				return nil
			}
			lang.Lang, dec.err = dec.unmarshalString(val)
		default:
			dec.err = fmt.Errorf("unexpected key '%s' for language object", key)
		}

		if dec.err != nil {
			return nil
		}

		if dec.isEndObject() {
			break
		}
		if err != nil {
			return nil
		}
	}
	return &lang
}

// TODO: implement
func (dec *decoder) decodeRepository() *model.Repository {
	return nil
}

// extractExprName looks for the first key of an object, which must be
// match the given key, and returns its value. The '{' character must have been
// previously consumed and value corresponding to the key must be a string.
//
// Errors are reported in dec.err and value corresponding to the key must be a string.
func (dec *decoder) extractFirstKey(key string) string {
	k, err := dec.scan.nextKey()
	if err != nil {
		if err == io.EOF {
			dec.err = errors.New("unexpected EOF")
			return ""
		}
		dec.err = err
		return ""
	}
	if k != key {
		dec.err = fmt.Errorf("expected key to be '%s', found '%s'", key, k)
		return ""
	}

	val, tok, err := dec.scan.nextValue()
	if err != nil {
		dec.err = err
		return ""
	}
	if tok != scanStringLit {
		dec.err = fmt.Errorf("expected 'string literal', found '%v'", tok)
		return ""
	}
	return string(val)
}

// unmarshalInt unmarshal integer value into an int 64.
func (dec *decoder) unmarshalInt64(data []byte) (int64, error) {
	num, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// unmarshalFloat unmarshal floating point number into a float 64.
func (dec *decoder) unmarshalFloat64(data []byte) (float64, error) {
	num, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0.0, err
	}
	return num, nil
}

// unmarshalString unmarshals a bytes slice into a string.
func (dec *decoder) unmarshalString(data []byte) (string, error) {
	if data == nil {
		return "", errors.New("unable to unmarshal string: data is nil")
	}
	return string(data), nil
}

// unmarshalBool unmarshals a bytes slice into a boolean.
func (dec *decoder) unmarshalBool(data []byte) (bool, error) {
	if data == nil {
		return false, errors.New("unable to unmarshal boolean: data is nil")
	}
	switch str := string(data); str {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("unable to unmarshal boolean: value '%s' is not a boolean", str)
	}
}

// assertNewObject makes sure that the next value is a new object. In other
// words, the next value must begin with a '{'. If it is not, it will set
// dec.err and return false.
func (dec *decoder) assertNewObject() bool {
	// Since Language is a JSON Object, we expect to find a '{' character.
	_, tok, err := dec.scan.nextValue()
	if err != nil {
		dec.err = err
		return false
	}
	if tok != scanBeginObject {
		dec.err = fmt.Errorf("expected object, found '%v'", tok)
		return false
	}
	return true
}

// assertNewArray makes sure that the next value is a new array. In order
// words, the next value must begin with a '['. If it is not, it will set
// dec.err and return false.
func (dec *decoder) assertNewArray() bool {
	_, tok, err := dec.scan.nextValue()
	if err != nil {
		dec.err = err
		return false
	}
	if tok != scanBeginArray {
		dec.err = fmt.Errorf("expected array, found '%v'", tok)
		return false
	}
	return true
}

// isEndObject returns true if the next value marks the end of the object
// ('}') and false otherwise. If it is false, the next value must be a
// comma. If not, it will set dec.err accordingly.
func (dec *decoder) isEndObject() bool {
	_, tok, err := dec.scan.nextValue()
	if err != nil {
		dec.err = err
		return false
	}
	if tok == scanEndObject {
		return true
	}
	if tok != scanComma {
		dec.err = fmt.Errorf("expected 'comma', found '%v'", tok)
	}
	return false
}

// isEndArray returns true if the next value marks the end of the array
// (']') and false otherwise. If it is false, the next value must be a
// comma. If not, it will set dec.err accordingly.
func (dec *decoder) isEndArray() bool {
	_, tok, err := dec.scan.nextValue()
	if err != nil {
		dec.err = err
		return false
	}
	if tok == scanEndArray {
		return true
	}
	if tok != scanComma {
		dec.err = fmt.Errorf("expected 'comma', found '%s'", tok)
	}
	return false
}

// isEmptyObject tests if the object is empty (no key/value pairs inside).
//
// This method does not consume any byte except when the object is empty.
// In this case, it consumes the '}'.
//
// If an error occurs, it returns false and set dec.err.
func (dec *decoder) isEmptyObject() bool {
	// The object can be empty, so we have to check for that and without
	// consuming the next byte.
	if b, err := dec.scan.peek(); err != nil {
		if err == io.EOF {
			dec.err = errors.New("unexpected EOF")
		} else {
			dec.err = err
		}
		return false
	} else if b == '}' {
		// We need to read the next byte here because if the caller accept empty
		// object, it will continue the decoding and won't expect to find a '}'.
		_, dec.err = dec.scan.read()
		return true
	}
	return false
}

// isEmptyArray tests if the object is empty (no values inside).
//
// This method does not consume any byte except when the array is empty.
// In this case, it consumes the ']'.
//
// If an error occurs, it returns false and set dec.err.
func (dec *decoder) isEmptyArray() bool {
	if b, err := dec.scan.peek(); err != nil {
		if err == io.EOF {
			dec.err = errors.New("unexpected EOF")
		} else {
			dec.err = err
		}
		return false
	} else if b == ']' {
		// We need to read the next byte here because if the caller accept empty
		// array, it will continue the decoding and won't expect to find a ']'.
		_, dec.err = dec.scan.read()
		return true
	}
	return false
}

// isNull tests if the next value is null.
//
// If the next value is null, it consumes it.
//
// It returns true if the next value is null, false otherwise.
// If an error occurs, it returns false and set dec.err.
func (dec *decoder) isNull() bool {
	if b, err := dec.scan.peek(); err != nil {
		if err == io.EOF {
			dec.err = errors.New("unexpected EOF")
		} else {
			dec.err = err
		}
		return false
	} else if b == 'n' {
		_, tok, err := dec.scan.nextValue()
		if err != nil {
			dec.err = err
			return false
		}
		if tok != scanNullVal {
			if dec.err != nil {
				dec.err = fmt.Errorf("expected 'null', found '%v'", tok)
				return false
			}
		}
		return true
	}
	return false
}
