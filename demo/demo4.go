/*
* @Author: wuhailin
* @Date:   2017-10-07 17:17:03
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-10-07 17:17:39
 */
package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

type Data struct {
	test []string
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "POST")
		if http.MethodPost == r.Method {
			if err := r.ParseForm(); nil != err {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Println(`post data:`, r.PostForm)
			for key, values := range r.PostForm {
				if `test[]` == key {
					go checkUrl(values)
				}
			}
			io.WriteString(w, `{"status":1, "msg":"success"}`)
		} else {
			http.NotFound(w, r)
		}
	})
	http.ListenAndServe("127.0.0.1:8051", nil)
}

func checkUrl(values []string) {
	urls := []string{}
	for _, value := range values {
		if _, err := url.ParseRequestURI(value); nil != err {
			log.Println(`url:`, value, ` parse fail`)
			continue
		}
		urls = append(urls, value)
	}
	log.Println(urls)
}
