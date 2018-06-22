// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"bytes"
	"errors"
	"testing"

	"github.com/DevMine/srcanlzr/src/ast"
	"github.com/DevMine/srcanlzr/src/token"
)

func TestDecodeStringsList(t *testing.T) {
	expected := []string{"foo", "bar"}
	buf := bytes.NewBufferString(`["foo","bar"]`)
	dec := newDecoder(buf)
	sl := dec.decodeStrings()
	if dec.err != nil {
		t.Fatal(dec.err)
	}
	if !stringsSlicesEquals(sl, expected) {
		t.Errorf("decodeStringsList: found %v, expected %v", sl, expected)
	}
}

func TestDecodeLanguages(t *testing.T) {
	expected := []*Language{
		&Language{
			Lang:      Go,
			Paradigms: []string{Compiled, Concurrent, Imperative, Structured},
		},
		&Language{
			Lang:      Java,
			Paradigms: []string{ObjectOriented, Structured, Imperative, Generic, Reflective, Concurrent},
		},
	}
	buf := bytes.NewBufferString(`[
		{
			"language": "go",
			"paradigms": [
				"compiled",
				"concurrent",
				"imperative",
				"structured"
			]
		},
		{
			"language": "java",
			"paradigms": [
				"object oriented",
				"structured",
				"imperative",
				"generic",
				"reflective",
				"concurrent"
			]
		}
	]`)
	dec := newDecoder(buf)
	langs := dec.decodeLanguages()
	if dec.err != nil {
		t.Fatal(dec.err)
	}
	if l := len(langs); l != 2 {
		t.Fatalf("decodeLanguages: found %d languages, expected 2", l)
	}
	compareLanguages(t, langs[0], expected[0])
	compareLanguages(t, langs[1], expected[1])
}

func TestDecodeLanguage(t *testing.T) {
	expected := &Language{
		Lang:      Go,
		Paradigms: []string{Compiled, Concurrent, Imperative, Structured},
	}
	buf := bytes.NewBufferString(`{"language": "go", "paradigms": ["compiled", "concurrent", "imperative", "structured"]}`)
	dec := newDecoder(buf)
	lang := dec.decodeLanguage()
	if dec.err != nil {
		t.Fatal(dec.err)
	}
	compareLanguages(t, lang, expected)
}

func TestDecodeIdent(t *testing.T) {
	expected := &ast.Ident{
		ExprName: token.IdentName,
		Name:     "foo",
	}
	buf := bytes.NewBufferString(`{"expression_name": "IDENT", "name": "foo"}`)
	dec := newDecoder(buf)
	ident := dec.decodeIdent()
	if dec.err != nil {
		t.Fatal(dec.err)
	}
	if ident.ExprName != expected.ExprName {
		t.Errorf("decodeIdent: found '%s', expected '%s'", ident.ExprName, expected.ExprName)
	}
	if ident.Name != expected.Name {
		t.Errorf("decodeIdent: found '%s', expected '%s'", ident.Name, expected.Name)
	}
}

func TestExtractFirstKey(t *testing.T) {
	// test success
	buf := bytes.NewBufferString(`"expression_name": "IDENT"`)
	dec := newDecoder(buf)
	if val := dec.extractFirstKey("expression_name"); val != "IDENT" {
		t.Errorf("extractFirstKey 'expression_name': found '%s', expected 'IDENT'", val)
	}
	if dec.err != nil {
		t.Error(dec.err)
	}

	// test failure

	buf = bytes.NewBufferString(``)
	dec = newDecoder(buf)
	expectedErr := errors.New("unexpected EOF")
	_ = dec.extractFirstKey("expression_name")
	if dec.err.Error() != expectedErr.Error() {
		t.Errorf("extractFirstKey 'expression_name': found error \"%v\", expected error \"%v\"",
			dec.err, expectedErr)
	}

	buf = bytes.NewBufferString(`"foo": "bar"`)
	dec = newDecoder(buf)
	expectedErr = errors.New("expected key to be 'expression_name', found 'foo'")
	_ = dec.extractFirstKey("expression_name")
	if dec.err.Error() != expectedErr.Error() {
		t.Errorf("extractFirstKey 'expression_name': found error \"%v\", expected error \"%v\"",
			dec.err, expectedErr)
	}

	buf = bytes.NewBufferString(`"expression_name": 42}`)
	dec = newDecoder(buf)
	expectedErr = errors.New("expected 'string literal', found 'integer literal'")
	_ = dec.extractFirstKey("expression_name")
	if dec.err.Error() != expectedErr.Error() {
		t.Errorf("extractFirstKey 'expression_name': found error \"%v\", expected error \"%v\"",
			dec.err, expectedErr)
	}
}

