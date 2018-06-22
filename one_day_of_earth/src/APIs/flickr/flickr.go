package flickr

import (
	"APIs"
	"config"
	"fmt"
	"lib"
	"strconv"
)

//https://www.flickr.com/services/rest/?method=flickr.photos.search&api_key=790ebc6499faaab832fca792b24788db&radius=20&lat=37.7624499&lon=-122.4602593&min_upload_date=2015-01-23&max_upload_date=2015-01-24&format=json&nojsoncallback=1&per_page=500
// Flickr distance from 1 to 20 as miles
func SearchImages(lat, lng, minTimeStamp, distance string) (images []APIs.ApiImage, cerr *lib.CError) {
	cerr = nil
	data, cerr := APIs.API_call(config.GET, "https://www.flickr.com/services/rest/", map[string]string{
		"method":          "flickr.photos.search",
		"api_key":         config.FLICKR_API_KEY,
		"lat":             lat,
		"lon":             lng,
		"min_upload_date": minTimeStamp,
		//		"max_upload_date": minTimeStamp + 86000, // Flickr not filtering by time if there only min time
		"radius":         distance,
		"format":         "json",
		"nojsoncallback": "1",
		"per_page":       "500",
	})
	if cerr != nil {
		return
	}
	i := data.(map[string]interface{})
	//Checking if there fail in request
	if i["stat"].(string) != "ok" {
		cerr = &lib.CError{}
		cerr.SetMessage(i["message"].(string))
		return
	}

	p := i["photos"].(map[string]interface{})

	for _, img_interface := range p["photo"].([]interface{}) {
		flickr_image := APIs.ApiImage{}
		img := img_interface.(map[string]interface{})
		flickr_image.Id = img["id"].(string)
		flickr_image.Instalink = "" //don't need for Flickr
		//Getting image URLS using Flickr URL API https://www.flickr.com/services/api/misc.urls.html
		//https://farm{farm-id}.staticflickr.com/{server-id}/{id}_{secret}_[mstzb].jpg
		flickr_image.Standard = APIs.Image{
			Url:    fmt.Sprintf("https://farm%s.staticflickr.com/%s/%s_%s_%s.jpg", strconv.Itoa(int(img["farm"].(float64))), img["server"].(string), img["id"].(string), img["secret"].(string), "z"),
			Width:  640,
			Height: 640,
		}

		flickr_image.Thumb = APIs.Image{
			Url:    fmt.Sprintf("https://farm%s.staticflickr.com/%s/%s_%s_%s.jpg", strconv.Itoa(int(img["farm"].(float64))), img["server"].(string), img["id"].(string), img["secret"].(string), "t"),
			Width:  -1, //Size is not fixed
			Height: -1, //Size is not fixed
		}
		flickr_image.Small = APIs.Image{
			Url:    fmt.Sprintf("https://farm%s.staticflickr.com/%s/%s_%s_%s.jpg", strconv.Itoa(int(img["farm"].(float64))), img["server"].(string), img["id"].(string), img["secret"].(string), "m"),
			Width:  -1, //Size is not fixed
			Height: -1, //Size is not fixed
		}
		//For Flickr need to make API call to
		// https://www.flickr.com/services/rest/?method=flickr.people.getInfo&api_key=790ebc6499faaab832fca792b24788db&user_id=37166181@N00
		// For Getting user information with img["owner"] parameter as user_id
		// in this version will be without user info

		flickr_image.User = APIs.User{
			UserName:     "",
			FullName:     "",
			ProfilePhoto: "",
			ID:           img["owner"].(string),
		}
		images = append(images, flickr_image)
	}

	return
}
