package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func CreateGist(desc, content string) ([]byte, error) {

	g := Gist{
		Description: desc,
		Public:      true,
		Files:       make(map[string]File),
	}
	g.Files["src.go"] = File{Content: content}

	json, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/gists", bytes.NewBuffer(json))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func ParseResponse(data []byte) (map[string]interface{}, error) {
	r := make(map[string]interface{})

	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return r, nil
}
