package conditionalhttp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Response struct {
	Content      []byte
	LastModified time.Time
}

func ConditionalGet(url string, lastModified time.Time) (bool, Response, error) {
	modified, newLastModified, err := hasChangedSince(url, lastModified)

	if modified && err == nil {
		res, err := http.Get(url)
		if err != nil {
			return false, Response{}, err
		}

		defer res.Body.Close()

		if res.StatusCode != 200 {
			return false, Response{}, fmt.Errorf("GET request response %d", res.StatusCode)
		}

		buff, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return false, Response{}, err
		}

		return true, Response{buff, newLastModified}, nil
	} else {
		return false, Response{}, err
	}
}

func hasChangedSince(url string, lastModified time.Time) (bool, time.Time, error) {
	res, err := http.Head(url)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("HEAD request failed (%s)", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return false, time.Time{}, fmt.Errorf("HEAD request response %d", res.StatusCode)
	}

	header := res.Header.Get("Last-Modified")

	t, err := time.Parse(time.RFC1123, header)
	if err != nil {
		return false, time.Time{}, fmt.Errorf("Failed to parse Last-Modified header (%s)", err)
	}

	if t.After(lastModified) {
		return true, t, nil
	}

	return false, t, nil
}
