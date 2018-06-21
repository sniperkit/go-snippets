package controller

import (
	"../common"
	"io/ioutil"
	"net/http"
)

type VPN struct {
	http.Handler
}

func (t *VPN) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		common.Render().Text(w, http.StatusMethodNotAllowed, "method not allowed")
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	req.ParseForm()
	url := req.PostFormValue("data")
	if url != "" {
		response, err := http.Get(url)
		if err != nil {
			common.Render().Text(w, http.StatusInternalServerError, err.Error())
		}
		body, _ := ioutil.ReadAll(response.Body)
		common.Render().Text(w, response.StatusCode, string(body))
	}
}
