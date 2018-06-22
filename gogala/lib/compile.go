package lib

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type CompileResponse struct {
	Errors string
	Events []map[string]interface{}
}

func (c CompileResponse) Message() interface{} {

	for i := 0; i < len(c.Events); i++ {
		for k, v := range c.Events[i] {
			if k == "Message" {
				return v
			}
		}
	}
	return nil

}

func ParseCompileResponse(data []byte) (CompileResponse, error) {
	var cr CompileResponse
	err := json.Unmarshal(data, &cr)
	if err != nil {
		return cr, err
	}
	return cr, nil
}

func Compile(body string) ([]byte, error) {

	resp, err := http.PostForm("http://golang.org/compile",
		url.Values{"version": {"2"}, "body": {body}})

	resp.Header.Add("User-Agent", "Gopher-Gala-2015@julienc")

	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return data, nil
}
