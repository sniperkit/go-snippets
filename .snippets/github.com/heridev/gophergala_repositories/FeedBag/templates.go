package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Template struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Event     string `json:"event"`
	Condition string `json:"condition"`
	Template  []byte
}

func ParseTemplatesDir(templatesDir string) ([]*Template, error) {
	files, err := ioutil.ReadDir(templatesDir)
	checkErr(err, "Read templates dir failed:")

	templates := []*Template{}

	for _, fileInfo := range files {
		tmpl, err := parseTemplate(templatesDir, fileInfo)
		if err != nil {
			log.Println(err.Error())
		} else {
			templates = append(templates, tmpl)
		}
	}

	return templates, nil
}

func parseTemplate(dir string, file os.FileInfo) (*Template, error) {
	var err error
	var template *Template

	if strings.HasSuffix(file.Name(), "tmpl") {

		filePath := fmt.Sprintf("%s/%s", dir, file.Name())
		log.Println("Parsing template:", filePath)

		file, err := os.Open(filePath)
		if err != nil {
			log.Println("Error opening template file:", filePath)
			return nil, err
		}
		defer file.Close()
		reader := bufio.NewReader(file)

		// Read header delimiter of <!--
		line, err := reader.ReadString('\n')
		if line != "<!--\n" {
			return nil, errors.New("Template file didn't begin with <!--")
		}

		// Read the json config package up until the --> line
		config := []byte{}
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					return nil, err
				} else {
					return nil, errors.New("Invalid template file")
				}
			}

			if string(line) == "-->\n" {
				break
			}

			config = append(config, line...)
		}

		err = json.Unmarshal(config, &template)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid template config: %q", string(config)))
		}

		if template.Id == "" {
			return nil, errors.New(`Template is missing "id" field`)
		}

		if template.Name == "" {
			return nil, errors.New(`Template is missing "name" field`)
		}

		if template.Event == "" {
			return nil, errors.New(`Template is missing "event" field`)
		}

		// Now store the actual DOM for the template output
		dom := []byte{}
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					return nil, err
				} else {
					break
				}
			}

			dom = append(dom, line...)
		}
		template.Template = dom
	}

	return template, err
}
