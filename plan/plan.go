package plan

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/bronzdoc/orbi/definition"
	"github.com/spf13/viper"
)

// ExecCommand allows to Have higher control of exec.Command,
// and will allow us to mock it easier in tests...
var ExecCommand = exec.Command

// Plan represents a plan to execute
type Plan struct {
	definition *definition.Definition
}

// New creates a new Plan
func New(definition *definition.Definition) *Plan {
	return &Plan{
		definition: definition,
	}
}

// Factory constructs a Plan with a default Definition struct
func Factory(planName string, options map[string]interface{}) *Plan {
	definitionPath := fmt.Sprintf("%s/%s/definition.yml", viper.GetString("PlansPath"), planName)
	definition := definition.New(definitionPath, options)
	return New(definition)
}

// Execute generates a project definnition in file system
func (p *Plan) Execute() error {
	return p.definition.Create()
}

// List gets all the plans in the plans path
func List() (list []string) {
	plansPath := viper.GetString("PlansPath")
	files, _ := ioutil.ReadDir(plansPath)
	for _, f := range files {
		list = append(list, f.Name())
	}
	return list
}

// Edit allows to edit a plan definition
func Edit(planName string) error {
	if !planExists(planName) {
		return fmt.Errorf(`plan "%s" doesn't exists`, planName)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf(`$EDITOR is empty, could not edit "%s" plan`, planName)
	}

	definitionPath := fmt.Sprintf("%s/%s/definition.yml", viper.GetString("PlansPath"), planName)
	cmd := ExecCommand(editor, definitionPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

// Definition gets a definition struct for a new plan
func Definition(planName string, options map[string]interface{}) *definition.Definition {
	// Default structure for a new plan
	resources := []map[interface{}]interface{}{
		{
			"dir": map[interface{}]interface{}{
				"name": planName,
				"dir": map[interface{}]interface{}{
					"name": "templates",
				},
				"files": []interface{}{
					"definition.yml",
				},
			},
		},
	}

	mapDefinition := map[interface{}]interface{}{
		"context":   viper.GetString("PlansPath"),
		"resources": resources,
	}

	def := definition.New(mapDefinition, options)
	resource := def.Search(planName + "/" + "definition.yml")

	file := resource.(*definition.File)
	file.SetContent(
		[]byte(`---
context: .
resources:
  - dir:
     name: dir_name_sample
     files:
      - file_name_sample`),
	)

	return def
}

// Get downloads a plan from a plan repo url
func Get(planURL string) error {
	return Clone(planURL)
}

func planExists(planName string) bool {
	plansPath := viper.GetString("PlansPath")
	files, _ := ioutil.ReadDir(plansPath)
	for _, f := range files {
		if f.Name() == planName {
			return true
		}
	}
	return false
}
