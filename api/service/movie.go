package service

import (
	"github.com/PeppyS/what-to-watch/proto"
)

type elasticsearchClient interface {
	AutocompleteMovies(s string) ([]*proto.MoviesList_Movie, error)
	BulkIndexMovies(movies []*proto.MoviesList_Movie) error
	GetAllMovies() ([]*proto.MoviesList_Movie, error)
}

// MovieService defines service structure
type MovieService struct {
	esClient elasticsearchClient
}

// NewMovie instantiates MovieService
func NewMovie(esClient elasticsearchClient) *MovieService {
	return &MovieService{esClient}
}

func (s *MovieService) Autocomplete(text string) ([]*proto.MoviesList_Movie, error) {
	return s.esClient.AutocompleteMovies(text)
}

func (s *MovieService) BulkIndex(movies []*proto.MoviesList_Movie) error {
	return s.esClient.BulkIndexMovies(movies)
}

func (s *MovieService) GetAll() ([]*proto.MoviesList_Movie, error) {
	return s.esClient.GetAllMovies()
}
