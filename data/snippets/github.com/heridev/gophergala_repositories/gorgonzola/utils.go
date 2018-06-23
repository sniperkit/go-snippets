package gorgonzola

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"appengine"
	"appengine/urlfetch"
)

func getJSONJobsDoc(c appengine.Context, url string) ([]byte, error) {
	client := urlfetch.Client(c)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error fetching file. Server returned status: %d %s", resp.StatusCode, resp.Status)
	}
	return data, nil
}

func validateURL(rawurl string) error {
	if _, err := url.Parse(rawurl); err != nil {
		return err
	}
	return nil
}
