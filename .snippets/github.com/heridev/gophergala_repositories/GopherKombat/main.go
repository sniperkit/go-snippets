package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	FILE_SERVE_PATH = "/var/static"
)

func main() {
	InitConfig()
	env := flag.String("env", "dev", "Environment")
	flag.Parse()
	fmt.Println(env)
	fileServePath := Config["fileserve-path-"+*env].(string)
	http.HandleFunc("/login/callback", loginCallback)
	panic(http.ListenAndServe(":8080", http.FileServer(http.Dir(fileServePath))))
}

var Config map[string]interface{}

func InitConfig() {
	Config = make(map[string]interface{})
	Config["fileserve-path-dev"] = "static"
	Config["fileserve-path-prod"] = "/var/static"
}

func loginCallback(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	code := params["code"][0]

	hc := http.Client{}

	form := url.Values{}
	form.Add("client_id", "fe6528d512e0697b7883")
	form.Add("client_secret", "035673ffc1d62ec1c6870df76952eb3a0f4cfb1f")
	form.Add("code", code)
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", nil)
	//strings.NewReader(form.Encode())
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	hc.Do(req)

}
