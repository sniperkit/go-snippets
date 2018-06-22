package APIs

import (
	"config"
	"encoding/json"
	"io/ioutil"
	"lib"
	"net/http"
)

type Image struct {
	Width  int
	Height int
	Url    string
}

type User struct {
	UserName     string
	FullName     string
	ID           string
	ProfilePhoto string //URl to profile photo
}

type ApiImage struct {
	User      User
	Id        string
	Instalink string //Link to Instagram for that image
	Standard  Image
	Thumb     Image
	Small     Image
}

//Private function only for this package, just for making GET REST API calls
//TODO: Need to Add Pagination for all API calls

func API_call(method int, api_url string, params map[string]string) (ret_data interface{}, cerr *lib.CError) {
	var (
		resp *http.Response
		err  error
	)
	cerr = nil
	switch method {
	case config.GET:
		{
			resp, err = http.Get(lib.GenerateURL(api_url, params))
		}
	case config.POST:
		{
			resp, err = http.PostForm(api_url, lib.UrlValues(params))
		}
	default:
		cerr = &lib.CError{}
		cerr.SetMessage("Unknown Request type")
		return
	}

	if err != nil {
		cerr = &lib.CError{}
		cerr.SetMessage(err.Error())
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		cerr = &lib.CError{}
		cerr.SetMessage(err.Error())
		return
	}
	json.Unmarshal(body, &ret_data)
	return
}
