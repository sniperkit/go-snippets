package hunter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gophergala/cobs/instrumenter"
)

type (
	SearchResults struct {
		NumPages   int      `json:"num_pages"`
		NumResults int      `json:"num_results"`
		PageSize   int      `json:"page_size"`
		Page       int      `json:"page"`
		Query      string   `json:"query"`
		Results    []Result `json:"results"`
	}

	Result struct {
		IsAutomated bool   `json:"is_automated"`
		Name        string `json:"name"`
		IsTrusted   bool   `json:"is_trusted"`
		IsOfficial  bool   `json:"is_official"`
		StarCount   int    `json:"star_count'`
		Description string `json:"description"`
	}
)

var rc redis.Conn

func SearchDockerRegistry(q string, params ...string) []Result {
	r := "https://index.docker.io/v1/search?q="
	if len(params) > 0 {
		r = params[0]
	}
	u := r + url.QueryEscape(q)

	res, err := http.Get(u)
	if err != nil {
		log.Fatalf("error retrieving: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading body: %s", err)
	}

	var results SearchResults
	json.Unmarshal(body, &results)

	return results.Results

}

func GetRawDockerfileFromRegistry(name string, params ...string) string {
	r := "https://registry.hub.docker.com/u/"
	if len(params) > 0 {
		r = params[0]
	}
	u := r + name + "/dockerfile/raw"

	res, err := http.Get(u)
	if err != nil {
		log.Fatalf("error retrieving: %s", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error reading body: %s", err)
	}

	if res.StatusCode != 200 {
		body = []byte("")
	}

	return string(body)
}

func GoHunting(imageId string) {
	var err error
	log.Println("Connecting to Redis")
	rc, err = redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatalf("Issues with redis: %s", err)
	}
	defer rc.Close()

	u, err := redis.String(rc.Do("GET", "url-"+imageId))
	if err != nil {
		log.Fatalf("Not found: %s", err)
	}
	results := SearchDockerRegistry(u)
	df := ""
	if len(results) > 0 {
		res := results[1]
		df = GetRawDockerfileFromRegistry(res.Name)
	}

	if len(strings.TrimSpace(df)) > 0 {
		rc.Do("SET", "dockerfile-"+imageId, df)
		go instrumenter.Run(imageId)

	}

}
