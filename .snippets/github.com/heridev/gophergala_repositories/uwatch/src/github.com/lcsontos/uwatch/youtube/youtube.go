package youtube

import (
	"net/http"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"

	"github.com/lcsontos/uwatch/catalog"
)

const _DEVELOPER_KEY = "AIzaSyDdkE9YERsVIFgH-l7mTxpBgHLDkmkPyMA"
const _PART = "id,snippet"

type Service struct {
	httpClient     *http.Client
	youTubeService *youtube.Service
}

func New() (*Service, error) {
	service, err := NewWithRoundTripper(nil)

	return service, err
}

func NewWithRoundTripper(roundTripper http.RoundTripper) (*Service, error) {
	transport := &transport.APIKey{
		Key:       _DEVELOPER_KEY,
		Transport: roundTripper,
	}

	client := &http.Client{
		Transport: transport,
	}

	service, err := youtube.New(client)

	if err != nil {
		return nil, err
	}

	return &Service{httpClient: client, youTubeService: service}, nil
}

func (service Service) SearchByID(videoId string) (*catalog.VideoRecord, error) {
	call := service.getVideosListCall(videoId)

	response, err := call.Do()

	if err != nil {
		return nil, err
	}

	if cap(response.Items) == 0 {
		return nil, &catalog.NoSuchVideoError{videoId}
	}

	videoRecord := catalog.VideoRecord{}

	videoRecord.VideoId = response.Items[0].Id
	videoRecord.Description = response.Items[0].Snippet.Description

	// TODO implement time conversion
	// videoRecord.PublishedAt = response.Items[0].Snippet.PublishedAt

	videoRecord.Title = response.Items[0].Snippet.Title

	return &videoRecord, nil
}

func (service Service) SearchByTitle(title string, maxResults int64) ([]catalog.VideoRecord, error) {
	call := service.getSearchListCall(title, maxResults)

	response, err := call.Do()

	if err != nil {
		return nil, err
	}

	items := response.Items
	itemCount := cap(items)

	var videoRecords []catalog.VideoRecord

	if itemCount == 0 {
		videoRecords = []catalog.VideoRecord{}
	} else {
		videoRecords = make([]catalog.VideoRecord, itemCount)

		for index, item := range items {
			// TODO implement time conversion

			videoRecord := catalog.NewVideoRecord(
				item.Id.VideoId, item.Snippet.Title, item.Snippet.Description, time.Now())

			videoRecords[index] = *videoRecord
		}
	}

	return videoRecords, nil
}

func (service Service) getSearchListCall(searchTerm string, maxResults int64) *youtube.SearchListCall {
	call := service.youTubeService.Search.List(_PART)

	call.Q(searchTerm)
	call.MaxResults(maxResults)

	return call
}

func (service Service) getVideosListCall(videoId string) *youtube.VideosListCall {
	call := service.youTubeService.Videos.List(_PART)

	return call.Id(videoId)
}
