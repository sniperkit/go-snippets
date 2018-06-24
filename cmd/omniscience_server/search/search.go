package search

import (
	"context"

	"github.com/keelerh/omniscience/protos"
	"github.com/olivere/elastic"
)

type searchService struct {
	elasticClient *elastic.Client
}

func NewSearchService(esClient *elastic.Client) *searchService {
	return &searchService{
		elasticClient: esClient,
	}
}

func (s *searchService) ListServices(ctx context.Context, in *omniscience.ListServicesRequest) (*omniscience.ListServicesResponse, error) {
	return &omniscience.ListServicesResponse{
		Services: []*omniscience.Service{},
	}, nil
}
