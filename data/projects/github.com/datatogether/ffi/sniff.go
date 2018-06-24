package ffi

import (
	"net/http"
)

func DetectContentType(urlstr string, body []byte) string {
	sniff := http.DetectContentType(body)
	// if fn, err := FilenameFromUrlString(urlstr); err == nil {
	// }
	return sniff
}
