package flickr

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type FlickrResponse struct {
	Sizes struct {
		Size []json.RawMessage `json:"size"`
	} `json:"sizes"`
}

type ImageInfoType1 struct {
	Width  int    `json:"width,string"`
	Height int    `json:"height,string"`
	Link   string `json:"source"`
}

type ImageInfoType2 struct {
	Width  int    `json:"width"`
	Height int    `json:"height,string"`
	Link   string `json:"source"`
}

type ImageInfoType3 struct {
	Width  int    `json:"width,string"`
	Height int    `json:"height"`
	Link   string `json:"source"`
}

type ImageInfo struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Link   string `json:"source"`
}

func GetImageInfo(c context.Context, id string) (*ImageInfo, error) {
	client := urlfetch.Client(c)
	url := "https://api.flickr.com/services/rest/?method=flickr.photos.getSizes&api_key=" + API_KEY + "&photo_id=" + id + "&format=json&nojsoncallback=1"
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	var parsed FlickrResponse
	err = decoder.Decode(&parsed)
	// since flickr api returns some string widths and some int widths :((((( NOTE TO SELF: NEVER DESIGN AN API LIKE THIS
	if err != nil {
		return nil, err
	}
	if len(parsed.Sizes.Size) == 0 {
		return nil, errors.New("Flickr API returned no sizes for image with id: " + id)
	}
	// need to loop through and convert all to same type
	var infos []ImageInfo
	for _, element := range parsed.Sizes.Size {
		var info *ImageInfo = &ImageInfo{}
		if err := json.Unmarshal(element, info); err != nil {
			// try another way
			var info2 ImageInfoType1
			if err = json.Unmarshal(element, &info2); err == nil {
				log.Infof(c, "Successfully rescued a flickr json on attempt 2")
				info = &ImageInfo{Width: info2.Width, Height: info2.Height, Link: info2.Link}
			} else {
				// try another way
				var info3 ImageInfoType2
				if err = json.Unmarshal(element, &info3); err == nil {
					log.Infof(c, "Successfully rescued a flickr json on attempt 3")
					info = &ImageInfo{Width: info3.Width, Height: info3.Height, Link: info3.Link}
				} else {
					// try another way
					var info4 ImageInfoType3
					if err = json.Unmarshal(element, &info4); err == nil {
						log.Infof(c, "Successfully rescued a flickr json on attempt 4")
						info = &ImageInfo{Width: info4.Width, Height: info4.Height, Link: info4.Link}
					}
				}
			}
		}
		if info != nil {
			infos = append(infos, *info)
		}
	}
	// need to loop through and find the max
	max_width := infos[0].Width
	max_img_index := 0
	for idx, element := range infos {
		if element.Width > max_width {
			max_width = element.Width
			max_img_index = idx
		}
	}
	return &infos[max_img_index], nil
}
