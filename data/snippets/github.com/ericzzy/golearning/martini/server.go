package main

import "github.com/go-martini/martini"

func main(){
  m := martini.Classic()
  m.Get("/", func() string {
    return `{"message": "hello", "status": 200}`
  })
  m.Run()
}
