// Package tmpl wraps a lot of html/template for easier use with common
// layout setup and sub-templates
package tmpl

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

// FuncMap aliases the html/template FuncMap for easier local use
type FuncMap template.FuncMap

// DefaultTemplateFunctions is a usable function map for getting some free,
// commonly needed template functionality
//
// Explanation of default functions:
//
//   - nl2br(string): For presenting user-entered text, escapes the string and
//     then converts all newlines to "<br />"
//   - raw(string): Turns a string into template.HTML so it can be displayed
//     without escaping.  This should be used carefully!
//   - comment(string): Prints an HTML comment for cases you want a comment
//     that's visible in the final HTML
var DefaultTemplateFunctions = FuncMap{
	"nl2br": func(s string) template.HTML {
		var escaped = template.HTMLEscaper(s)
		var replaced = strings.Replace(escaped, "\n", "<br />", -1)
		return template.HTML(replaced)
	},
	"raw":     func(s string) template.HTML { return template.HTML(s) },
	"comment": func(s string) template.HTML { return template.HTML(fmt.Sprintf("<!-- %s -->", s)) },
}

// Template wraps html/template's type in order to provide a name so a single
// template can be self-contained
type Template struct {
	*template.Template
	Name string
}

// Clone wraps html/template.Clone to also clone the name
func (t *Template) Clone() (*Template, error) {
	var tmpl, err = t.Template.Clone()
	return &Template{tmpl, t.Name}, err
}

// TRoot wraps template.Template for use to spawn "real" templates.  The TRoot
// is never meant to be directly rendered itself, but a top-level object for
// collecting the template path on disk, a layout template and shared templates
// (e.g., sidebar), and template functions for reuse in renderable templates
type TRoot struct {
	template *Template
	Path     string
}

// Root creates a new TRoot for use in spawning templates.  The name should
// match the main layout's name (as defined in the layout template) so
// execution of templates doesn't require a template.Lookup call, which can be
// somewhat error prone.
func Root(name, path string) *TRoot {
	var tmpl = &Template{template.New(name), name}
	var t = &TRoot{tmpl, path}

	return t
}

// Funcs allows adding template function maps to TRoots; this should be done
// before creating any templates, or else previously created templates won't
// get the newest function maps
func (t *TRoot) Funcs(fnList FuncMap) *TRoot {
	t.template.Funcs(template.FuncMap(fnList))
	return t
}

// Clone creates a copy of the TRoot for ease of creating sub-layouts.  Since
// TRoots cannot be executed externally, we don't have the possibility of
// returning an error.
func (t *TRoot) Clone() *TRoot {
	var clone, _ = t.template.Clone()
	return &TRoot{clone, t.Path}
}

// Name exposes the underlying template's name
func (t *TRoot) Name() string {
	return t.template.Name
}

// ReadPartials parses the given files into the TRoot instance for gathering
// things like the top-level layout, navigation elements, etc.  The list of
// files is relative to the TRoot's Path.  Returns on the first error
// encountered, if any.
func (t *TRoot) ReadPartials(files ...string) error {
	for _, file := range files {
		var _, err = t.template.ParseFiles(filepath.Join(t.Path, file))
		if err != nil {
			return err
		}
	}

	return nil
}

// MustReadPartials calls ReadPartials and panics on any error
func (t *TRoot) MustReadPartials(files ...string) {
	var err = t.ReadPartials(files...)
	if err != nil {
		panic(err)
	}
}

// Build clones the root (for layout, funcs, etc) and parses the given file in
// the clone.  The returned template is the clone, and is safe to alter without
// worrying about breaking the root.
func (t *TRoot) Build(path string) (*Template, error) {
	var tNew, err = t.template.Clone()
	if err != nil {
		return nil, err
	}

	_, err = tNew.ParseFiles(filepath.Join(t.Path, path))
	if err != nil {
		return nil, err
	}

	tNew.Name = path
	return tNew, nil
}

// MustBuild calls Build and panics on any error
func (t *TRoot) MustBuild(path string) *Template {
	var tmpl, err = t.Build(path)
	if err != nil {
		panic(err)
	}
	return tmpl
}

// Template returns an empty template associated with this TRoot so we can use
// it for errors without fake "empty" templates
func (t *TRoot) Template() *Template {
	return t.Clone().template
}
