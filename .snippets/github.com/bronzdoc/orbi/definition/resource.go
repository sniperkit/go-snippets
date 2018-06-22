package definition

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

// Resource represent a resource common actions
type Resource interface {
	Create(map[string]interface{}) error
	Name() string
	Children() []Resource
	ID() string
}

// Directory resource
type Directory struct {
	name     string
	id       string
	children []Resource
}

// NewDirectory creates a new Directory resource
func NewDirectory(name, id string, children []Resource) *Directory {
	return &Directory{
		name:     name,
		id:       id,
		children: children,
	}
}

// Create creates a new Directory resource in file system
func (d *Directory) Create(options map[string]interface{}) error {
	os.Mkdir(d.id, 0776)
	return nil
}

// Name gets a Directory resource name
func (d *Directory) Name() string {
	return d.name
}

// Children gets a Diretory resource children
func (d *Directory) Children() []Resource {
	return d.children
}

// ID gets a Directory resource id
func (d *Directory) ID() string {
	return d.id
}

// File resource
type File struct {
	name    string
	id      string
	content []byte
}

// NewFile creates a new File resource
func NewFile(name, id string, content []byte) *File {
	return &File{
		name:    name,
		id:      id,
		content: content,
	}
}

// Create creates a new File resource in file system
func (f *File) Create(options map[string]interface{}) error {
	file, err := os.Create(f.id)

	if err != nil {
		return fmt.Errorf("File Create: %s", err)
	}
	defer file.Close()

	// TODO improve how a definition.template and a definition.File interacts
	if f.isContentEmpty() {
		templatesPath := viper.GetString("TemplatesPath")

		// Check if there is a template for the file
		if f.hasTemplate(templatesPath) {
			vars := options["vars"].(map[string]string)
			path := templatesPath + "/" + f.name

			content, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("File Create: %s", err)
			}

			_, err = NewTemplate(f.name, content, vars).Execute(file)
			if err != nil {
				return fmt.Errorf("File Create: %s", err)
			}
		}
	} else {
		var vars map[string]string
		NewTemplate(f.name, f.content, vars).Execute(file)
	}

	return nil
}

// Content gets a File resource content
func (f *File) Content() []byte {
	return f.content
}

// SetContent sets a file resource content
func (f *File) SetContent(content []byte) {
	f.content = content
}

// Name gets a File resource name
func (f *File) Name() string {
	return f.name
}

// Children gets a File resource children
func (f *File) Children() []Resource {
	return nil
}

// ID gets a File resource id
func (f *File) ID() string {
	return f.id
}

func (f *File) isContentEmpty() bool {
	return f.content == nil
}

func (f *File) hasTemplate(templatesPath string) bool {
	_, err := os.Stat(templatesPath + "/" + f.name)
	return err == nil
}
