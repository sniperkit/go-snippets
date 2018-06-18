package main

import (
	"fmt"

	"github.com/Masterminds/cookoo"
)

func HelloWorld(ctx cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	ctx.Log("hello-world:", "maybe this will show")
	return true, nil
}

func main() {
	fmt.Println("chain of command...")

	// Build a new cookoo
	registry, router, context := cookoo.Cookoo()

	// Add cookoo routes
	registry.AddRoutes(
		cookoo.Route{
			Name: "TEST",
			Help: "A test route",
			Does: cookoo.Tasks{
				cookoo.Cmd{
					Name: "hi",
					Fn:   HelloWorld,
				},
			},
		},
	)

	// Execute the route
	router.HandleRequest("TEST", context, false)
}
