package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var (
	InternalServerError = errors.New("InternalServerError")
	WrongOutput         = errors.New("Wrong Output")
)

func process(snippetId string, answer string) (string, error) {
	f, err := ioutil.ReadFile(snippetId)
	if err != nil {
		return "", InternalServerError
	}
	fStr := strings.Replace(string(f), Placeholder, answer, 1)
	output, err := execute([]byte(fStr))
	if err != nil && output == "" {
		return output, InternalServerError
	} else if err != nil {
		return output, errors.New(output)
	}
	f, err = ioutil.ReadFile(snippetId + ".output")
	if err != nil {
		return output, InternalServerError
	}
	if strings.TrimSpace(output) != strings.TrimSpace(string(f)) {
		return output, WrongOutput
	}
	return output, nil
}

func execute(srcCode []byte) (out string, err error) {
	f, err := ioutil.TempFile(os.TempDir(), "gophergame")
	if err != nil {
		return
	}
	fileName := f.Name() + ".go"
	f, err = os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(fileName, srcCode, os.ModeTemporary) // f.Write(srcCode)
	if err != nil {
		return
	}
	goBinary, err := exec.LookPath("go")
	if err != nil {
		return
	}
	c := exec.Command(goBinary, "run", fileName)
	buffer := bytes.NewBufferString("")
	c.Stdout = buffer
	c.Stderr = buffer
	err = c.Start()
	if err != nil {
		return
	}
	err = c.Wait()

	os.Remove(fileName)

	//TODO find better way to filter output
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "$GOROOT") < 0 && strings.Index(line, "$GOPATH") < 0 && strings.Index(line, "/go/src/") < 0 {
			out += line + "\n"
		}
	}

	out = strings.Replace(out, fileName, "Line", -1)
	out = strings.Replace(out, "# command-line-arguments", "", -1)
	out = strings.TrimSpace(out)
	return
}
