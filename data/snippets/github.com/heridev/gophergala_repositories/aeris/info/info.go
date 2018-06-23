package info

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/gophergala/aeris/format"
)

type Info struct {
	Id                  string
	streams             []*Stream
	playerJsUrl         string
	decryptedSignatures bool
}

type Stream struct {
	Url       string
	signature string
	Format    *format.YoutubeFormat
}

const WATCH_PAGE_URL = "http://www.youtube.com/watch?v="

var configRegex = regexp.MustCompile(`ytplayer\.config = (.*);ytplayer\.load`)

func NewInfo(id string) *Info {
	return &Info{
		Id:                  id,
		decryptedSignatures: false,
	}
}

func (i *Info) Streams() []*Stream {
	return i.streams
}

func (i *Info) Fetch() error {

	res, err := http.Get(WATCH_PAGE_URL + i.Id)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	res.Body.Close()

	match := configRegex.FindSubmatch(buf.Bytes())
	if match == nil {
		return errors.New("Could not match yt player config in player page")
	}

	var config = struct {
		Args struct {
			UrlEncodedFmtStreamMap string `json:"url_encoded_fmt_stream_map"`
		}
		Assets struct {
			Js string
		}
	}{}

	err = json.Unmarshal(match[1], &config)
	if err != nil {
		return err
	}

	i.playerJsUrl = "http:" + config.Assets.Js

	err = i.parseStreams(config.Args.UrlEncodedFmtStreamMap)
	if err != nil {
		return err
	}

	return nil
}

func (i *Info) parseStreams(streams string) error {
	i.streams = nil
	i.decryptedSignatures = false

	formats := format.YoutubeFormats()

	for _, encodedStream := range strings.Split(streams, ",") {
		streamInfo, err := url.ParseQuery(encodedStream)
		if err != nil {
			return err
		}

		itag := streamInfo.Get("itag")
		if format, ok := formats[itag]; ok {
			stream := &Stream{
				Url:       streamInfo.Get("url"),
				signature: streamInfo.Get("s"),
				Format:    format,
			}

			signature := streamInfo.Get("s")
			if signature != "" {
				err = stream.buildSignatureUrl(signature)
				if err != nil {
					return err
				}
			}

			i.streams = append(i.streams, stream)
		}
	}

	return nil
}

func (s *Stream) buildSignatureUrl(sig string) error {
	u, err := url.Parse(s.Url)
	if err != nil {
		return err
	}

	q := u.Query()

	q.Set("signature", sig)

	u.RawQuery = q.Encode()

	s.Url = u.String()

	s.signature = sig

	return nil
}
