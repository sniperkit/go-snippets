package router

import (
	"faygo-demo2/ctrl"
	"faygo-demo2/ctrl/user"
	"github.com/henrylee2cn/faygo"
)

func Router(frame *faygo.Framework) {
	frame.Route(
		frame.NewGET("/", new(ctrl.Home)),
		frame.NewGroup("api",
			frame.NewGroup("users",
				frame.NewGET("/", new(user.Index)),
			),
		),
	)
}
