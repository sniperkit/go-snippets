package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
)

func NewAppData() interface{} {
	obj := js.Global.Get("Object").New()
	obj.Set("showModal", false)
	return obj
}

func main() {
	InitModal()

	o := vue.NewOption()
	o.Data = NewAppData()
	v := o.NewViewModel()
	v.Mount("#app")
}
