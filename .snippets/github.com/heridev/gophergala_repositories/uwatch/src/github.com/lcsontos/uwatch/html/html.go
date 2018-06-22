package html

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/lcsontos/uwatch/store"
	"github.com/lcsontos/uwatch/util"
)

type invalidLongUrlError struct {
	longUrl string
}

var page *template.Template
var pattern *regexp.Regexp

func (err *invalidLongUrlError) Error() string {
	return fmt.Sprintf("\"%s\" is an invalid long video URL", err.longUrl)
}

func ProcessTemplate(rw http.ResponseWriter, req *http.Request) {
	tc := make(map[string]interface{})

	urlId, _, err := getPathTokens(req.URL)

	if util.HandleError(rw, req, err, err, true) {
		return
	}

	tc["URL"], err = getVideoUrl(urlId, req)

	if util.HandleError(rw, req, err, err, true) {
		return
	}

	if err := page.Execute(rw, tc); err != nil {
		util.HandleError(rw, req, err, nil, false)
	}
}

func getPathTokens(url *url.URL) (int64, string, error) {
	matches := pattern.FindStringSubmatch(url.Path)

	if matches == nil {
		return -1, "", &invalidLongUrlError{url.RequestURI()}
	}

	urlIdString, normalizedTitle := matches[1], matches[2]

	urlId, err := strconv.ParseInt(urlIdString, 10, 64)

	if err != nil {
		return -1, "", err
	}

	return urlId, normalizedTitle, nil
}

func getVideoUrl(urlId int64, req *http.Request) (string, error) {
	videoRecord, err := store.GetVideoRecord(urlId, req)

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoRecord.VideoId)

	return url, nil
}

func init() {
	page = template.Must(template.ParseFiles("templates/html/index.html"))
	pattern = regexp.MustCompile("/(\\S+)/(\\S+)/?")
}
