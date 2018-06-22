package lib

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type CError struct {
	message string
}

func (e *CError) Message() string {
	return e.message
}

func (e *CError) SetMessage(s string) {
	e.message = s
}

func GenerateURL(api_url string, params map[string]string) string {
	var buffer bytes.Buffer //Using Buffer because it should be faster
	buffer.WriteString(api_url)
	i := strings.Index(api_url, "?")
	if i == -1 {
		buffer.WriteString("?")
	} else if i < (len(api_url) - 1) {
		buffer.WriteString("?")
	}
	buffer.WriteString(URLEncodeParams(params))
	return buffer.String()
}

func URLEncodeParams(params map[string]string) string {
	val := UrlValues(params)
	return val.Encode()
}

func UrlValues(params map[string]string) url.Values {
	val := url.Values{}
	for k, v := range params {
		val.Add(k, v)
	}
	return val
}

func YesterdayTime() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), (now.Day() - 1), 0, 0, 0, 0, time.UTC)
}

func MD5strings(data ...string) string {
	md := md5.New()
	for _, s := range data {
		io.WriteString(md, s)
	}

	return hex.EncodeToString(md.Sum(nil))
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	ret_str := strconv.FormatFloat(input_num, 'f', 6, 64)
	for {
		last := len(ret_str) - 1
		if last >= 0 && ret_str[last] == '0' {
			ret_str = ret_str[:last]
		} else {
			break
		}
	}
	return ret_str
}
