package common

import (
	"net/rpc"
	"strings"
	"net/url"
)

// A Conn struct contains the hostPort of connection endpoint and the dialer
// to this hostport
type Conn struct {
	HostPort []string
	Dialer   *rpc.Client
}

// removeHttp returns a raw url with no http:// or https:// as prefix
func RemoveHttp(url string) string {
	//Remove http:// and https:// from the absolute url
	removeHttp := strings.Replace(url, "http://", "", -1)
	removeHttps := strings.Replace(removeHttp, "https://", "", -1)
	return removeHttps
}

// GetAbsoluteUrl returns the absolute url based on href contained in the webpage
// and the base url information. (e.g. it turns a /language_tool href link in
// http://www.google.com into http://www.google.com/language_tool)
func GetAbsoluteUrl(href, base string) string {
	uri, err := url.Parse(href)
	//abandon broken href
	if err != nil {
		return ""
	}
	//abandon broken base url
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	//resolve and generate absolute url
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}
