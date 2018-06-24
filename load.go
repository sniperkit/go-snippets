package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type unspecifiedWidget []struct {
	Type  string
	Name  string
	Value interface{}
}

func loadJsonData(url string) (unspecifiedWidget, error) {
	errorMsg := fmt.Sprintf("Can't reach provided URL: %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return unspecifiedWidget{}, errors.New(errorMsg)
	}

	req.Header.Set("User-Agent", fmt.Sprintf("dashco %s", version))

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return unspecifiedWidget{}, errors.New(errorMsg)
	}

	buffer, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return unspecifiedWidget{}, errors.New(errorMsg)
	}

	widgets := unspecifiedWidget{}
	err = json.Unmarshal(buffer, &widgets)

	if err != nil {
		return unspecifiedWidget{}, errors.New("Crash Boom Death! Invalid JSON response, data can't be parsed :skull:")
	}

	return widgets, nil
}
