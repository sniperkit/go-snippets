// Package dtsproto provide protocol of service request dts
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package dtsproto

import (
	"encoding/json"
	"io/ioutil"
	pts "mdts/protocols"

	"github.com/gin-gonic/gin"
)

type respcode uint8

const (
	SUCCESS respcode = 0
	FAILED  respcode = 1
)

// CommResp is the common response struct
type CommResp struct {
	Code    respcode
	Message string
	Data    json.RawMessage
}

// * S2T

// HeadS2T is the Request Head from service to third
type HeadS2T struct {
	TID     string `binding:"alphanum,len=8"`
	APICODE string `binding:"number,len=8"`
	Sender  string `binding:"serversion"`
}

// ParseS2T parse request from service to third
func ParseS2T(c *gin.Context) (*HeadS2T, []byte, error) {
	var (
		h HeadS2T
	)
	h.TID = c.GetHeader("TID")
	h.APICODE = c.GetHeader("APICODE")
	h.Sender = c.GetHeader("Sender")

	if err := pts.Validator.ValidateStruct(h); err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, nil, err
	}

	return &h, body, nil
}

// * T2S

// HeadT2S is the Request Head from third to service
type HeadT2S struct {
	Version string `binding:"serversion"`
	APICODE string `binding:"number,len=8"`
	Sender  string `binding:"alphanum,len=8"`
}

// ParseS2T parse request from third to service
func ParseT2S(c *gin.Context) (*HeadT2S, []byte, error) {
	var (
		h HeadT2S
	)
	h.Version = c.GetHeader("Version")
	h.APICODE = c.GetHeader("APICODE")
	h.Sender = c.GetHeader("Sender")

	if err := pts.Validator.ValidateStruct(h); err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, nil, err
	}

	return &h, body, nil
}
