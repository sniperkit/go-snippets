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

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

const (
	deftAddress  = "http://127.0.0.1:7080/v1/"
	deftHostport = ":9010"
	deftVersion  = "1.0.0"
)

var (
	dtsAddress = deftAddress
	hostport   = deftHostport
	version    = deftVersion
)

func init() {
	if address := os.Getenv("DTS_ADDRESS"); address != "" {
		dtsAddress = address
	}
	if hp := os.Getenv("DTS_HOSTPORT"); hp != "" {
		hostport = hp
	}
	if v := os.Getenv("DTS_VERSION"); v != "" {
		version = v
	}

}

var req = gorequest.New()

func sendCancelOrder(c *gin.Context) {
	data := `{"orderCode":"1111000011110000", "remark":"Just A Test"}`

	log.Println(color.Green("发送撤销工单请求："))
	_, body, errs := req.Post(dtsAddress+"s2t").Type("json").
		Set("Sender", "1.0.0").
		Set("TID", "xioxioxi").
		Set("APICODE", "00000002").
		Send(data).EndBytes()
	if errs != nil {
		log.Printf("Ping Error: %v.", errs)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Ping Error.",
		})
		return
	}
	log.Println(string(data))

	log.Println(color.Yellow("\n收到撤单响应："))
	var resData pts.CommResp
	if err := json.Unmarshal(body, &resData); err != nil {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: err.Error(),
		})
		return
	}
	if resData.Code != pts.SUCCESS {
		log.Printf("Send Cancel Order Error: FAILED, %s.\n", resData.Message)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Send Cancel Order Error: FAILED\n",
		})
		return
	}
	v, _ := json.MarshalIndent(&resData, "", "   ")
	log.Println(string(v))

	c.Data(200, "application/json", []byte{'{', '}'})
}

func run(c *gin.Context) {
	data := `{"orderCode":"12345678", "remark":"Just A Test"}`

	log.Printf("%s : %s", color.Green("发送请求数据"), data)
	_, body, errs := req.Post(dtsAddress+"s2t").Type("json").
		Set("Sender", "1.0.0").
		Set("TID", "a3x77n02").
		Set("APICODE", "00000002").
		Send(data).EndBytes()
	if errs != nil {
		log.Printf("Ping Error: %v.", errs)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Ping Error.",
		})
		return
	}

	var resData pts.CommResp
	if err := json.Unmarshal(body, &resData); err != nil {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: err.Error(),
		})
		return
	}
	if resData.Code != pts.SUCCESS {
		log.Printf("Send Cancel Order Error: FAILED, %s.\n", resData.Message)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Send Cancel Order Error: FAILED\n",
		})
		return
	}

	log.Printf("%s : %s .\n", color.Yellow("得到响应数据"), string(body))

	c.Data(200, "application/json", []byte{'{', '}'})
}

func receiveRefuseOrder(c *gin.Context) {
	var body struct {
		OrderCode  string   `json:"orderCode"`
		ReasonType []string `json:"reasonType"`
		Remark     string   `json:"remark"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: err.Error(),
		})
		return
	}

	// Show Ping Data
	log.Printf("%s : {code: '%s', reasonType: %v, remark: '%s'}.", color.Cyan("收到请求数据"),
		body.OrderCode, body.ReasonType, body.Remark)

	c.JSON(200, &pts.CommResp{
		Code:    pts.SUCCESS,
		Message: "",
		Data:    []byte("{}"),
	})
}

func main() {
	router := gin.Default()
	router.POST("/rescue/refuseOrder", receiveRefuseOrder)
	router.GET("/run/default", run)
	router.GET("/run/cancelOrder", sendCancelOrder)

	if err := router.Run(hostport); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
