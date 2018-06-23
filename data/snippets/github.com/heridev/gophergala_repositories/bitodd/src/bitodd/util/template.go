package util

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateModel struct {
	Params   map[interface{}]interface{}
}

func GetTemplateBase(name string) *template.Template {

	t := template.New(filepath.Base(name))

	return t
}

func GetNamedTemplate(name string) *template.Template {
	return template.Must(GetTemplateBase(name).ParseFiles(name))
}

func RenderTemplate(tmpl *template.Template, w http.ResponseWriter, model *TemplateModel) {

	err := tmpl.Execute(w, model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
