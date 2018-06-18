package tests

import (
	"testing"
	beeTest "github.com/astaxie/beego/testing"
	"io/ioutil"
	"github.com/jiangew/hancock/shorturl/controllers"
	"encoding/json"
)

func TestShort(t *testing.T) {
	request := beeTest.Post("/v1/shorten")
	request.Param("longurl", "http://jiangew.github.io")
	response, _ := request.Response()
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	var short controllers.ShortResult
	json.Unmarshal(content, &short)

	if short.UrlShort == "" {
		t.Fatal("shorturl is empty")
	}
}

func TestExpand(t *testing.T) {
	request := beeTest.Get("/v1/expand")
	request.Param("shorturl", "5laZF")
	response, _ := request.Response()
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)
	var short controllers.ShortResult
	json.Unmarshal(content, &short)

	if short.UrlLong == "" {
		t.Fatal("longurl is empty")
	}
}
