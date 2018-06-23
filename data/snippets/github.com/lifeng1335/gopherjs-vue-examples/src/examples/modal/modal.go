package main

import (
	"github.com/oskca/gopherjs-vue"
	"vueutil"
)

func InitModal() {
	o := vue.NewOption()
	o.Name = "modal"
	o.Template = vueutil.GetTemplateById("modal-template")
	o.NewComponent().Register("modal")
}
