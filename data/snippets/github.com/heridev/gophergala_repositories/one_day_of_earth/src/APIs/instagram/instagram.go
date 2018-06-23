package instagram

import (
	"APIs"
	"config"
	"lib"
)

func SearchImages(lat, lng, minTimeStamp, distance string) (images []APIs.ApiImage, cerr *lib.CError) {
	cerr = nil
	data, cerr := APIs.API_call(config.GET, "https://api.instagram.com/v1/media/search", map[string]string{
		"lat":           lat,
		"lng":           lng,
		"min_timestamp": minTimeStamp,
		"distance":      distance,
		"access_token":  config.INSTAGRAM_ACCESS_TOKEN,
	})
	if cerr != nil {
		return
	}
	i := data.(map[string]interface{})
	//Getting Status code for request
	s := i["meta"].(map[string]interface{})
	if int(s["code"].(float64)) != 200 {
		cerr = &lib.CError{}
		cerr.SetMessage(s["error_message"].(string))
		return
	}

	for _, img_interface := range i["data"].([]interface{}) {
		insta_image := APIs.ApiImage{}
		img := img_interface.(map[string]interface{})
		if img["type"].(string) != "image" {
			continue
		}
		insta_image.Id = img["id"].(string)
		insta_image.Instalink = img["link"].(string)
		imgs := img["images"].(map[string]interface{})
		im_standard := imgs["standard_resolution"].(map[string]interface{})
		insta_image.Standard = APIs.Image{
			Width:  int(im_standard["width"].(float64)),
			Height: int(im_standard["height"].(float64)),
			Url:    im_standard["url"].(string),
		}
		im_thumbnail := imgs["thumbnail"].(map[string]interface{})
		insta_image.Standard = APIs.Image{
			Width:  int(im_thumbnail["width"].(float64)),
			Height: int(im_thumbnail["height"].(float64)),
			Url:    im_thumbnail["url"].(string),
		}
		im_small := imgs["low_resolution"].(map[string]interface{})
		insta_image.Standard = APIs.Image{
			Width:  int(im_small["width"].(float64)),
			Height: int(im_small["height"].(float64)),
			Url:    im_small["url"].(string),
		}

		user_data := img["user"].(map[string]interface{})
		insta_image.User = APIs.User{
			UserName:     user_data["username"].(string),
			FullName:     user_data["full_name"].(string),
			ProfilePhoto: user_data["profile_picture"].(string),
			ID:           user_data["id"].(string),
		}
		images = append(images, insta_image)
	}

	return
}
