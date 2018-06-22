package cassandredis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

//go:generate stringer -type=respType
type respType uint

const (
	respSimpleString respType = iota + 1
	respError
	respInteger
	respBulkString
	respArray
	respUnknown

	crlf = "\r\n"
)

type respSerializable interface {
	SerializeTo(w io.Writer) error
}

type respValue struct {
	respType respType
	value    interface{}
}

func (c *respValue) String() string {
	return fmt.Sprintf("type: %s value: %s", c.respType, c.value)
}

type respSimpleStringValue struct {
	value string
}

func (r *respSimpleStringValue) String() string {
	return r.value
}

type respErrorValue struct {
	message string
}

func (r *respErrorValue) String() string {
	return r.message
}

func (r *respErrorValue) SerializeTo(w io.Writer) error {
	var buf bytes.Buffer

	buf.WriteRune('-')
	buf.WriteString(r.message)
	buf.WriteString(crlf)

	_, err := io.Copy(w, &buf)
	return err
}

type respIntegerValue struct {
	value int64
}

func (r *respIntegerValue) SerializeTo(w io.Writer) error {
	var buf bytes.Buffer

	buf.WriteRune(':')
	buf.WriteString(strconv.FormatInt(r.value, 10))
	buf.WriteString(crlf)

	_, err := io.Copy(w, &buf)
	return err
}

type respBulkStringValue struct {
	value []byte
}

func (r *respBulkStringValue) SerializeTo(w io.Writer) error {
	var buf bytes.Buffer

	buf.WriteRune('$')
	buf.WriteString(strconv.Itoa(len(r.value)))
	buf.WriteString(crlf)
	buf.Write(r.value)
	buf.WriteString(crlf)

	_, err := io.Copy(w, &buf)
	return err
}

func (r *respBulkStringValue) String() string {
	return string(r.value)
}

type respArrayValue struct {
	length int
	values []interface{}
}

func (r *respArrayValue) String() string {
	var buf bytes.Buffer
	buf.WriteRune('[')
	for i, v := range r.values {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(v.(fmt.Stringer).String())
	}
	buf.WriteRune(']')

	return buf.String()
}

func (r *respArrayValue) SerializeTo(w io.Writer) error {
	var buf bytes.Buffer

	buf.WriteRune('*')
	buf.WriteString(strconv.Itoa(len(r.values)))
	buf.WriteString(crlf)
	for _, v := range r.values {
		v.(respSerializable).SerializeTo(&buf)
	}

	_, err := io.Copy(w, &buf)
	return err
}

type respArrayStreamingValue struct {
	length int
}

type protocol struct {
	br *bufio.Reader
}

func (p *protocol) read() (*respValue, error) {
	respType, value, err := p.readRespValue()
	if err != nil {
		return nil, err
	}

	return &respValue{
		respType: respType,
		value:    value,
	}, nil
}

func (p *protocol) readRespValue() (respType, interface{}, error) {
	b, err := p.br.ReadByte()
	if err != nil {
		return respUnknown, nil, err
	}

	respType, err := parseRespType(b)
	if err != nil {
		return respUnknown, nil, err
	}

	var value interface{}
	switch respType {
	case respArray:
		value, err = p.readRespArrayValue()
		if err != nil {
			return respType, nil, err
		}
	case respBulkString:
		value, err = p.readRespBulkString()
		if err != nil {
			return respType, nil, err
		}
	}

	return respType, value, nil
}

func (p *protocol) readRespArrayValue() (*respArrayValue, error) {
	bytes, err := p.br.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	line := string(bytes[:len(bytes)-2])
	length, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	v := &respArrayValue{length: length}

	for i := 0; i < length; i++ {
		_, iv, err := p.readRespValue()
		if err != nil {
			return nil, err
		}
		v.values = append(v.values, iv)
	}

	return v, nil
}

func (p *protocol) readRespBulkString() (*respBulkStringValue, error) {
	bytes, err := p.br.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	line := string(bytes[:len(bytes)-2])
	length, err := strconv.Atoi(line)
	if err != nil {
		return nil, err
	}

	bytes, err = p.br.ReadBytes('\n')
	if err != nil {
		return nil, err
	} else if len(bytes)-2 != length {
		return nil, fmt.Errorf("should have read %d bytes but only read %d", length, len(bytes))
	}

	return &respBulkStringValue{bytes[:len(bytes)-2]}, nil
}

func parseRespType(b byte) (respType, error) {
	switch b {
	case '+':
		return respArray, nil
	case '-':
		return respError, nil
	case ':':
		return respInteger, nil
	case '$':
		return respBulkString, nil
	case '*':
		return respArray, nil
	default:
		return respUnknown, fmt.Errorf("unable to identify RESP type from rune '%x'", b)
	}
}
