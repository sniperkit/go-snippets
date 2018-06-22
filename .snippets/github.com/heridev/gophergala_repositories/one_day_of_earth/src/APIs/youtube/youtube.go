package youtube

import (
	"APIs"
	"config"
	"fmt"
	"lib"
)

type Video struct {
	ID          string
	Title       string
	Description string
	Thumb       string
}

//https://www.googleapis.com/youtube/v3/search?part=snippet&key=AIzaSyD4DKmXQIeNUtqzsaKP3il8q6jK-bfa-C0&location=37.42307,-122.08427&locationRadius=100km&type=video&maxResults=50&publishedAfter=2015-01-23T00:00:00Z&order=rating
func SearchVideos(lat, lng, datetime_RFC_3339, distance string, recursive bool, next_page_token string) (videos []Video, err *lib.CError) {
	params := map[string]string{
		"part":           "snippet",
		"key":            config.YOUTUBE_API_KEY,
		"location":       fmt.Sprintf("%s,%s", lat, lng),
		"locationRadius": fmt.Sprintf("%skm", distance),
		"type":           "video",
		"maxResults":     "50",
		"minResults":     "50",
		"publishedAfter": datetime_RFC_3339,
		"order":          "rating",
	}
	if len(next_page_token) > 1 {
		params["pageToken"] = next_page_token
	}
	data_body, err := APIs.API_call(config.GET, "https://www.googleapis.com/youtube/v3/search", params)
	if err != nil {
		return
	}
	v := data_body.(map[string]interface{})
	if v["items"] == nil {
		err = &lib.CError{}
		err.SetMessage(fmt.Sprint(data_body))
		return
	}
	for _, vi := range v["items"].([]interface{}) {
		yVideo := Video{}
		video := vi.(map[string]interface{})

		vid := video["id"].(map[string]interface{})
		yVideo.ID = vid["videoId"].(string)

		snipet := video["snippet"].(map[string]interface{})
		yVideo.Title = snipet["title"].(string)
		yVideo.Description = snipet["description"].(string)

		thumb := snipet["thumbnails"].(map[string]interface{})
		th := thumb["high"].(map[string]interface{})
		yVideo.Thumb = th["url"].(string)
		videos = append(videos, yVideo)
	}

	if v["nextPageToken"] != nil && recursive {
		vs, er := SearchVideos(lat, lng, datetime_RFC_3339, distance, true, v["nextPageToken"].(string))
		if er != nil {
			return
		}
		videos = append(videos, vs...)
	}

	return
}
