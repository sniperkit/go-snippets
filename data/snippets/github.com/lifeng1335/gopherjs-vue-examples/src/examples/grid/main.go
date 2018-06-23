package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"strings"
)

type Data struct {
	*js.Object
	Name  string `js:"name"`
	Power int    `js:"power"`
}

func NewData(Name string, Power int) *Data {
	d := &Data{
		Object: js.Global.Get("Object").New(),
	}
	d.Name = Name
	d.Power = Power
	return d
}

type AppData struct {
	*js.Object
	SearchQuery string   `js:"searchQuery"`
	GridColumns []string `js:"gridColumns"`
	GridData    []*Data  `js:"gridData"`
}

func NewAppData() interface{} {
	ad := &AppData{
		Object: js.Global.Get("Object").New(),
	}
	ad.SearchQuery = ""
	ad.GridColumns = []string{"name", "power"}
	ad.GridData = []*Data{
		NewData("Bruce Lee", 9000),
		NewData("Jackie Chan", 7000),
		NewData("Chuck Norris", 9999),
		NewData("Jet Li", 8000),
	}
	return ad
}

func main() {
	RegisterFilter()
	InitGrid()

	o := vue.NewOption()
	o.Data = NewAppData()

	v := o.NewViewModel()
	v.Mount("#demo")
}

func RegisterFilter() {
	vue.NewFilter(func(v *js.Object) interface{} {
		return strings.ToUpper(v.String())
	}).Register("capitalize")
}
