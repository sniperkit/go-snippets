// Package main provide ping third
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package main

import (
	"color"
	"encoding/json"
	"log"
	"os"
	pts "players/protocols/dtsproto"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	deftAddress  = "http://127.0.0.1:7081/v1/"
	deftHostport = ":9000"
	deftTid      = "a3x77n02"
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

var req = gorequest.New()

func receiveCancelOrder(c *gin.Context) {
	var body struct {
		OrderCode string `json:"orderCode"`
		Remark    string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Printf("Receive Cancel Order Error: %v.", err)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: err.Error(),
		})
		return
	}

	log.Printf("%s : {code: '%s', remark: '%s'}.", color.Cyan("收到请求数据"),
		body.OrderCode, body.Remark)

	go func() {

		time.Sleep(time.Millisecond * 500)

		data := struct {
			OrderCode  string   `json:"orderCode"`
			ReasonType []string `json:"reasonType"`
			Remark     string   `json:"remark"`
		}{
			OrderCode:  body.OrderCode,
			ReasonType: []string{"999", "666"},
			Remark:     "Just A Test",
		}

		log.Printf(`%s : {"orderCode": %s, "reasonType": ["999", "666"], "remark": "Just A Test"}`, color.Green("发送请求数据"), body.OrderCode)
		_, body, errs := req.Post(dtsAddress+"t2s").Type("json").
			Set("Version", "1.0.0").
			Set("APICODE", "00000003").
			Set("Sender", tid).
			SendStruct(&data).EndBytes()
		if errs != nil {
			log.Printf("Ping Error: %v.", errs)
			c.JSON(200, &pts.CommResp{
				Code:    pts.FAILED,
				Message: "Ping Error",
			})
			return
		}

		var resData pts.CommResp
		if err := json.Unmarshal(body, &resData); err != nil {
			log.Printf("Send Refuse Order Error: %v.\n", err)
			c.JSON(200, &pts.CommResp{
				Code:    pts.FAILED,
				Message: err.Error(),
			})
			return
		}

		if resData.Code != pts.SUCCESS {
			log.Printf("Send Refuse Order Error: FAILED, %s.\n", resData.Message)
			c.JSON(200, &pts.CommResp{
				Code:    pts.FAILED,
				Message: "Send Refuse Order Error: FAILED.\n",
			})
			return
		}

		log.Printf("%s : %s .\n", color.Yellow("得到响应数据"), string(body))
	}()

	c.JSON(200, &pts.CommResp{
		Code:    pts.SUCCESS,
		Message: "",
		Data:    []byte("{}"),
	})
}

func main() {
	router := gin.Default()
	router.POST("/cancelOrder", receiveCancelOrder)

	if err := router.Run(hostport); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
