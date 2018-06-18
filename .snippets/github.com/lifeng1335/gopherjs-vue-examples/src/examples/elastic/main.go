package main

import (
	"github.com/oskca/gopherjs-vue"
	"vueutil"
)

func main() {
	InitDraggableHeader()

	o := vue.NewOption()
	o.Data = vueutil.EmptyDataFunc()
	v := o.NewViewModel()
	v.Mount("#app")
}
