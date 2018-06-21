package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PeppyS/what-to-watch/proto"
)

const (
	movieIndex   = "movies"
	movieMapping = "movie"
)

type ElasticsearchClient struct {
	httpClient *http.Client
	url        string
}

func NewElasticsearchClient(c *http.Client, url string) *ElasticsearchClient {
	return &ElasticsearchClient{c, url}
}

func (c *ElasticsearchClient) AutocompleteMovies(text string) ([]*proto.MoviesList_Movie, error) {
	b := []byte(`{
		"suggest": {
			"movie-suggestions" : {
				"prefix" : "` + text + `", 
				"completion" : { 
					"field" : "title.completion"
				}
			}
		}
	}`)

	resp, err := http.Post(fmt.Sprintf("http://%s/%s/_search", c.url, movieIndex), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("Failed sending request: %v", err)
	}

	var AutocompleteResults struct {
		Suggest struct {
			MovieSuggestions []struct {
				Options []struct {
					Movie *proto.MoviesList_Movie `json:"_source"`
				} `json:"options"`
				Movie *proto.MoviesList_Movie `json:"_source"`
			} `json:"movie-suggestions"`
		} `json:"suggest"`
	}

	err = json.NewDecoder(resp.Body).Decode(&AutocompleteResults)
	if err != nil {
		return nil, fmt.Errorf("Failed json decoding search results: %v", err)
	}

	var movies []*proto.MoviesList_Movie

	for _, suggestions := range AutocompleteResults.Suggest.MovieSuggestions {
		for _, option := range suggestions.Options {
			movies = append(movies, option.Movie)
		}
	}

	return movies, nil
}

func (c *ElasticsearchClient) ClearMovieIndex() error {
	b := []byte(`{
		"query": {
			"match_all" : {}
		}
	}`)

	_, err := http.Post(fmt.Sprintf("http://%s/%s/_delete_by_query", c.url, movieIndex), "application/json", bytes.NewBuffer(b))

	return err
}

func (c *ElasticsearchClient) BulkIndexMovies(movies []*proto.MoviesList_Movie) error {
	err := c.ClearMovieIndex()
	if err != nil {
		return fmt.Errorf("Problem clearing movie index: %v", err)
	}

	for _, movie := range movies {
		b, err := json.Marshal(movie)
		if err != nil {
			return fmt.Errorf("Failed to encode payload: %v", err)
		}

		_, err = c.httpClient.Post(
			fmt.Sprintf("http://%s/%s/%s", c.url, movieIndex, movieMapping),
			"application/json",
			bytes.NewBuffer(b),
		)
		if err != nil {
			return fmt.Errorf("Failed sending request: %v", err)
		}
	}

	return nil
}

func (c *ElasticsearchClient) GetAllMovies() ([]*proto.MoviesList_Movie, error) {
	query := "_search/?size=1000"

	b := []byte(`
		{
			"sort": [{
				"imdb_meta.score": {
					"nested_path": "imdb_meta",
					"order": "desc"
				},
				"rotten_tomatoes_meta.tomato_score": {
					"nested_path": "rotten_tomatoes_meta",
					"order": "desc"
				},
				"rotten_tomatoes_meta.popcorn_score": {
					"nested_path": "rotten_tomatoes_meta",
					"order": "desc"
				}
			}]
		}
	`)

	resp, err := c.httpClient.Post(
		fmt.Sprintf("http://%s/%s/%s", c.url, movieIndex, query),
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, fmt.Errorf("Failed sending request: %v", err)
	}

	var MovieSearchResults struct {
		Hits struct {
			Hits []struct {
				Movie *proto.MoviesList_Movie `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(resp.Body).Decode(&MovieSearchResults)
	if err != nil {
		return nil, fmt.Errorf("Failed json decoding search results: %v", err)
	}

	var movies []*proto.MoviesList_Movie

	for _, movieHit := range MovieSearchResults.Hits.Hits {
		movies = append(movies, movieHit.Movie)
	}

	return movies, nil
}
