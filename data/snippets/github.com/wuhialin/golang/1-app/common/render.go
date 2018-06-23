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
	r = render.New(RenderOptions())
	return r
}

func RenderOptions() render.Options {
	o := render.Options{}
	o.Extensions = []string{".html"}
	dir, _ := os.Getwd()
	o.Directory = filepath.Join(dir, "templates")
	return o
}

func RenderHTML(o render.Options) *render.Render {
	return render.New(o)
}
