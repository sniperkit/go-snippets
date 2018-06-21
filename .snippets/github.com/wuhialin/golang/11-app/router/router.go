package router

import (
	"../ctrl"
	"../ctrl/warehouse"
	"github.com/henrylee2cn/faygo"
	"../middleware"
)

func Router(frame *faygo.Framework) {
	frame.Use(new(middleware.StartTimeRequest))
	frame.Route(
		frame.NewAPI("GET", "/", new(ctrl.Index)),
		frame.NewGroup("warehouse",
			frame.NewGET("/", new(warehouse.Index)),
		),
	)
	frame.Use(new(middleware.EndTimeRequest))
}
