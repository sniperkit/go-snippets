// Package request provide http client to request platforms
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package request

import (
	"crypto/tls"
	"fmt"
	"log"
	"mdts/dts/conf"
	"net/http"

	greq "github.com/parnurzeal/gorequest"
)

// goreqs gorequest池
var goreqs Pool

type actionFunc func(*greq.SuperAgent) (*http.Response, []byte, error)

func useReqDo(action actionFunc) (*http.Response, []byte, error) {
	req, err := goreqs.Get()
	defer func() {
		puterr := goreqs.Put(req)
		if puterr != nil {
			log.Println(puterr)
		}
	}()

	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	return action(req.(*greq.SuperAgent))
}

// Get send get request
func Get(hostport string, params string) (*http.Response, []byte, error) {
	action := func(req *greq.SuperAgent) (resp *http.Response, body []byte, err error) {

		var errs []error
		if params == "" {
			resp, body, errs = req.Get(hostport).Type("json").EndBytes()
		} else {
			resp, body, errs = req.Get(hostport).Query(params).Type("json").EndBytes()
		}

		if errs != nil {
			err = fmt.Errorf("RequestError: %v", errs)
		}

		return
	}

	return useReqDo(action)
}

// PostXML send xml post request and get response
func PostXML(hostport string, v string) (*http.Response, []byte, error) {

	action := func(req *greq.SuperAgent) (resp *http.Response, body []byte, err error) {

		var errs []error
		resp, body, errs = req.Post(hostport).Type("xml").
			Send(v).EndBytes()

		if errs != nil {
			err = fmt.Errorf("RequestError: %v", errs)
		}

		return
	}

	return useReqDo(action)
}

// PostBytes send post request and get response
func PostBytes(hostport string, contentType string, bs []byte) (*http.Response, []byte, error) {

	action := func(req *greq.SuperAgent) (resp *http.Response, body []byte, err error) {

		if contentType == "" {
			contentType = "application/json"
		}

		var errs []error
		resp, body, errs = req.Post(hostport).
			Set("Content-Type", contentType).
			Send(string(bs)).EndBytes()

		if errs != nil {
			err = fmt.Errorf("RequestError: %v", errs)
		}

		return
	}

	return useReqDo(action)
}

// Post send post request and get response
func Post(hostport string, json interface{}) (*http.Response, []byte, error) {

	action := func(req *greq.SuperAgent) (resp *http.Response, body []byte, err error) {

		var errs []error
		if json == nil {
			resp, body, errs = req.Post(hostport).Type("json").
				Send("").EndBytes()
		} else {
			resp, body, errs = req.Post(hostport).Type("json").
				SendStruct(json).EndBytes()
		}

		if errs != nil {
			err = fmt.Errorf("RequestError: %v", errs)
		}

		return
	}

	return useReqDo(action)
}

// PostToBody send post request and get response body to the struct
func PostToBody(hostport string, json interface{}, v interface{}) (*http.Response, []byte, error) {

	action := func(req *greq.SuperAgent) (*http.Response, []byte, error) {

		resp, body, errs := req.Post(hostport).Type("json").
			SendStruct(json).EndStruct(v)
		if errs != nil {
			return nil, nil, fmt.Errorf("RequestError: %v", errs)
		}

		return resp, body, nil
	}

	return useReqDo(action)
}

// InitReqPool 初始化request pool
// 如果 tlsConf == nil 则使用http客户端
func InitReqPool(tlsConf *tls.Config) {
	grs, err := NewHoldPool(conf.ReqPoolSize, func() (Object, error) {
		if conf.Usehttps && tlsConf != nil {
			return greq.New().TLSClientConfig(tlsConf).Timeout(conf.ReqClientTimeOut), nil
		}

		return greq.New().Timeout(conf.ReqClientTimeOut), nil
	})
	if err != nil {
		log.Fatalln("GoRequst Pool Init Failed!")
	}

	goreqs = grs
}
