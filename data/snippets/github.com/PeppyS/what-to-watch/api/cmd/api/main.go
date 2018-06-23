package main

import (
	"log"
	"net/http"
	"os"

	"github.com/PeppyS/what-to-watch/api/controller"
	"github.com/PeppyS/what-to-watch/api/server"
	"github.com/PeppyS/what-to-watch/api/service"
)

func main() {
	grpcPort := ":" + os.Getenv("GRPC_PORT")
	httpPort := ":" + os.Getenv("HTTP_PORT")
	elasticSearchURL := os.Getenv("ELASTICSEARCH_URL")

	healthAPI := controller.NewHealth()
	movieAPI := controller.NewMovie(
		service.NewMovie(
			service.NewElasticsearchClient(
				http.DefaultClient,
				elasticSearchURL,
			),
		),
	)

	go server.ListenAndServe(grpcPort, movieAPI, healthAPI)

	log.Fatal(server.ListenAndServeHTTPGateway(grpcPort, httpPort))
}
