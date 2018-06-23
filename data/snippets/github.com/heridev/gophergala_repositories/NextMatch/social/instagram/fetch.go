package instagram

import (
	"encoding/json"
	"net/http"
	"os"
	"log"
)

var (
	at  = `INSTAGRAM_ACCESS_TOKEN` // access_token
	api = `https://api.instagram.com/v1`                    //api_url
)

func init() {
	// Allow override of access_token and api urk
	if envAt := os.Getenv(`INSTAGRAM_ACCESS_TOKEN`); len(envAt) > 0 {
		at = envAt
	}

	if envAPI := os.Getenv(`INSTAGRAM_API`); len(envAPI) > 0 {
		api = envAPI
	}
}

var cache = make(map[string]Obj)

func ByTag(tag string) (data Obj, err error) {
	var resp *http.Response
    uri := api + `/tags/` + tag + `/media/recent?access_token=` + at
    if _, ok := cache[uri]; ok {
		log.Printf("cache  for URI %s", uri)
        return cache[uri], nil
    }
	if resp, err = http.Get(uri); err != nil {
		return
	}
	log.Printf("%s for URI %s", resp.Status, uri)



	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&data); err != nil {
		return data, err
	}
    cache[uri] = data
	return data, err
}

func BuildTag(away, home string) string {
	return away + home
}
