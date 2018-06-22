package main

import (
		"gopkg.in/kataras/iris.v6"
		"gopkg.in/kataras/iris.v6/adaptors/httprouter" // <--- this line
	)

func main(){
  app := iris.New()
  app.Get("/hi", hi)
  // right below the iris.New()
  app.Adapt(httprouter.New())
  app.Listen("0.0.0.0:8080")
}

func hi(ctx *iris.Context){
   ctx.Writef("Hi %s", "iris")
}
