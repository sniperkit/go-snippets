package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"strconv"
	"vueutil"
)

type DraggableHeaderData struct {
	*js.Object
	Dragging bool   `js:"dragging"`
	C        *Point `js:"c"`
	Start    *Point `js:"start"`
}

func NewDraggableHeaderData() interface{} {
	dh := &DraggableHeaderData{
		Object: js.Global.Get("Object").New(),
	}
	dh.Dragging = false
	dh.C = NewPoint(160.0, 160.0)
	dh.Start = NewPoint(0.0, 0.0)
	return dh
}

type DraggableHeaderProps struct {
	*js.Object
}

type DraggableHeader struct {
	Data  *DraggableHeaderData
	Props *DraggableHeaderProps
}

func (dh *DraggableHeader) HeaderPath() string {
	return `M0,0 L320,0 320,160` +
		`Q` + strconv.FormatFloat(dh.Data.C.X, 'f', 2, 64) + `,` +
		strconv.FormatFloat(dh.Data.C.Y, 'f', 2, 64) +
		` 0,160`
}

func (dh *DraggableHeader) ContentPosition() interface{} {
	dy := dh.Data.C.Y - 160.0
	dampen := 0.0
	if dy > 0 {
		dampen = 2.0
	} else {
		dampen = 4.0
	}
	transform := `translate3d(0,` + strconv.FormatFloat(dy/dampen, 'f', 2, 64) + `px,0)`
	obj := js.Global.Get("Object").New()
	obj.Set("transform", transform)

	return obj
}

func (dh *DraggableHeader) StartDrag(event *js.Object) {
	e := js.Undefined
	if event.Get("changedTouches") != js.Undefined {
		e = event.Get("changedTouches").Index(0)
	} else {
		e = event
	}
	dh.Data.Dragging = true
	dh.Data.Start.X = e.Get("pageX").Float()
	dh.Data.Start.Y = e.Get("pageY").Float()
}

func (dh *DraggableHeader) OnDrag(event *js.Object) {
	e := js.Undefined
	if event.Get("changedTouches") != js.Undefined {
		e = event.Get("changedTouches").Index(0)
	} else {
		e = event
	}
	if dh.Data.Dragging {
		dh.Data.C.X = 160.0 + (e.Get("pageX").Float() - dh.Data.Start.X)
		dy := e.Get("pageY").Float() - dh.Data.Start.Y
		dampen := 4.0
		if dy > 0 {
			dampen = 1.5
		}
		dh.Data.C.Y = 160.0 + dy/dampen
	}
}

func (dh *DraggableHeader) StopDrag() {
	if dh.Data.Dragging {
		dh.Data.Dragging = false
		dynamics := js.Global.Get("dynamics")
		//dynamics.Call("animate")
		ArgA := NewPoint(160.0, 160.0)
		ArgB := js.Global.Get("Object").New()
		ArgB.Set("type", dynamics.Get("spring"))
		ArgB.Set("duration", 700)
		ArgB.Set("friction", 280)
		dynamics.Call("animate", dh.Data.C, ArgA, ArgB)
	}
}

func (dh *DraggableHeader) SyncViewModel(vm *vue.ViewModel) {
	// [Vue warn]: Avoid replacing instance root $data. Use nested data properties instead.
	//vm.Data = t.Date.Object
	keys := js.Keys(dh.Data.Object)
	for _, v := range keys {
		vm.Data.Set(v, dh.Data.Get(v))
	}

	vm.Get("$options").Set("propsData", dh.Props.Object)
}

func NewDraggableHeader(vm *vue.ViewModel) *DraggableHeader {
	return &DraggableHeader{
		Data: &DraggableHeaderData{
			Object: vm.Data,
		},
		Props: &DraggableHeaderProps{
			Object: vueutil.PropsData(vm),
		},
	}
}

func InitDraggableHeader() {
	o := vue.NewOption()
	o.Name = "draggable-header-view"
	o.Template = vueutil.GetTemplateById("header-view-template")
	o.Data = NewDraggableHeaderData
	o.AddComputed("headerPath", func(vm *vue.ViewModel) interface{} {
		dh := NewDraggableHeader(vm)
		return dh.HeaderPath()
	})
	o.AddComputed("contentPosition", func(vm *vue.ViewModel) interface{} {
		dh := NewDraggableHeader(vm)
		return dh.ContentPosition()
	})
	o.AddMethod("startDrag", func(vm *vue.ViewModel, args []*js.Object) {
		dh := NewDraggableHeader(vm)
		dh.StartDrag(args[0])
		dh.SyncViewModel(vm)
	})
	o.AddMethod("onDrag", func(vm *vue.ViewModel, args []*js.Object) {
		dh := NewDraggableHeader(vm)
		dh.OnDrag(args[0])
		dh.SyncViewModel(vm)
	})
	o.AddMethod("stopDrag", func(vm *vue.ViewModel, args []*js.Object) {
		dh := NewDraggableHeader(vm)
		dh.StopDrag()
		dh.SyncViewModel(vm)
	})
	o.NewComponent().Register("draggable-header-view")
}

type Point struct {
	*js.Object
	X float64 `js:"x"`
	Y float64 `js:"y"`
}

func NewPoint(x float64, y float64) *Point {
	p := &Point{
		Object: js.Global.Get("Object").New(),
	}
	p.X = x
	p.Y = y
	return p
}
