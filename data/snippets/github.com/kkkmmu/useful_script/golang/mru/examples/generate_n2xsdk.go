package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"util"
)

var MethodR = regexp.MustCompile("[[:space:]]*(?P<object>[[:alnum:]]+)[[:space:]]+(?P<method>[[:alnum:]]+)[[:space:]]+{(?P<params>[[:alnum:][:space:]]*)}")

type Method struct {
	Object     string
	Name       string
	Parameters string
}

func main() {
	content, err := ioutil.ReadFile("n2xapilist.txt")
	if err != nil {
		panic(err)
	}

	structs := make(map[string]string, 1000)
	methods := make(map[string][]*Method, 1000)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "ListMethodParameters") {
			continue
		}

		i := strings.Index(line, "}")
		if i < 0 {
			continue
		}

		line = line[:i+1]
		line = strings.Replace(line, "\"", "", -1)
		line = strings.Replace(line, ":", "", -1)
		//fmt.Println(line)

		line = strings.TrimSpace(line)
		fields := strings.Split(line, " ")
		if _, ok := structs[fields[0]]; !ok {
			structs[fields[0]] = fields[0]
		}

		if _, ok := methods[fields[0]]; !ok {
			methods[fields[0]] = make([]*Method, 0, 100)
		}

		matches := MethodR.FindStringSubmatch(line)
		if len(matches) == 4 {
			methods[fields[0]] = append(methods[fields[0]], &Method{Object: matches[1], Name: matches[2], Parameters: matches[3]})
		}
	}

	for obj, ms := range methods {
		os.Remove(fmt.Sprintf("%s.go", strings.Replace(obj, "Agt", "", -1)))
		util.AppendToFile(fmt.Sprintf("%s.go", strings.Replace(obj, "Agt", "", -1)), []byte("package n2xsdk\n\n"))
		util.AppendToFile(fmt.Sprintf("%s.go", strings.Replace(obj, "Agt", "", -1)), []byte(fmt.Sprintf("type %s struct {\n Handler string\n}\n\n", strings.Replace(obj, "Agt", "", -1))))

		for _, m := range ms {
			if strings.Contains(m.Name, "Get") || strings.Contains(m.Name, "List") {
				util.AppendToFile(fmt.Sprintf("%s.go", strings.Replace(obj, "Agt", "", -1)), []byte(fmt.Sprintf("func(np *%s) %s ()(string, error) {\n //parameters: %s\n //%s %s\n return \"\", nil\n}\n\n", strings.Trim(m.Object, "Agt"), m.Name, m.Parameters, m.Object, m.Name)))
			} else {
				util.AppendToFile(fmt.Sprintf("%s.go", strings.Replace(obj, "Agt", "", -1)), []byte(fmt.Sprintf("func(np *%s) %s () error {\n //parameters: %s\n //%s %s\n return nil\n}\n\n", strings.Trim(m.Object, "Agt"), m.Name, m.Parameters, m.Object, m.Name)))
			}
		}
	}
}
