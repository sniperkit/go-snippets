package service

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/lcsontos/uwatch/catalog"
	"github.com/lcsontos/uwatch/normalizer"
	"github.com/lcsontos/uwatch/store"
)

type VideoType int

const (
	YouTube VideoType = (iota)

	// Reserved for future implementation
	Vimeo
	Youku
	Rutube

	// Internal use only!
	unknown = -1
)

type InvalidVideoUrlError struct {
	VideoUrl string
}

type InvalidVideoTypeNameError struct {
	VideoTypeName string
}

type ParsedVideoUrl struct {
	VideoId   string
	VideoType VideoType
}

type LengthenedVideoUrl struct {
	ParsedVideoUrl
	Title   string
	UrlId   int64
	UrlPath string
}

type UnsupportedVideoTypeError struct {
	VideoType VideoType
}

type urlPattern struct {
	videoType VideoType
	pattern   *regexp.Regexp
}

var urlPatterns = []urlPattern{
	urlPattern{YouTube, regexp.MustCompile("http.+youtube\\.com\\/watch\\?v=(\\S+)")},
	urlPattern{YouTube, regexp.MustCompile("http.+youtu\\.be\\/(\\S+)")},
}

var videoTypesLookupMap = make(map[string]VideoType)

var videoTypesStringMap = map[VideoType]string{
	YouTube: "YouTube",
	Vimeo:   "Vimeo",
	Youku:   "Youku",
	Rutube:  "Rutube",
}

func (err *InvalidVideoUrlError) Error() string {
	return fmt.Sprintf("\"%s\" is an invalid video URL", err.VideoUrl)
}

func (err *InvalidVideoTypeNameError) Error() string {
	return fmt.Sprintf("\"%s\" is an invalid video name", err.VideoTypeName)
}

func (err *UnsupportedVideoTypeError) Error() string {
	return fmt.Sprintf("\"%s\" is an invalid video type", err.VideoType)
}

func GetVideoTypeByName(videoTypeName string) (VideoType, error) {
	if videoType, ok := videoTypesLookupMap[videoTypeName]; !ok {
		return unknown, &InvalidVideoTypeNameError{videoTypeName}
	} else {
		return videoType, nil
	}
}

func LongVideoUrl(videoCatalog catalog.VideoCatalog, videoType VideoType, videoId string, req *http.Request) (*LengthenedVideoUrl, error) {
	if videoType != YouTube {
		return nil, &UnsupportedVideoTypeError{videoType}
	}

	videoRecord, err := videoCatalog.SearchByID(videoId)

	if err != nil {
		return nil, err
	}

	normalizedTitle := normalizer.Normalize(videoRecord.Title)

	urlId, err := store.PutVideoRecord(videoRecord, req)

	if err != nil {
		return nil, err
	}

	urlPath := fmt.Sprintf("%d/%s", urlId, normalizedTitle)

	LengthenedVideoUrl := &LengthenedVideoUrl{
		ParsedVideoUrl{videoId, videoType},
		videoRecord.Title, urlId, urlPath,
	}

	return LengthenedVideoUrl, nil
}

func ParseVideoUrl(videoUrl string) (*ParsedVideoUrl, error) {
	if videoUrl == "" {
		return nil, &InvalidVideoUrlError{""}
	}

	parsedVideoUrl := &ParsedVideoUrl{VideoType: unknown}

	for _, urlPattern := range urlPatterns {
		matches := urlPattern.pattern.FindStringSubmatch(videoUrl)

		if matches != nil {
			parsedVideoUrl.VideoId = matches[1]
			parsedVideoUrl.VideoType = urlPattern.videoType

			break
		}
	}

	if parsedVideoUrl.VideoType == unknown {
		return nil, &InvalidVideoUrlError{videoUrl}
	}

	return parsedVideoUrl, nil
}

func (videoType VideoType) String() string {
	return videoTypesStringMap[videoType]
}

func (url *LengthenedVideoUrl) String() string {
	return ""
}

func init() {
	for videoType, videoTypeName := range videoTypesStringMap {
		videoTypesLookupMap[videoTypeName] = videoType
	}
}
