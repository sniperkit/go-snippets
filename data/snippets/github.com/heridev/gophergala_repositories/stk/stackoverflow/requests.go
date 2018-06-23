package stackoverflow

import (
	"bytes"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type SearchRequest struct {
	SiteID   string
	Sort     string
	Query    string
	Tags     []string
	Accepted bool
}

type AnswerRequest struct {
	AnswerIDS []int
	SiteID    string
}

type apiRequest struct {
	what   string
	ids    string
	params *url.Values
}

func makeURL(request *apiRequest) string {
	var buf *bytes.Buffer
	buf = bytes.NewBufferString(APIBase)
	buf.WriteByte('/')
	buf.WriteString(request.what)

	if len(request.ids) > 0 {
		buf.WriteByte('/')
		buf.WriteString(request.ids)
		buf.WriteByte('/')

	}

	// add stackoverflow key
	// according to StackoverFlow API
	// This is not considered a secret, and may be safely embed in client side code or distributed binaries.
	request.params.Set("key", "vfuajcl*3rqAKABKGqWsGA((")

	buf.WriteByte('?')
	buf.WriteString(request.params.Encode())

	log.Println(buf.String())

	return buf.String()
}

func makeSearchRequest(request *SearchRequest) string {
	tagged := strings.Join(request.Tags, ";")

	args := url.Values{}
	args.Set("order", "desc")
	args.Set("sort", request.Sort)
	args.Set("site", request.SiteID)
	args.Set("title", request.Query)
	args.Set("tagged", tagged)

	return makeURL(&apiRequest{
		what:   "similar",
		params: &args,
	})
}

func makeAnswerRequest(request *AnswerRequest) string {
	var ids = make([]string, len(request.AnswerIDS))

	for i, val := range request.AnswerIDS {
		ids[i] = strconv.Itoa(val)
	}

	escaped := url.QueryEscape(strings.Join(ids, ";"))

	args := url.Values{}
	args.Set("order", "desc")
	args.Set("filter", "withbody")
	args.Set("site", request.SiteID)

	return makeURL(&apiRequest{
		what:   "answers",
		ids:    escaped,
		params: &args,
	})
}
