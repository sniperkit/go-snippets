package main

import (
  "fmt"
  "github.com/codegangsta/controller"
  "github.com/codegangsta/negroni"
  "github.com/gophergala/gotoolbox/controllers"
  "github.com/gorilla/mux"
  "github.com/markbates/goth"
  "github.com/markbates/goth/gothic"
  "github.com/markbates/goth/providers/github"
  "os"
)

func main() {
  pwd, err := os.Getwd()
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  //fmt.Println("Working directory:")
  fmt.Println(pwd)

  goth.UseProviders(
    github.New(
      os.Getenv("GOTOOLBOX_GITHUB_KEY"),
      os.Getenv("GOTOOLBOX_GITHUB_SECRET"),
      os.Getenv("GOTOOLBOX_GITHUB_CALLBACK")),
  )
  gothic.GetProviderName = controllers.GetProviderNameForMux

  r := mux.NewRouter()
  r.Handle("/", controller.Action((*controllers.ApplicationController).Index))
  r.Handle("/auth/{provider}/callback", controller.Action((*controllers.AuthController).Create))
  r.Handle("/auth/logout", controller.Action((*controllers.AuthController).Delete))
  r.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler)

  // links routes
  r.Handle("/projects/new", controller.Action((*controllers.ProjectsController).New))
  projects := r.PathPrefix("/projects/").Subrouter()
  projects.Methods("POST").Handler(controller.Action((*controllers.ProjectsController).Create))

  // @TODO handle that with the DELETE verb
  r.Handle("/projects/{id}/delete", controller.Action((*controllers.ProjectsController).Delete))

  n := negroni.Classic()
  n.UseHandler(r)
  n.Run(":8080")
}
