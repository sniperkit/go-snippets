// Package webutil holds very simple functions and data that other packages may
// need in order to generate URLs and include assets.
package webutil

import (
	"fmt"
	"html/template"
	"path"

	"github.com/uoregon-libraries/gopkg/tmpl"
)

// Webroot must be set externally to tell us where we are within the main
// website, such as "/reports", and is used to generate absolute paths to
// various handlers and site assets
var Webroot string

// Static must be set externally if CSS, JS, images, and other static
// assets should have a custom path (default is "static")
var Static = "static"

// FuncMap offers up these helpers as a FuncMap suitable for use with the
// "github.com/uoregon-libraries/gopkg/tmpl" package
var FuncMap = tmpl.FuncMap{
	"FullPath":   FullPath,
	"StaticPath": StaticPath,
	"HomePath":   HomePath,
	"ImageURL":   ImageURL,
	"IncludeCSS": IncludeCSS,
	"RawCSS":     RawCSS,
	"IncludeJS":  IncludeJS,
	"RawJS":      RawJS,
	"Webroot":    func() string { return Webroot },
}

// FullPath uses the webroot, if not empty, to join together all the path parts
// with a slash, returning an absolute path to something
func FullPath(parts ...string) string {
	if Webroot != "" {
		parts = append([]string{Webroot}, parts...)
	}
	return path.Join(parts...)
}

// StaticPath returns the absolute path to static assets (CSS, JS, etc)
func StaticPath(dir, file string) string {
	return FullPath(Static, dir, file)
}

// HomePath returns the absolute path to the home page
func HomePath() string {
	return FullPath("")
}

// ImageURL constructs a path to <webroot>/<static>/images/filename
func ImageURL(filename string) string {
	return StaticPath("images", filename)
}

// IncludeCSS generates a <link> tag with an absolute path for including the
// given file's CSS.  ".css" is automatically appended to the filename for less
// verbose use.
func IncludeCSS(file string) template.HTML {
	var path = StaticPath("css", file+".css")
	return template.HTML(fmt.Sprintf(`<link rel="stylesheet" type="text/css" href="%s" />`, path))
}

// RawCSS generates a <link> tag with an absolute path for including the given
// file's CSS.  It doesn't assume the path is /css, and it doesn't auto-append
// ".css", but it does assume the file will be under the static path.
func RawCSS(file string) template.HTML {
	var path = StaticPath("", file)
	return template.HTML(fmt.Sprintf(`<link rel="stylesheet" type="text/css" href="%s" />`, path))
}

// IncludeJS generates a <script> tag with an absolute path for including the
// given file's JS.  ".js" is automatically appended to the filename for less
// verbose use.
func IncludeJS(file string) template.HTML {
	var path = StaticPath("js", file+".js")
	return template.HTML(fmt.Sprintf(`<script src="%s"></script>`, path))
}

// RawJS generates a <script> tag with an absolute path for including the given
// file's JS.  It doesn't assume the path is /js, and it doesn't auto-append
// ".js", but it does assume the file will be under the static path.
func RawJS(file string) template.HTML {
	var path = StaticPath("", file)
	return template.HTML(fmt.Sprintf(`<script src="%s"></script>`, path))
}
