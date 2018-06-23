package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// convertFromFunc defines a function which returns
// items which should be converted by convertToFunc
type convertFromFunc func() (map[string]string, error)

// convertToFunc defines a function which converts
// items which were returned by convertFromFunc
type convertToFunc func(snippets map[string]string) error

// converterMap is a struct which holds conversion rules
type converterMap struct {
	Converters map[string]ConvertRule
}

// NewConverterMap returns new converterMap instance
func NewConverterMap() *converterMap {
	cm := converterMap{
		map[string]ConvertRule{},
	}

	return &cm
}

func (cm *converterMap) addConvertRule(ruleKey string, ruleFrom convertFromFunc, ruleTo convertToFunc) {
	cm.Converters[ruleKey] = ConvertRule{
		ruleFrom,
		ruleTo,
	}
}

func (cm *converterMap) convert() {
	for k, v := range cm.Converters {
		fmt.Printf("Processing rule under key '%s' \n", k)

		snippets, err := v.RuleFrom()
		if err != nil {
			fmt.Printf("/!\\ Error getting snippets under key '%s': %s \n", k, err)
			continue
		}

		v.RuleTo(snippets)

		fmt.Printf("Processing rule under key '%s' completed \n", k)
	}
}

// ConvertRule holds two rules - how to get input data and how to get output data
type ConvertRule struct {
	RuleFrom convertFromFunc
	RuleTo   convertToFunc
}

const (
	defDirMode          = 0755
	defFileMode         = 0644
	snippetExtension    = ".sublime-snippet"
	snippetTemplateFile = "sublime-snippet.sample"
)

var sublimeUserPath string
var snippetTemplate string
var cm *converterMap
var cfg *ini.File
var replacer *strings.Replacer

func init() {
	cm = NewConverterMap()
	sublimeUserPath = os.Getenv("SUBLIME_USER_PATH")

	cfg = ini.Empty()
	cfg.Append(os.Getenv("GEANY_SNIPPETS_CONF"))

	// Can move to a certain function and call once if `snippetTemplate` is empty
	curPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("/!\\ Err finding current path!")
	}
	snippetTemplateBytes, err := ioutil.ReadFile(curPath + "/" + snippetTemplateFile)
	if err != nil {
		panic("/!\\ Err reading file contents!")
	}
	snippetTemplate = string(snippetTemplateBytes)

	replacer = strings.NewReplacer(
		"\\n", "\n",
		"\\t", "\t",
		"\\s", " ",
		"%brace_open%", "{\n\t",
		"%brace_close%", "}\n",
		"%cursor%", "${1:content}",
		"%block%", "\n{\n\t${1:content}\n}",
		"%block_cursor%", "{\n\t${1:content}\n}\n",
	)
}

func main() {
	//cm.addConvertRule("convert_general", convFrom, convTo)
	//cm.addConvertRule("convert_php", convPHPFrom, convPHPTo)
	//cm.addConvertRule("convert_javascript", convFrom, convTo)
	cm.addConvertRule("convert_go", convGoFrom, convGoTo)
	cm.convert()

	fmt.Println("All convertation rules proceeded")
}

// convPHPFrom gets all snippets under PHP config section
func convPHPFrom() (map[string]string, error) {
	return getSectionKeys("PHP")
}

func convPHPTo(snippets map[string]string) error {
	snippetsFolder := sublimeUserPath + "/php"
	var err error

	err = createFolder(snippetsFolder)
	if err != nil {
		return err
	}

	return createSnippetsInFolder(snippets, snippetsFolder, "source.php")
}

func convGoFrom() (map[string]string, error) {
	return getSectionKeys("Go")
}

func convGoTo(snippets map[string]string) error {
	snippetsFolder := sublimeUserPath + "/go"
	var err error

	err = createFolder(snippetsFolder)
	if err != nil {
		return err
	}

	return createSnippetsInFolder(snippets, snippetsFolder, "source.go")
}

// getSectionKeys gets keys under certain section and error if something went wrong
func getSectionKeys(sectionName string) (map[string]string, error) {
	keys := map[string]string{}

	section, err := cfg.GetSection(sectionName)
	if err != nil {
		return keys, err
	}

	keyNames := section.KeyStrings()
	for _, v := range keyNames {
		keys[v] = cfg.Section(sectionName).Key(v).String()
	}

	return keys, nil
}

// createFolder creates a folder with mode equals defDirMode
func createFolder(path string) error {
	return os.Mkdir(path, defDirMode)
}

// createSnippetsInFolder creates snippets files
func createSnippetsInFolder(snippets map[string]string, folderPath, snippetScope string) error {
	var err error

	//var counter int
	for k, v := range snippets {
		//counter++

		snippetContent := getSnippetContent(k, v, snippetScope)
		createSnippetFile(folderPath, k, snippetContent)
		// todo remove
		/*if counter == 10 {
			break
		}*/
	}

	return err
}

func getSnippetContent(trigger, content, scope string) string {
	return fmt.Sprintf(
		snippetTemplate,
		replacer.Replace(content),
		trigger,
		scope,
	)
}

func createSnippetFile(path, filename, content string) error {
	var err error
	name := path + "/" + filename + snippetExtension

	if fh, err := os.Create(name); err == nil {
		_ = os.Chmod(name, defFileMode)

		_, err = fh.Write([]byte(content))

		fh.Close()
	}

	return err
}
