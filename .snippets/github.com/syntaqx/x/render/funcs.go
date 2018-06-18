// Package render is a package that provides functionality for easily rendering
// JSON, XML, binary data, and HTML templates.
package render

import (
	"html/template"
)

// defaultFuncs adds additional functions those built-in by Go's lightweight,
// but extensible template engine. Functions included this way are intended to
// provide solutions to problems consistently encountered.
var defaultFuncs = template.FuncMap{}
