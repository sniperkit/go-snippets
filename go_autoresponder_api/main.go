package main

import (
  "github.com/codegangsta/negroni"
  "github.com/heridev/go_autoresponder_api/database"
  "github.com/heridev/go_autoresponder_api/routes"
)

func main() {
  database.InitDb()

  router := routes.Create()

  n := negroni.New()
  n.UseHandler(router)
  n.Run(":" + "3000")
}