func TestUnmarshalString(t *testing.T) {
	dec := newDecoder(new(bytes.Buffer))
	if str, err := dec.unmarshalString([]byte("foo")); err != nil {
		t.Error(err)
	} else if str != "foo" {
		t.Errorf("unmarshalString 'foo': found '%s', expected 'foo'", str)
	}

	expectedErr := errors.New("unable to unmarshal string: data is nil")
	if _, err := dec.unmarshalString(nil); err == nil {
		t.Errorf("unmarshalString nil: found no error, expected \"%v\"", expectedErr)
	} else if err.Error() != expectedErr.Error() {
		t.Errorf("unmarshalString nil: found error \"%v\", expected \"%v\"", err, expectedErr)
	}
}

func TestUnmarshalInt64(t *testing.T) {
	dec := newDecoder(new(bytes.Buffer))
	if num, err := dec.unmarshalInt64([]byte("42")); err != nil {
		t.Error(err)
	} else if num != 42 {
		t.Errorf("unmarshalInt64 '42': found %d, expected 42", num)
	}
}

func TestUnmarshalFloat64(t *testing.T) {
	dec := newDecoder(new(bytes.Buffer))
	if num, err := dec.unmarshalFloat64([]byte("42.42")); err != nil {
		t.Error(err)
	} else if num-42.42 > 0.0001 {
		t.Errorf("unmarshalFloat64 '42.42': found %f, expected 42.42", num)
	}
}

func TestUnmarshalBool(t *testing.T) {
	dec := newDecoder(new(bytes.Buffer))
	if num, err := dec.unmarshalBool([]byte("true")); err != nil {
		t.Error(err)
	} else if !num {
		t.Error("unmarshalBool 'true': found 'false', expected 'true'")
	}

	dec = newDecoder(new(bytes.Buffer))
	if num, err := dec.unmarshalBool([]byte("false")); err != nil {
		t.Error(err)
	} else if num {
		t.Error("unmarshalBool 'false': found 'true', expected 'false'")
	}

	expectedErr := errors.New("unable to unmarshal boolean: value 'none' is not a boolean")
	dec = newDecoder(new(bytes.Buffer))
	if num, err := dec.unmarshalBool([]byte("none")); err == nil {
		t.Errorf("unmarshalBool 'none': found no error, expected \"%v\"", expectedErr)
	} else if err.Error() != expectedErr.Error() {
		t.Errorf("unmarshalBool 'none': found \"%v\", expected \"%v\"", err, expectedErr)
	} else if num {
		t.Error("unmarshalBool 'none': found 'true', expected 'false'")
	}
}

func TestAssertNewObject(t *testing.T) {
	type returnStatus struct {
		status bool
		err    error
	}
	input := map[string]returnStatus{
		"{": returnStatus{true, nil},
		"":  returnStatus{false, errors.New("expected value, found EOF")},
		"[": returnStatus{false, errors.New("expected object, found '['")},
	}

	for in, ret := range input {
		buf := bytes.NewBufferString(in)
		dec := newDecoder(buf)
		if status := dec.assertNewObject(); status != ret.status {
			t.Errorf("assertNewObject '%s': found %v, expected %v", in, status, ret.status)
		}
		if err := dec.err; err != nil && err.Error() != ret.err.Error() {
			t.Errorf("assertNewObject '%s': found error \"%v\", expected error \"%v\"", in, dec.err, ret.err)
		}
	}
}

func TestAssertNewArray(t *testing.T) {
	type returnStatus struct {
		status bool
		err    error
	}
	input := map[string]returnStatus{
		"[": returnStatus{true, nil},
		"":  returnStatus{false, errors.New("expected value, found EOF")},
		"{": returnStatus{false, errors.New("expected array, found '{'")},
	}

	for in, ret := range input {
		buf := bytes.NewBufferString(in)
		dec := newDecoder(buf)
		if status := dec.assertNewArray(); status != ret.status {
			t.Errorf("assertNewArray '%s': found %v, expected %v", in, status, ret.status)
		}
		if err := dec.err; err != nil && err.Error() != ret.err.Error() {
			t.Errorf("assertNewArray '%s': found error \"%v\", expected error \"%v\"", in, dec.err, ret.err)
		}
	}
}

