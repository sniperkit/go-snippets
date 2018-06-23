package main

import (

	"bufio"
	"fmt"
	"os"
	"regexp"

)

// XXX fix extensions to handle more than one

type Template struct {
	Name string
	Extensions string
	Structures map[string]string
	Types map[string]string
	Names map[string]string
	Comments string
	Keywords map[string]string
	Expressions map[string]string
	Libraries map[string]string
}

func templateParse(templatePath string, verbosePtr bool, quietPtr bool, statusPtr bool) map[string]*Template {

	count := 0

	files   := getFiles(templatePath)
	r, _    := regexp.Compile("^templates/\\w+\\.template$")
	skip, _ := regexp.Compile("^templates/example\\.template$")

	for _, file := range files {
		// XXX
		if(r.MatchString(file) && ! skip.MatchString(file)) {
			count++
			if (!quietPtr && verbosePtr) {
				fmt.Println("found a specific template:", file)
			}
		}
	}
	// XXX fix handling count
	templates := map[string]*Template{}

	if (!quietPtr && verbosePtr) {
		fmt.Println("found a few templates:", count)
	}

	for _, file := range files {
		// XXX
		if(r.MatchString(file) && ! skip.MatchString(file)) {
			t := templateLoad(file, verbosePtr, quietPtr, statusPtr)
			templates[t.Name] = t
		}
	}

	return templates

}

func templateLoad(templateFile string, verbosePtr bool, quietPtr bool, statusPtr bool) *Template {

	t := new(Template)
	nameparts := regexp.MustCompile("[/.]").Split(templateFile, -1)
	t.Name = nameparts[1]
	//fmt.Println("setting name to", nameparts[1], "from ", templateFile)
	section := "unknown"

	sectionr, _ := regexp.Compile("^[\\w\\s]+\\:$")
	cleanr, _   := regexp.Compile("^\\s*")
	//paramr, _   := regexp.Compile("^\\s+(\\w+)\\: (.*)$")

	file, err := os.Open(templateFile)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if(sectionr.MatchString(scanner.Text())) {
			section = scanner.Text()[:len(scanner.Text())-1]
			if (!quietPtr && verbosePtr) {
				fmt.Println("found section", section, "from file", templateFile)
			}
		} else {
			parts := regexp.MustCompile(": ").Split(scanner.Text(), 2)
			param := cleanr.ReplaceAllString(parts[0], "")
			if(len(parts) > 1) {
				if (!quietPtr && verbosePtr) {
					fmt.Println("adding content match", parts[1], "to param", param)
				}
				// XXX literal issue
				if section == "names" {
					if (!quietPtr && verbosePtr) {
						fmt.Println(colorize("found names section"))
					}
					if param == "comments" {
						if (!quietPtr && verbosePtr) {
							fmt.Println("found comments section(s)", parts[1])
						}
						t.Comments = parts[1]
					}
				}
			} else {
				if(len(param) > 0) {
					if (!quietPtr && verbosePtr) {
						fmt.Println("adding single param", param, "to section", section)
					}
					// XXX literal issue here, ugly
					if section == "extensions" {
						t.Extensions = param
					}
				} else {
					//fmt.Println("skipping blank line", param)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		check(err)
	}

	//fmt.Printf("template has this struct: %v\n", t)

	return t

}
