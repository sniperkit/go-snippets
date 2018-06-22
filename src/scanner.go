// Copyright 2014-2015 The project AUTHORS. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package src

import (
	"errors"
	"fmt"
	"io"
)

const bufsize = 32 * 1024

type scanToken int

const (
	scanIllegalToken scanToken = iota
	scanBeginObject
	scanEndObject
	scanBeginArray
	scanEndArray
	scanStringLit
	scanBoolLit
	scanInt64Lit
	scanFloat64Lit
	scanNullVal
	scanComma
)

var scanTokenToString = map[scanToken]string{
	scanIllegalToken: "illegal token",
	scanBeginObject:  "{",
	scanEndObject:    "}",
	scanBeginArray:   "[",
	scanEndArray:     "]",
	scanStringLit:    "string literal",
	scanBoolLit:      "boolean literal",
	scanInt64Lit:     "integer literal",
	scanFloat64Lit:   "float literal",
	scanNullVal:      "null",
	scanComma:        "comma",
}

func (tok scanToken) String() string {
	if str, ok := scanTokenToString[tok]; ok {
		return str
	}
	return "invalid scan token"
}

type scanner struct {
	r   io.Reader
	err error

	globPos int64 // position in the JSON input

	buf []byte
	pos int // position inside the buffer; must be -1 by default

	eof bool // true when the reader has been reached
}

func newScanner(r io.Reader) *scanner {
	return &scanner{r: r, buf: make([]byte, bufsize), pos: -1}
}

func (scan *scanner) nextKey() (string, error) {
	if err := scan.ignoreWhitespaces(); err != nil {
		if err == io.EOF {
			return "", errors.New("expected key, found EOF")
		}
		return "", err
	}

	if c, err := scan.read(); err != nil {
		return "", err
	} else if c != '"' {
		return "", fmt.Errorf("expected '\"', found '%c'", c)
	}

	var key string
	origPos := scan.pos
	if origPos > len(scan.buf)-1 {
		origPos = 0
	}
	for {
		c, err := scan.read()
		if err == io.EOF {
			return "", errors.New("expected key, found EOF")
		} else if err != nil {
			return "", err
		}

		// TODO: check character

		if c == '"' {
			break
		}

		if scan.pos == 0 || scan.pos > len(scan.buf)-1 {
			key += string(scan.buf[origPos:])
			origPos = 0
		}
	}

	key += string(scan.buf[origPos : scan.pos-1])

	if err := scan.ignoreWhitespaces(); err != nil {
		if err == io.EOF {
			return "", errors.New("expected ':', found EOF")
		}
		return "", err
	}

	if c, err := scan.read(); err == io.EOF {
		return "", errors.New("expected ':', found EOF")
	} else if err != nil {
		return "", err
	} else if c != ':' {
		return "", fmt.Errorf("expected ':', found '%c'", c)
	}

	return key, nil
}

func (scan *scanner) nextValue() ([]byte, scanToken, error) {
	if err := scan.ignoreWhitespaces(); err != nil {
		if err == io.EOF {
			return nil, scanIllegalToken, errors.New("expected value, found EOF")
		}
		return nil, scanIllegalToken, err
	}

	c, err := scan.read()
	if err != nil {
		if err == io.EOF {
			return nil, scanIllegalToken, errors.New("expected value, found EOF")
		}
		return nil, scanIllegalToken, err
	}

	switch {
	case isDigit(c) || c == '-' || c == '+':
		scan.back()
		return scan.readNumber()
	case c == 't' || c == 'f':
		b, err := scan.readBool(c)
		return b, scanBoolLit, err
	case c == '{': // beginning of an object literal
		return nil, scanBeginObject, nil
	case c == '}': // ending of an object literal
		return nil, scanEndObject, nil
	case c == '[': // beginning of an array literal
		return nil, scanBeginArray, nil
	case c == ']': // ending of an array literal
		return nil, scanEndArray, nil
	case c == '"': // beginning or ending of a string literal
		str, err := scan.readString()
		return str, scanStringLit, err
	case c == 'n': // null
		tok, err := scan.readNull()
		return nil, tok, err
	case c == ',':
		return nil, scanComma, nil
	}
	return nil, scanIllegalToken, fmt.Errorf("illegal value '%c'", c)
}

// peek reads the next value without consuming it.
func (scan *scanner) peek() (byte, error) {
	c, err := scan.read()
	if err != nil {
		return 0, err
	}
	scan.back()
	return c, nil
}

