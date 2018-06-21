package common

import (
	"github.com/unrolled/render"
	"os"
	"path/filepath"
)

var r *render.Render

func Render() *render.Render {
	if r != nil {
		return r
	}
	o := render.Options{}
	o.Extensions = []string{".html"}
	dir, _ := os.Getwd()
	o.Directory = filepath.Join(dir, "4-pages", "templates")
	r = render.New(o)
	return r
}
