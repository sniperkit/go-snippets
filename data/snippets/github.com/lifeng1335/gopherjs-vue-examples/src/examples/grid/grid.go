package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-json"
	"github.com/oskca/gopherjs-vue"
	"sort"
	"strconv"
	"strings"
	"vueutil"
)

type GridData struct {
	*js.Object
	SortKey    string         `js:"sortKey"`
	SortOrders map[string]int `js:"sortOrders"`
}

type GridProps struct {
	*js.Object
	FilterKey string   `js:"filterKey"`
	Columns   []string `js:"columns"`
	Data      []*Data  `js:"data"`
}

type Grid struct {
	Data  *GridData
	Props *GridProps
}

func NewGrid(vm *vue.ViewModel) *Grid {
	return &Grid{
		Data: &GridData{
			Object: vm.Data,
		},
		Props: &GridProps{
			Object: vueutil.PropsData(vm),
		},
	}
}

func (g *Grid) SyncViewModel(vm *vue.ViewModel) {
	// [Vue warn]: Avoid replacing instance root $data. Use nested data properties instead.
	//vm.Data = t.Date.Object
	keys := js.Keys(g.Data.Object)
	for _, v := range keys {
		vm.Data.Set(v, g.Data.Get(v))
	}
	vm.Get("$options").Set("propsData", g.Props.Object)
}

func NewGridData() interface{} {
	g := &GridData{
		Object: js.Global.Get("Object").New(),
	}
	g.SortKey = ""
	g.SortOrders = nil
	return g
}

func (g *Grid) InitGridData() {
	SortOrders := map[string]int{}
	for _, v := range g.Props.Columns {
		SortOrders[v] = 1
	}
	g.Data.SortKey = ""
	g.Data.SortOrders = SortOrders
	println("InitGrid", json.Stringify(g))
}

func (g *Grid) SortBy(keyObj *js.Object) {
	key := keyObj.String()
	g.Data.SortKey = key

	order := g.Data.SortOrders
	order[key] = order[key] * -1

	g.Data.SortOrders = order
}

func (g *Grid) FilteredData() interface{} {
	fk := g.Props.FilterKey
	data := []*Data{}
	if fk != "" {
		for _, v := range g.Props.Data {
			if !strings.Contains(strings.ToLower(v.Name), fk) && !strings.Contains(strconv.Itoa(v.Power), fk) {
				continue
			}
			data = append(data, v)
		}
	} else {
		data = g.Props.Data
	}

	// Sort
	if g.Data.SortKey == "power" {
		if g.Data.SortOrders[g.Data.SortKey] == 1 {
			sort.Slice(data, func(i, j int) bool { return data[i].Power < data[j].Power })
		} else {
			sort.Slice(data, func(i, j int) bool { return data[i].Power > data[j].Power })
		}
	} else {
		if g.Data.SortOrders[g.Data.SortKey] == 1 {
			sort.Slice(data, func(i, j int) bool { return strings.Compare(data[i].Name, data[j].Name) >= 0 })
		} else {
			sort.Slice(data, func(i, j int) bool { return strings.Compare(data[i].Name, data[j].Name) < 0 })
		}
	}
	return data
}

func InitGrid() {
	o := vue.NewOption()
	o.Template = vueutil.GetTemplateById("grid-template")
	o.AddProp("data", "columns", "filterKey")
	o.Data = NewGridData
	o.OnLifeCycleEvent(vue.EvtCreated, func(vm *vue.ViewModel) {
		g := NewGrid(vm)
		g.InitGridData()
		g.SyncViewModel(vm)
	})
	o.AddComputed("filteredData", func(vm *vue.ViewModel) interface{} {
		g := NewGrid(vm)
		return g.FilteredData()
	})
	o.AddMethod("SortBy", func(vm *vue.ViewModel, args []*js.Object) {
		grid := NewGrid(vm)
		grid.SortBy(args[0])
		grid.SyncViewModel(vm)
	})
	o.NewComponent().Register("demo-grid")
}
