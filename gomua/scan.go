package gomua

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"os"
	"path/filepath"
)

var msgs []*Message

func processFile(filename string, in io.Reader, stdin bool) error {
	if in == nil {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		in = bytes.NewReader(b)
		/*
			f, err := os.Open(filename)
			if err != nil {
				return err
			}
			defer f.Close()
			in = f
		*/
	}

	m, err := mail.ReadMessage(in)
	if err != nil {
		return err
	}

	msg := ReadMessage(m)
	msg.Filename, err = filepath.Abs(filename)
	if err != nil {
		return err
	}

	msgs = append(msgs, msg)
	return nil
}

func visitFile(path string, f os.FileInfo, err error) error {
	if err == nil {
		err = processFile(path, nil, false)
	}
	if err != nil {
		// This was printing errors when reaading a directory, nerfing the msg for now.
		//fmt.Printf("%s", err)
	}
	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func Scan(path string) []*Message {
	switch dir, err := os.Stat(path); {
	case err != nil:
		fmt.Printf("%s", err)
	case dir.IsDir():
		walkDir(path)
	default:
		if err := processFile(path, nil, false); err != nil {
			fmt.Printf("%s", err)
		}
	}
	return msgs
}
