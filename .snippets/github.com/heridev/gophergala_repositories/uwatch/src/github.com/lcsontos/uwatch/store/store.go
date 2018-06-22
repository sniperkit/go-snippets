package store

import (
	"appengine"
	"appengine/datastore"

	"net/http"

	"github.com/lcsontos/uwatch/catalog"
)

const _KIND = "VideoRecord"

func GetVideoRecord(id int64, req *http.Request) (*catalog.VideoRecord, error) {
	context := appengine.NewContext(req)

	key := datastore.NewKey(context, _KIND, "", id, nil)

	var videoRecord catalog.VideoRecord

	err := datastore.Get(context, key, &videoRecord)

	if err != nil {
		return nil, err
	}

	return &videoRecord, nil
}

func PutVideoRecord(videoRecord *catalog.VideoRecord, req *http.Request) (int64, error) {
	context := appengine.NewContext(req)

	key := datastore.NewIncompleteKey(context, _KIND, nil)

	// Avoid error: Property Description is too long. Maximum length is 500
	videoRecord.Description = ""

	key, err := datastore.Put(context, key, videoRecord)

	if err != nil {
		return -1, err
	}

	return key.IntID(), nil
}
