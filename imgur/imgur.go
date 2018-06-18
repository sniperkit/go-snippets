package imgur

import (
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

type ImgurResponse struct {
	Data ImageInfo `json:"data"`
}

type ImageInfo struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Link   string `json:"link"`
}

func GetImageInfo(c context.Context, id string) (*ImageInfo, error) {
	client := urlfetch.Client(c)
	url := "https://api.imgur.com/3/image/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Client-ID "+CLIENT_ID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var parsed ImgurResponse
	err = decoder.Decode(&parsed)
	if err != nil {
		return nil, err
	}
	return &parsed.Data, nil
}
