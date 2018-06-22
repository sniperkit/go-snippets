package definition

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Definition represents a defined project srtucture
type Definition struct {
	Context      string
	Resources    []map[interface{}]interface{}
	ResourceTree *Tree
	Options      map[string]interface{}
}

// New creates a new Definiton
func New(definitionFormat interface{}, options map[string]interface{}) *Definition {
	var def *Definition

	switch ftype := definitionFormat.(type) {
	default:
		log.Fatalf("%s: Is an invalid format type to create a definition", ftype)
	case string:
		def = newFromFile(definitionFormat.(string), options)
	case map[interface{}]interface{}:
		def = newFromMap(definitionFormat.(map[interface{}]interface{}), options)
	}

	def.ResourceTree = NewTree(def.Context, def.Resources)

	return def
}

// Create creates a defined project structure in file system
func (d *Definition) Create() error {
	errChan := make(chan error)
	tree := d.ResourceTree

	// Check definition context exists
	if _, err := os.Stat(tree.Root().ID()); err != nil {
		return fmt.Errorf(
			`Definition Create: Expected context: "%s" to exist.`,
			tree.Root().Name(),
		)
	}

	go func() {
		defer close(errChan)
		tree.Traverse(func(r Resource) {
			errChan <- r.Create(d.Options)
		})
	}()

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("Definition Create: %s", err)
		}
	}

	return nil
}

// Search search for a Resource in a Definition struct
func (d *Definition) Search(pattern string) Resource {
	resource := make(chan Resource)
	tree := d.ResourceTree

	go func() {
		defer close(resource)
		tree.Traverse(func(r Resource) {
			resource <- r
		})
	}()

	for resourceWanted := range resource {
		pattern := fmt.Sprint(d.Context, "/", pattern)
		if resourceWanted.ID() == pattern {
			return resourceWanted
		}
	}

	return nil
}

func newFromFile(fileName string, options map[string]interface{}) *Definition {
	def := Definition{Options: options}

	definitionContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	if err = yaml.Unmarshal(definitionContent, &def); err != nil {
		log.Fatal(err)
	}

	return &def
}

func newFromMap(mapDefinition map[interface{}]interface{}, options map[string]interface{}) *Definition {
	return &Definition{
		Context:   mapDefinition["context"].(string),
		Resources: mapDefinition["resources"].([]map[interface{}]interface{}),
		Options:   options,
	}
}
