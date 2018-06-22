package xml

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
)

var dataPath = "sample/sample.xml"
var dataStr = ""

func init() {
	absPath, _ := filepath.Abs(dataPath)
	bytes, err := ioutil.ReadFile(absPath)

	if err == nil {
		dataStr = string(bytes)
	} else {
		panic(err)
	}
}

func TestRelativePath(t *testing.T) {
	absPath, _ := filepath.Abs(dataPath)
	log.Printf("data path ==> %s\n", absPath)
}

func TestMarshal(t *testing.T) {
	type book struct {
		Name   string `xml:"name"`
		Author string `xml:"author"`
		Press  string `xml:"press"`
	}

	v := book{Name: "Go", Author: "Google", Press: "Github"}
	expectedContent := "<book><name>Go</name><author>Google</author><press>Github</press></book>"

	bs, err := xml.Marshal(&v)
	if err != nil {
		log.Printf("error: %s", err)
	} else {
		log.Printf("content: %s", string(bs))
		if string(bs) != expectedContent {
			log.Printf("content by marshalling is not as expected, marshalling: [%s], expected: [%s]",
				string(bs), expectedContent)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	type book struct {
		Name   string `xml:"name"`
		Author string `xml:"author"`
		Press  string `xml:"press"`
	}

	xmlStr := "<book>" +
		"<name>The Go Programming Language</name>" +
		"<author>Brian W. Kernighan/Alan Donovan</author>" +
		"<press>Addison-Wesley Professional</press>" +
		"</book>"

	v := book{}
	err := xml.Unmarshal([]byte(xmlStr), &v)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	log.Printf("book: %v", v)
}
