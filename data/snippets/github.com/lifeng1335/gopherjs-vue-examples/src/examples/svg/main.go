package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"vueutil"
)

type SvgData struct {
	*js.Object
	NewLabel string  `js:"newLabel"`
	Stats    []*Stat `js:"stats"`
}

func NewSvgData() interface{} {
	stats := []*Stat{}
	stats = append(stats,
		NewStat("A", 100),
		NewStat("B", 100),
		NewStat("C", 100),
		NewStat("D", 100),
		NewStat("E", 100),
		NewStat("F", 100),
	)

	svg := &SvgData{
		Object: js.Global.Get("Object").New(),
	}
	svg.NewLabel = ""
	svg.Stats = stats
	return svg
}

type SvgProps struct {
	*js.Object
}

type Svg struct {
	Data  *SvgData
	Props *SvgProps
}

func NewSvg(vm *vue.ViewModel) *Svg {
	return &Svg{
		Data: &SvgData{
			Object: vm.Data,
		},
		Props: &SvgProps{
			Object: vueutil.PropsData(vm),
		},
	}
}

func (s *Svg) SyncViewModel(vm *vue.ViewModel) {
	// [Vue warn]: Avoid replacing instance root $data. Use nested data properties instead.
	//vm.Data = t.Date.Object
	keys := js.Keys(s.Data.Object)
	for _, v := range keys {
		vm.Data.Set(v, s.Data.Get(v))
	}

	vm.Get("$options").Set("propsData", s.Props.Object)
}

func (s *Svg) Add() {
	if s.Data.NewLabel == "" {
		return
	}
	s.Data.Stats = append(s.Data.Stats, NewStat(s.Data.NewLabel, float64(100)))
	s.Data.NewLabel = ""
}

func (s *Svg) Remove(args []*js.Object) {
	stat := &Stat{
		Object: args[0],
	}
	if len(s.Data.Stats) > 3 {
		index := IndexOf(s.Data.Stats, stat)
		stats := []*Stat{}
		for key, value := range s.Data.Stats {
			if key != index {
				stats = append(stats, value)
			}
		}
		s.Data.Stats = stats
	} else {
		js.Global.Call("alert", `Can't delete more!`)
	}
}

func main() {
	InitPolygraph()

	o := vue.NewOption()
	o.Data = NewSvgData()
	o.AddMethod("add", func(vm *vue.ViewModel, args []*js.Object) {
		// Use As: e.preventDefault()
		args[0].Call("preventDefault")
		svg := NewSvg(vm)
		svg.Add()
		svg.SyncViewModel(vm)
	})
	o.AddMethod("remove", func(vm *vue.ViewModel, args []*js.Object) {
		svg := NewSvg(vm)
		svg.Remove(args)
		svg.SyncViewModel(vm)
	})
	v := o.NewViewModel()
	v.Mount("#demo")
}

// Help func
func IndexOf(stats []*Stat, stat *Stat) int {
	for key, s := range stats {
		if s.Object.Get("__ob__").Get("dep").Get("id") == stat.Object.Get("__ob__").Get("dep").Get("id") {
			return key
		}
	}
	return -1
}
