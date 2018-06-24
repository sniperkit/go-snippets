package common

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	pb "github.com/keelerh/omniscience/protos"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	fIngestionServiceAddress = flag.String(
		"ingestion_service_address",
		"localhost:50051",
		"Address at which the omniscience ingestion service can be found.")
	fModifiedSince = flag.String(
		"modified_since",
		"02 Jan 06 15:04 MST",
		"The ingester will only fetch documents modified after this date.")
)

func CreateIngesterCLI(fetcherFactory DocumentFetcherFactory) {
	flag.Parse()
	
	// TODO(ains): more sensible ctx?
	ctx := context.Background()

	fetcher, err := fetcherFactory()
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create fetcher"))
	}

	modifiedSince, err := time.Parse(time.RFC822, *fModifiedSince)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to parse modified since timestamp"))
	}

	fetchedDocuments, err := fetcher.Fetch(modifiedSince)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to fetch documents"))
	}

	cc, err := grpc.Dial(*fIngestionServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to connect to IngestionService"))
	}
	defer cc.Close()

	client := pb.NewIngesterClient(cc)
	stream, err := client.Ingest(ctx)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to open Ingest stream"))
	}

	for _, doc := range fetchedDocuments {
		if err := stream.Send(doc); err != nil {
			// TODO(ains): something better to do here than fatal(?)
			log.Fatal(errors.Wrap(err, "failed to send document"))
		}
	}

	if _, err := stream.CloseAndRecv(); err != nil {
		// We expect io.EOF once the stream has closed.
		if err != io.EOF {
			log.Fatal(errors.Wrap(err, "failed to close stream"))
		}
	}
}
