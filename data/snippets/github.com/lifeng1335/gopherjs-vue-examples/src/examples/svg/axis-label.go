package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"math"
	"vueutil"
)

type Stat struct {
	*js.Object
	Label string  `js:"label"`
	Value float64 `js:"value"`
}

func NewStat(label string, value float64) *Stat {
	s := &Stat{
		Object: js.Global.Get("Object").New(),
	}
	s.Label = label
	s.Value = value
	return s
}

type AxisLabelData struct {
	*js.Object
}

type AxisLabelProps struct {
	*js.Object
	Stat  Stat    `js:"stat"`
	Index int     `js:"index"`
	Total float64 `js:"total"`
}

type AxisLabel struct {
	Data  *AxisLabelData
	Props *AxisLabelProps
}

func NewAxisLabel(vm *vue.ViewModel) *AxisLabel {
	return &AxisLabel{
		Data: &AxisLabelData{
			Object: vm.Data,
		},
		Props: &AxisLabelProps{
			Object: vueutil.PropsData(vm),
		},
	}
}

func InitAxisLabel() *vue.Component {
	o := vue.NewOption()
	o.Name = "axis-label"
	o.Template = vueutil.GetTemplateById("axis-label-template")
	o.AddProp("stat", "index", "total")
	o.Data = vueutil.EmptyDataFunc
	o.AddComputed("point", func(vm *vue.ViewModel) interface{} {
		axisLabel := NewAxisLabel(vm)
		return ValueToPoint(
			axisLabel.Props.Stat.Value+float64(10),
			axisLabel.Props.Index,
			axisLabel.Props.Total)
	})
	return o.NewComponent()
}

type Point struct {
	*js.Object
	X float64 `js:"x"`
	Y float64 `js:"y"`
}

// math helper
func ValueToPoint(value float64, index int, total float64) *Point {
	x := float64(0)
	y := -0.8 * value
	angle := math.Pi * 2 / total * float64(index)
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	tx := x*cos - y*sin + float64(100)
	ty := x*sin + y*cos + float64(100)

	point := &Point{
		Object: js.Global.Get("Object").New(),
	}
	point.X = tx
	point.Y = ty
	return point
}
