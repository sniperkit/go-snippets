package gorgonzola

import (
	"html/template"
	"net/http"
)

type templateData map[string]interface{}

// Template is the user facing template helper structure
type Template struct {
	data   templateData
	layout *template.Template
	w      http.ResponseWriter
}

// NewTemplate creates new Template
func NewTemplate(w http.ResponseWriter) *Template {
	return &Template{
		data:   make(templateData),
		layout: template.New("layout.html"),
		w:      w,
	}
}

func (t *Template) render(filenames ...string) error {
	return template.Must(t.layout.ParseFiles(filenames...)).Execute(t.w, t.data)
}

func (t *Template) set(key string, data interface{}) {
	t.data[key] = data
}
