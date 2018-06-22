package webservice

import (
	"appengine"
	"appengine/urlfetch"

	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/lcsontos/uwatch/catalog"
	"github.com/lcsontos/uwatch/service"
	"github.com/lcsontos/uwatch/util"
	"github.com/lcsontos/uwatch/youtube"
)

var videoCatalogRegistry = make(map[service.VideoType]catalog.VideoCatalog)

var videoCatalogRegistryRWM sync.RWMutex

func GetLongVideoUrl(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	videoTypeName := vars["videoTypeName"]
	videoId := vars["videoId"]

	videoType, err := service.GetVideoTypeByName(videoTypeName)

	if apperr, isAppErr := err.(*service.InvalidVideoTypeNameError); util.HandleError(rw, req, err, apperr, isAppErr) {
		return
	}

	videoCatalog := getVideoCatalog(videoType, req)

	lengthenedVideoUrl, err := service.LongVideoUrl(videoCatalog, videoType, videoId, req)

	handledError := false

	switch apperr := err.(type) {
	case *service.UnsupportedVideoTypeError:
		handledError = util.HandleError(rw, req, err, apperr, true)
	case *catalog.NoSuchVideoError:
		handledError = util.HandleError(rw, req, err, apperr, true)
	default:
		handledError = util.HandleError(rw, req, err, apperr, false)
	}

	if handledError {
		return
	}

	json.NewEncoder(rw).Encode(*lengthenedVideoUrl)
}

func GetParseVideoUrl(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	videoUrl := vars["videoUrl"]

	log.Printf("videoUrl=%s", videoUrl)

	parsedVideoUrl, err := service.ParseVideoUrl(videoUrl)

	if apperr, isAppErr := err.(*service.InvalidVideoUrlError); util.HandleError(rw, req, err, apperr, isAppErr) {
		return
	}

	json.NewEncoder(rw).Encode(*parsedVideoUrl)
}

// func init() {
// 	var err error

// 	// TODO create factory for creating wrapper objects to
// 	// video sharing services
// 	videoCatalogRegistry[YouTube], err = youtube.New()

// 	if err != nil {
// 		panic(err)
// 	}
// }

// I needed this "hack", because app engine requires a http.Request object to
// instanciate Transport objects. Why on earth do I have to do this???
// Reference: https://cloud.google.com/appengine/docs/go/urlfetch/
func getVideoCatalog(videoType service.VideoType, req *http.Request) catalog.VideoCatalog {
	/*
		videoCatalogRegistryRWM.RLock()

		if videoCatalog, ok := videoCatalogRegistry[videoType]; ok {
			videoCatalogRegistryRWM.RUnlock()

			return videoCatalog
		}

		videoCatalogRegistryRWM.RUnlock()

		videoCatalogRegistryRWM.Lock()

		if videoCatalog, ok := videoCatalogRegistry[videoType]; ok {
			videoCatalogRegistryRWM.Unlock()

			return videoCatalog
		}
	*/
	context := appengine.NewContext(req)

	transport := &urlfetch.Transport{Context: context}

	videoCatalog, err := youtube.NewWithRoundTripper(transport)

	if err != nil {
		panic(err)
	}

	/*
		videoCatalogRegistry[videoType] = videoCatalog

		videoCatalogRegistryRWM.Unlock()
	*/

	return videoCatalog
}
