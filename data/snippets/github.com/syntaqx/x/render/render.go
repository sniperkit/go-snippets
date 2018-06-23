// Package render is a pckage that provides functionality for easily rendering
// JSON, XML, binary data, and HTML templates. The majority of functionality is
// implemented through unrolled/render with the exception of subtle differences
// that are preferred by our core contributors.
package render

import (
	"html/template"
	"io"

	renderer "github.com/unrolled/render"
)

// Options is a struct for specifying configuration options to a Render object.
type Options = renderer.Options

// HTMLOptions is a struct used to override HTML-specific options at render.
type HTMLOptions = renderer.HTMLOptions

// Engine is the generic interface for all responses.
type Engine = renderer.Engine

// Render is the primary canonical interface exposed by this package.
type Render interface {
	Render(w io.Writer, e Engine, data interface{}) error
	Data(w io.Writer, status int, v []byte) error
	HTML(w io.Writer, status int, name string, binding interface{}, htmlOpt ...HTMLOptions) error
	JSON(w io.Writer, status int, v interface{}) error
	JSONP(w io.Writer, status int, callback string, v interface{}) error
	Text(w io.Writer, status int, v string) error
	XML(w io.Writer, status int, v interface{}) error
}

// New constructs a new default Render instances with the given options.
func New(options Options) Render {
	options.Funcs = append(options.Funcs, []template.FuncMap{defaultFuncs}...)
	return &render{
		render: renderer.New(options),
	}
}

type render struct {
	render *renderer.Render
}

// Render is the generic function to be called by each Engine, and can even be
// called by custom implementations.
func (r *render) Render(w io.Writer, e Engine, data interface{}) error {
	return r.render.Render(w, e, data)
}

// Data writes out the raw bytes as binary data.
func (r *render) Data(w io.Writer, status int, v []byte) error {
	return r.render.Data(w, status, v)
}

// HTML builds up the response from the specified template and bindings.
func (r *render) HTML(w io.Writer, status int, name string, binding interface{}, htmlOpt ...HTMLOptions) error {
	return r.render.HTML(w, status, name, binding, htmlOpt...)
}

// JSON marshals the given interface object and writes the JSON response.
func (r *render) JSON(w io.Writer, status int, v interface{}) error {
	return r.render.JSON(w, status, v)
}

// JSONP marshals the given interface object and writes the JSON response.
func (r *render) JSONP(w io.Writer, status int, callback string, v interface{}) error {
	return r.render.JSONP(w, status, callback, v)
}

// Text writes out a string as plain text.
func (r *render) Text(w io.Writer, status int, v string) error {
	return r.render.Text(w, status, v)
}

// XML marshals the given interface object and writes the XML response.
func (r *render) XML(w io.Writer, status int, v interface{}) error {
	return r.render.XML(w, status, v)
}
