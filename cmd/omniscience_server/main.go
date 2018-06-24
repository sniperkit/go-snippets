package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/keelerh/omniscience/cmd/omniscience_server/ingestion"
	"github.com/keelerh/omniscience/cmd/omniscience_server/search"
	pb "github.com/keelerh/omniscience/protos"
	"github.com/olivere/elastic"
	"google.golang.org/grpc"
)

var (
	fElasticSearchServerURL = flag.String(
		"elastic_search_server_url",
		"http://elasticsearch:9200",
		"URL on which to access elastic search cluster for indexing and searching",
	)
	fGRPCServerPort = flag.Int(
		"grpc_server_port",
		50051,
		"Port on which gRPC server should listen",
	)
	fHTTPServerPort = flag.Int(
		"http_server_port",
		80,
		"Port on which the HTTP server should listen",
	)
)

func runGRPCServer(grpcListener net.Listener) {
	// Obtain a client and connect to the default Elasticsearch installation on 127.0.0.1:9200.
	esClient, err := elastic.NewClient(elastic.SetURL(*fElasticSearchServerURL))
	if err != nil {
		log.Fatalf("failed to initialise Elasticsearch client: %v", err)
	}

	ingester := ingestion.NewIngester(esClient)
	if err != nil {
		log.Fatalf("failed to initialise Ingestion service: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterIngesterServer(s, ingester)

	searchSvc := search.NewSearchService(esClient)
	pb.RegisterSearchServer(s, searchSvc)

	if err := s.Serve(grpcListener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runHTTPServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterSearchHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("localhost:%d", *fGRPCServerPort),
		opts,
	)
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", *fHTTPServerPort), mux)
}

func main() {
	flag.Parse()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *fGRPCServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go runGRPCServer(grpcListener)

	if err := runHTTPServer(); err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}
