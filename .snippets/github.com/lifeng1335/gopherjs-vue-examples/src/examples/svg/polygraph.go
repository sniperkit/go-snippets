package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/oskca/gopherjs-vue"
	"strconv"
	"vueutil"
)

type PolygraphData struct {
	*js.Object
}

type PolygraphProps struct {
	*js.Object
	Stats []Stat `js:"stats"`
}

type Polygraph struct {
	Data  *PolygraphData
	Props *PolygraphProps
}

func NewPolygraph(vm *vue.ViewModel) *Polygraph {
	return &Polygraph{
		Data: &PolygraphData{
			Object: vm.Data,
		},
		Props: &PolygraphProps{
			Object: vueutil.PropsData(vm),
		},
	}
}

func (p *Polygraph) Points() string {
	length := len(p.Props.Stats)
	PointsStr := ""
	for index, stat := range p.Props.Stats {
		point := ValueToPoint(stat.Value, index, float64(length))
		px := strconv.FormatFloat(point.X, 'f', 4, 64)
		py := strconv.FormatFloat(point.Y, 'f', 4, 64)
		PointsStr = PointsStr + " " + px + "," + py
	}
	return PointsStr
}

func InitPolygraph() {
	o := vue.NewOption()
	o.Name = "polygraph"
	o.Template = vueutil.GetTemplateById("polygraph-template")
	o.AddProp("stats")
	o.Data = vueutil.EmptyDataFunc
	o.AddSubComponent("axis-label", InitAxisLabel())
	o.AddComputed("points", func(vm *vue.ViewModel) interface{} {
		polygraph := NewPolygraph(vm)
		return polygraph.Points()
	})
	o.NewComponent().Register("polygraph")
}
