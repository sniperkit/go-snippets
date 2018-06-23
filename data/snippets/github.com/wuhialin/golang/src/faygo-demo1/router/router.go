package router

import (
	"github.com/henrylee2cn/faygo"
	"faygo-demo1/ctrl"
	"faygo-demo1/ctrl/warehouse"
	"faygo-demo1/middleware"
)

func Router(frame *faygo.Framework) {
	frame.Route(
		frame.NewGET("/", new(ctrl.Home)),
		frame.NewGET("/login", new(ctrl.Login)),
		frame.NewPOST("/login", new(ctrl.LoginPost)),
		frame.NewGroup("warehouse",
			frame.NewGET("/", new(warehouse.Index)),
		),
	)
	frame.Use(new(middleware.Common)).Use(new(middleware.CheckUser))
}
