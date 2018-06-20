// Package main provide ping third
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package main

import (
	"color"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	deftAddress  = "http://127.0.0.1:7081/v1/"
	deftHostport = ":9001"
	deftTid      = "xioxioxi"
)

var (
	dtsAddress = deftAddress
	hostport   = deftHostport
	tid        = deftTid
)

func init() {
	if address := os.Getenv("DTS_ADDRESS"); address != "" {
		dtsAddress = address
	}
	if hp := os.Getenv("DTS_HOSTPORT"); hp != "" {
		hostport = hp
	}
	if v := os.Getenv("DTS_TID"); v != "" {
		tid = v
	}

}

const (
	goiaNum       = "11001"
	xioRespFormat = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:impl="http://impl.webservcice.eis.com/"><soapenv:Body><impl:esbServiceOperationResponse><return>%s</return></impl:esbServiceOperationResponse></soapenv:Body></soapenv:Envelope>`
)

type xioRouter struct {
	SourceSysID string `xml:"SourceSysID"`
	ServiceID   string `xml:"ServiceID"`
	SerialNO    string `xml:"SerialNO"`
	ServiceTime string `xml:"ServiceTime"`
}

type xioData struct {
	Control string `xml:"Control"`
	Request struct {
		Value string `xml:",innerxml"`
	} `xml:"Request"`
	Response string `xml:"Response"`
}

type xioRespRouter struct {
	xioRouter
	ServiceResponse struct {
		Status string `xml:"Status"`
		Code   string `xml:"Code"`
		Desc   string `xml:"Desc"`
	} `xml:"ServiceResponse"`
}

type xioRespData struct {
	Control  string `xml:"Control"`
	Request  string `xml:"Request"`
	Response struct {
		Code    int    `xml:"code"`
		Message string `xml:"message"`
	} `xml:"Response"`
}

type xioRequest struct {
	XMLName xml.Name  `xml:"Envelope"`
	Router  xioRouter `xml:"Body>esbServiceOperation>arg0>Service>Route"`
	Data    xioData   `xml:"Body>esbServiceOperation>arg0>Service>Data"`
}

type xioResponse struct {
	XMLName xml.Name      `xml:"Service"`
	Router  xioRespRouter `xml:"Route"`
	Data    xioRespData   `xml:"Data"`
}

var req = gorequest.New()

func combineResponse(status string, code int, message string) []byte {
	var resp xioResponse
	resp.Router.ServiceResponse.Status = status
	resp.Data.Response.Code = code
	resp.Data.Response.Message = message

	output, err := xml.MarshalIndent(&resp, "", "   ")
	if err != nil {
		log.Fatalln(err)
	}

	return []byte(fmt.Sprintf(xioRespFormat, string(output)))
}

func receiveCancelOrder(c *gin.Context) {
	datas, err := c.GetRawData()
	if err != nil {
		log.Println(err)
		c.Data(200, "application/xml", combineResponse("FAIL", 100, ""))
		return
	}

	var req xioRequest
	err = xml.Unmarshal(datas, &req)
	if err != nil {
		log.Println(err)
		c.Data(200, "application/xml", combineResponse("FAIL", 100, ""))
		return
	}

	log.Println(color.Green("收到撤销工单请求："))
	output, err := xml.MarshalIndent(&req, "", "   ")
	if err != nil {
		log.Println(err)
		c.Data(200, "application/xml", combineResponse("FAIL", 100, ""))
		return
	}
	log.Println(string(output))

	log.Println(color.Yellow("\n响应撤单请求："))
	res := combineResponse("COMPLETE", 0, "")
	log.Println(string(res))

	c.Data(200, "application/xml", res)
}

func receiveAddOrder(c *gin.Context) {
	datas, err := c.GetRawData()
	if err != nil {
		log.Println(err)
		c.Data(200, "application/xml", combineResponse("FAIL", 100, ""))
		return
	}

	var req xioRequest
	err = xml.Unmarshal(datas, &req)
	if err != nil {
		log.Println(err)
		c.Data(200, "application/xml", combineResponse("FAIL", 100, ""))
		return
	}

	log.Println(color.Green("收到撤销工单请求："))
	output, err := xml.MarshalIndent(&req, "", "   ")
	if err != nil {
		log.Println(err)
		c.Data(200, "application/xml", combineResponse("FAIL", 100, ""))
		return
	}
	log.Println(string(output))

	log.Println(color.Yellow("\n响应撤单请求："))
	res := combineResponse("COMPLETE", 0, "")
	log.Println(string(res))

	c.Data(200, "application/xml", res)
}

func main() {
	router := gin.Default()
	router.POST("/cancelOrder", receiveCancelOrder)
	router.POST("/addOrder", receiveAddOrder)

	if err := router.Run(hostport); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
