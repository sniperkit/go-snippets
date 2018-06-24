package common

import (
	"time"

	pb "github.com/keelerh/omniscience/protos"
)

// DocumentFetcher implementations fetch documents modified since a set time and returns them
type DocumentFetcher interface {
	Fetch(modifiedSince time.Time) ([]*pb.Document, error)
}

// DocumentFetcherFactory are methods which create instances of the DocumentFetcher interface
type DocumentFetcherFactory func() (DocumentFetcher, error)
