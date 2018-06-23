// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"bytes"
	"errors"
	"testing"
)

func TestIgnoreWhitespaces(t *testing.T) {
	buf := bytes.NewBufferString("     x")
	scan := newScanner(buf)
	if err := scan.ignoreWhitespaces(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c := scan.buf[scan.pos:]; len(c) != 1 || c[0] != 'x' {
		t.Errorf("scan.ignoreWhitespace: found '%s', expected 'x'", c)
	}
}

func TestIsWhitespace(t *testing.T) {
	ws := []byte{' ', '\n', '\t', '\r'}
	for _, c := range ws {
		if !isWhitespace(c) {
			t.Errorf("'%c' should be considered as a whitespace", c)
		}
	}

	nows := []byte{'a', ',', ';', ':', '0', '\v'}
	for _, c := range nows {
		if isWhitespace(c) {
			t.Errorf("'%c' should not be considered as a whitespace", c)
		}
	}
}

func TestIsDigit(t *testing.T) {
	digits := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	for _, c := range digits {
		if !isDigit(c) {
			t.Errorf("'%c' should be considered as a digit", c)
		}
	}

	var c byte
	for c = 'a'; c <= 'z'; c++ {
		if isDigit(c) {
			t.Errorf("'%c' should not be considered as a digit", c)
		}
	}
}

func TestNextKey(t *testing.T) {
	validKeys := []string{
		`"foo":`,
		`"foo"  :`,
		`"foo": "bar"`,
	}
	for _, keyInput := range validKeys {
		buf := bytes.NewBufferString(keyInput)
		scan := newScanner(buf)
		key, err := scan.nextKey()
		if err != nil {
			t.Error(err)
		}
		if key != "foo" {
			t.Errorf("nextKey: found '%s', expected 'foo'", key)
		}
	}

	invalidKeys := map[string]error{
		`"foo"`: errors.New("expected ':', found EOF"),
		`"foo`:  errors.New("expected key, found EOF"),
		`foo`:   errors.New("expected '\"', found 'f'"),
	}
	for keyInput, expectedErr := range invalidKeys {
		buf := bytes.NewBufferString(keyInput)
		scan := newScanner(buf)
		_, err := scan.nextKey()
		if err == nil {
			t.Errorf("nextKey: '%s' is expected to return the error '%v'", keyInput, expectedErr)
			continue
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("nextKey: found error \"%v\", expected error \"%v\"", err, expectedErr)
		}
	}
}

func TestNextValue(t *testing.T) {
	// test valid values
	validKeys := map[string]scanToken{
		`[`:     scanBeginArray,
		`]`:     scanEndArray,
		`{`:     scanBeginObject,
		`}`:     scanEndObject,
		`42`:    scanInt64Lit,
		`42.42`: scanFloat64Lit,
		`"foo"`: scanStringLit,
		`null`:  scanNullVal,
		`true`:  scanBoolLit,
		`false`: scanBoolLit,
		`,`:     scanComma,
	}
	for val, expectedTok := range validKeys {
		buf := bytes.NewBufferString(val)
		scan := newScanner(buf)
		_, tok, err := scan.nextValue()
		if err != nil {
			t.Errorf("nextValue '%s': %v", string(val), err)
		}
		if expectedTok != tok {
			t.Errorf("nextValue '%s': found '%v', expected '%v'", string(val), tok, expectedTok)
		}
	}

	// test invalid values
	invalidKeys := map[string]error{
		`z`:  errors.New("illegal value 'z'"),
		``:   errors.New("expected value, found EOF"),
		`  `: errors.New("expected value, found EOF"),
	}
	for val, expectedErr := range invalidKeys {
		buf := bytes.NewBufferString(val)
		scan := newScanner(buf)
		_, tok, err := scan.nextValue()
		if err == nil {
			t.Errorf("nextValue '%s': found no error, expected \"%v\"", string(val), expectedErr)
			continue
		}
		if err.Error() != expectedErr.Error() {
			t.Errorf("nextValue '%s': found \"%v\", expected \"%v\"", string(val), err, expectedErr)
		}
		if tok != scanIllegalToken {
			t.Errorf("nextValue '%s': found '%v', expected '%v'", string(val), tok, scanIllegalToken)
		}
	}
}

func TestReadNumber(t *testing.T) {
	// test int64
	buf := bytes.NewBufferString("42}")
	scan := newScanner(buf)
	val, tok, err := scan.readNumber()
	if err != nil {
		t.Fatal("readNumber '42}': ", err)
	}
	if tok != scanInt64Lit {
		t.Errorf("readNumber '42}': found token '%v', expected '%v'", tok, scanInt64Lit)
	}
	if s := string(val); s != "42" {
		t.Errorf("readNumber '42': found '%s', expected '42'", s)
	}

	// test float 64
	buf = bytes.NewBufferString("42.24}")
	scan = newScanner(buf)
	val, tok, err = scan.readNumber()
	if err != nil {
		t.Fatal("readNumber '42.24}': ", err)
	}
	if tok != scanFloat64Lit {
		t.Errorf("readNumber '42.24}': found token '%v', expected '%v'", tok, scanFloat64Lit)
	}
	if s := string(val); s != "42.24" {
		t.Errorf("readNumber '42.24}': found '%s', expected '42.24'", s)
	}

	// test empty number
	expectedErr := errors.New("expected number, found nothing")
	buf = bytes.NewBufferString(", foo")
	scan = newScanner(buf)
	val, tok, err = scan.readNumber()
	if err == nil {
		t.Fatalf("readNumber ', foo': found no error, expected \"%v\"", expectedErr)
	}
	if tok != scanIllegalToken {
		t.Errorf("readNumber ', foo': found token '%v', expected '%v'", tok, scanIllegalToken)
	}
	if val != nil {
		t.Errorf("readNumber ', foo': found '%s', expected 'nil'", string(val))
	}

	// TODO: test other errors
}
func TestReadNull(t *testing.T) {
	// test 'null'
	buf := bytes.NewBufferString("ull")
	scan := newScanner(buf)
	tok, err := scan.readNull()
	if err != nil {
		t.Fatal("readNull 'ull': ", err)
	}
	if tok != scanNullVal {
		t.Errorf("readNull 'ull': found '%v', expected '%v'", tok, scanNullVal)
	}

	// test error
	expectedErr := errors.New("invalid null value: found 'f', expected 'u'")
	buf = bytes.NewBufferString("foo")
	scan = newScanner(buf)
	tok, err = scan.readNull()
	if err == nil {
		t.Fatalf("readNull 'foo': found no error, expected \"%v\"", expectedErr)
	}
	if err.Error() != expectedErr.Error() {
		t.Fatalf("readNull 'foo': found \"%v\", expected \"%v\"", err, expectedErr)
	}
	if tok != scanIllegalToken {
		t.Errorf("readNull 'foo': found '%v', expected '%v'", tok, scanIllegalToken)
	}
}

func TestReadBool(t *testing.T) {
	// test 'true'
	buf := bytes.NewBufferString("rue")
	scan := newScanner(buf)
	val, err := scan.readBool('t')
	if err != nil {
		t.Fatal("readBool 'rue': ", err)
	}
	if s := string(val); s != "true" {
		t.Errorf("readBool 'rue': found '%s', expected 'true'", s)
	}

	// test 'false'
	buf = bytes.NewBufferString("alse")
	scan = newScanner(buf)
	val, err = scan.readBool('f')
	if err != nil {
		t.Fatal("readBool 'alse': ", err)
	}
	if s := string(val); s != "false" {
		t.Errorf("readBool 'alse': found '%s', expected 'false'", s)
	}

	// test error
	expectedErr := errors.New("invalid boolean value: found 'o', expected 'a'")
	buf = bytes.NewBufferString("oo")
	scan = newScanner(buf)
	val, err = scan.readBool('f')
	if err == nil {
		t.Fatalf("readBool 'oo': found no error, expected \"%v\"", expectedErr)
	}
	if err.Error() != expectedErr.Error() {
		t.Fatalf("readBool 'oo': found \"%v\", expected \"%v\"", err, expectedErr)
	}
	if val != nil {
		t.Errorf("readBool 'oo': found '%s', expected 'nil'", string(val))
	}
}