func TestIsEndObject(t *testing.T) {
	type returnStatus struct {
		status bool
		err    error
	}
	input := map[string]returnStatus{
		"}":   returnStatus{true, nil},
		",":   returnStatus{false, nil},
		"42}": returnStatus{false, errors.New("expected 'comma', found 'integer literal'")},
	}

	for in, ret := range input {
		buf := bytes.NewBufferString(in)
		dec := newDecoder(buf)
		if status := dec.isEndObject(); status != ret.status {
			t.Errorf("isEndObject '%s': found %v, expected %v", in, status, ret.status)
		}
		if err := dec.err; err != nil && err.Error() != ret.err.Error() {
			t.Errorf("isEndObject '%s': found error \"%v\", expected error \"%v\"", in, dec.err, ret.err)
		}
	}
}

func TestIsEndArray(t *testing.T) {
	type returnStatus struct {
		status bool
		err    error
	}
	input := map[string]returnStatus{
		"]":   returnStatus{true, nil},
		",":   returnStatus{false, nil},
		"42}": returnStatus{false, errors.New("expected 'comma', found 'integer literal'")},
	}

	for in, ret := range input {
		buf := bytes.NewBufferString(in)
		dec := newDecoder(buf)
		if status := dec.isEndArray(); status != ret.status {
			t.Errorf("isEndArray '%s': found %v, expected %v", in, status, ret.status)
		}
		if err := dec.err; err != nil && err.Error() != ret.err.Error() {
			t.Errorf("isEndArray '%s': found error \"%v\", expected error \"%v\"", in, dec.err, ret.err)
		}
	}
}

func TestIsEmptyObject(t *testing.T) {
	type returnStatus struct {
		status bool
		err    error
	}
	input := map[string]returnStatus{
		"}":        returnStatus{true, nil},
		"":         returnStatus{false, errors.New("unexpected EOF")},
		"\"foo\":": returnStatus{false, nil},
	}

	for in, ret := range input {
		buf := bytes.NewBufferString(in)
		dec := newDecoder(buf)
		if status := dec.isEmptyObject(); status != ret.status {
			t.Errorf("isEmptyObject '%s': found %v, expected %v", in, status, ret.status)
		}
		if err := dec.err; err != nil && err.Error() != ret.err.Error() {
			t.Errorf("isEmptyObject '%s': found error \"%v\", expected error \"%v\"", in, dec.err, ret.err)
		}
	}
}

func TestIsEmptyArray(t *testing.T) {
	type returnStatus struct {
		status bool
		err    error
	}
	input := map[string]returnStatus{
		"]":  returnStatus{true, nil},
		"":   returnStatus{false, errors.New("unexpected EOF")},
		"42": returnStatus{false, nil},
	}

	for in, ret := range input {
		buf := bytes.NewBufferString(in)
		dec := newDecoder(buf)
		if status := dec.isEmptyArray(); status != ret.status {
			t.Errorf("isEmptyArray '%s': found %v, expected %v", in, status, ret.status)
		}
		if err := dec.err; err != nil && err.Error() != ret.err.Error() {
			t.Errorf("isEmptyArray '%s': found error \"%v\", expected error \"%v\"", in, dec.err, ret.err)
		}
	}
}

func TestIsNull(t *testing.T) {
	input := map[string]bool{
		"null":  true,
		"nan":   false,
		`"foo"`: false,
	}

	for in, expected := range input {
		buf := bytes.NewBufferString(in)
		dec := newDecoder(buf)
		if b := dec.isNull(); b != expected {
			t.Errorf("isNull '%s': found '%v', expected '%v'", in, b, expected)
		}
	}
}

func compareLanguages(t *testing.T, found, expected *Language) {
	if found == nil {
		t.Fatal("decodeLanguage: should not be nil")
	}

	if expected.Lang != found.Lang {
		t.Errorf("decodeLanguage.Lang: found '%s', expected '%s'", found.Lang, expected.Lang)
	}
	if !stringsSlicesEquals(expected.Paradigms, found.Paradigms) {
		t.Errorf("decodeLanguage.Paradigms: found '%v', expected '%v'", found.Paradigms, expected.Paradigms)
	}
}

func stringsSlicesEquals(sl1, sl2 []string) bool {
	if len(sl1) != len(sl2) {
		return false
	}
	for i := 0; i < len(sl1); i++ {
		if sl1[i] != sl2[i] {
			return false
		}
	}
	return true
}