func (scan *scanner) read() (byte, error) {
	scan.globPos++
	if scan.pos == -1 || (scan.pos > len(scan.buf)-1 && !scan.eof) {
		n, err := scan.r.Read(scan.buf)
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			scan.eof = true
		}
		if n < bufsize {
			scan.buf = scan.buf[:n]
		}
		scan.pos = 0
	}
	if scan.pos >= len(scan.buf) {
		return 0, io.EOF
	}
	b := scan.buf[scan.pos]
	scan.pos++
	return b, nil
}

// readString reads a string literal.
func (scan *scanner) readString() ([]byte, error) {
	// length of the read string
	// XXX: overflow?
	var strLen int

	buf := make([]byte, bufsize)
	for {
		c, err := scan.read()
		if err != nil {
			return nil, err
		}
		if c == '"' {
			if strLen == 0 {
				// empty string
				return []byte{}, nil
			}
			if buf[strLen-1] != '\\' {
				break
			}
		}
		if strLen == len(buf) {
			buf = append(buf, c)
		} else {
			buf[strLen] = c
		}
		strLen++
	}

	return buf[:strLen], nil
}

// readNumber reads either an int 64 or f float 64.
func (scan *scanner) readNumber() ([]byte, scanToken, error) {
	tok := scanInt64Lit
	numLen := 0
	buf := make([]byte, bufsize, bufsize)
	for {
		c, err := scan.read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, scanIllegalToken, err
		}
		if c == ',' || c == '}' || c == ']' {
			if numLen == 0 {
				return nil, scanIllegalToken, errors.New("expected number, found nothing")
			}
			scan.back()
			break
		}
		if c == '.' {
			// "." cannot be at the first place in a floating point number
			if numLen == 0 {
				return nil, scanIllegalToken, errors.New("expected digit, found '.'")
			}
			// floating point number can only contain one "."
			if tok == scanFloat64Lit {
				return nil, scanIllegalToken, errors.New("unexpected character '.'")
			}
			tok = scanFloat64Lit
		} else if c == '-' || c == '+' {
			if numLen != 0 {
				return nil, scanIllegalToken, fmt.Errorf("symbol '%c' can only be at the first position of a number", c)
			}
		} else if !isDigit(c) {
			return nil, scanIllegalToken, fmt.Errorf("expected digit, found '%c'", c)
		}
		// bufsize is already bigger that the maximum possible size for a number,
		// therefore, if the numLen is bigger than the busize, we return an
		// error.
		if numLen > bufsize {
			return nil, scanIllegalToken, errors.New("number too long")
		}
		buf[numLen] = c
		numLen++
	}
	return buf[:numLen], tok, nil
}

// readNull reads a null value.
func (scan *scanner) readNull() (scanToken, error) {
	var err error
	step := func(expected byte) {
		if err != nil {
			return
		}
		var c byte
		if c, err = scan.read(); err != nil {
			return
		}
		if c != expected {
			err = fmt.Errorf("invalid null value: found '%c', expected '%c'", c, expected)
		}
	}
	// 'n' has already been consumed
	step('u')
	step('l')
	step('l')
	if err != nil {
		return scanIllegalToken, err
	}
	return scanNullVal, nil
}

// readBool reads a boolean value.
func (scan *scanner) readBool(first byte) ([]byte, error) {
	var err error
	step := func(expected byte) {
		if err != nil {
			return
		}
		var c byte
		if c, err = scan.read(); err != nil {
			return
		}
		if c != expected {
			err = fmt.Errorf("invalid boolean value: found '%c', expected '%c'", c, expected)
		}
	}
	// first character ('t' or 'f') has already been consumed
	var val []byte
	if first == 't' {
		step('r')
		step('u')
		step('e')
		val = []byte("true")
	} else {
		step('a')
		step('l')
		step('s')
		step('e')
		val = []byte("false")
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (scan *scanner) back() {
	if scan.pos > 0 {
		scan.pos--
	}
}

// ignoreWhitespaces consumes all whitespaces.
func (scan *scanner) ignoreWhitespaces() error {
	var c byte
	var err error
	// XXX: refactor when the tests pass
	for {
		if c, err = scan.read(); err != nil || !isWhitespace(c) {
			break
		}
		// Nothing, we just skip whitespaces.
	}

	// Since a non-whitespace has been read, we have to put it back to the
	// buffer so that it can be read again.
	scan.back()

	if err != io.EOF {
		return err
	}
	return nil
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\n' || c == '\t' || c == '\r'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
