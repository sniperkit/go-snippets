package lex

import (
	"io/ioutil"
	"log"
	"path"
)

// LoadAllLangs loads all the language files.
func LoadAllLangs(langDir string) map[string]map[string]bool {
	langFiles, err := ioutil.ReadDir(langDir)
	if err != nil {
		log.Fatal(err)
	}
	langWords := make(map[string]map[string]bool)
	for _, langFile := range langFiles {
		name := langFile.Name()
		path := path.Join(langDir, name)
		langWords[name] = LoadLang(path)
	}
	return langWords
}
