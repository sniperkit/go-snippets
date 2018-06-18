package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"vueutil"
)

func NewAppData() interface{} {
	options := []*Option{
		NewOption(1, "Hello"),
		NewOption(2, "World"),
	}

	obj := js.Global.Get("Object").New()
	obj.Set("selected", 2)
	obj.Set("options", options)
	return obj
}

func main() {
	InitSelect2()

	o := vue.NewOption()
	o.Template = vueutil.GetTemplateById("demo-template")
	o.Data = NewAppData()
	v := o.NewViewModel()
	v.Mount("#el")
}
