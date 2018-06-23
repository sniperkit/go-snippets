package main

import (
	"github.com/ant0ine/go-json-rest/rest"
	eram "github.com/Onefootball/entity-rest-api/manager"
	era "github.com/Onefootball/entity-rest-api/api"
	"log"
	"net/http"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// serve ng-admin angular app
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/admin/", http.StripPrefix("/admin/", fs))

	db, err := sql.Open("sqlite3", "./blog.db")

	defer db.Close()

	if err != nil {
		panic(err)
	}

	api := rest.NewApi()
	api.Use(rest.DefaultProdStack...)

	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			return true
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Total-Count", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	// initialize Entity REST API
	entityManager := eram.NewEntityDbManager(db)
	entityRestApi := era.NewEntityRestAPI(entityManager)

	router, err := rest.MakeRouter(
		rest.Get("/api/:entity", entityRestApi.GetAllEntities),
		rest.Post("/api/:entity", entityRestApi.PostEntity),
		rest.Get("/api/:entity/:id", entityRestApi.GetEntity),
		rest.Put("/api/:entity/:id", entityRestApi.PutEntity),
		rest.Delete("/api/:entity/:id", entityRestApi.DeleteEntity),
	)

	if err != nil {
		log.Fatal(err)
	}

	api.SetApp(router)

	http.Handle("/api/", api.MakeHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
