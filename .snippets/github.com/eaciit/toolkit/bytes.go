package toolkit

import (
	"bytes"
	"encoding/gob"
	"errors"
	"strings"
)

func ToBytesWithError(data interface{}, encoderId string) ([]byte, error) {
	encoderId = strings.ToLower(encoderId)
	if encoderId == "" {
		encoderId = "json"
	}

	if encoderId == "json" {
		return Jsonify(data), nil
	} else if encoderId == "gob" {
		b, e := EncodeByte(data)
		if e != nil {
			return nil, errors.New(e.Error())
		} else {
			return b, nil
		}
	}
	return nil, errors.New("Invalid encoderId method")
}

func ToBytes(data interface{}, encoderId string) []byte {
	b, e := ToBytesWithError(data, encoderId)
	if e != nil {
		return []byte{}
	} else {
		return b
	}
}

func FromBytes(b []byte, decoderId string, out interface{}) error {
	var e error
	decoderId = strings.ToLower(decoderId)
	if decoderId == "" {
		decoderId = "json"
	}

	if decoderId == "json" {
		e = Unjson(b, out)
	} else {
		e = DecodeByte(b, out)

	}
	return e
}

func DecodeByte(bytesData []byte, result interface{}) error {
	buf := bytes.NewBuffer(bytesData)
	dec := gob.NewDecoder(buf)
	e := dec.Decode(result)
	return e
}

func GetEncodeByte(obj interface{}) []byte {
	b, e := EncodeByte(obj)
	if e != nil {
		return new(bytes.Buffer).Bytes()
	}
	return b
}

func EncodeByte(obj interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	gw := gob.NewEncoder(buf)
	err := gw.Encode(obj)
	if err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}
