package main

import (
	"color"
	"conf"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	bsdk "mdts/brokerSDK/base"
	"mdts/brokerSDK/s2t"
	bmsg "mdts/protocols/brokermsg"
	"strconv"
	"time"
)

const (
	goiaNum   = "11001"
	xioFormat = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:impl="http://impl.webservcice.eis.com/">
   <soapenv:Header/>
   <soapenv:Body>
      <impl:esbServiceOperation>
         <arg0>%s</arg0>
      </impl:esbServiceOperation>
   </soapenv:Body>
</soapenv:Envelope>`

	xioSuccess    = "COMPLETE"
	xioFailed     = "FAIL"
	xioRespFormat = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:impl="http://impl.webservcice.eis.com/"><soapenv:Body><impl:esbServiceOperationResponse><return>%s</return></impl:esbServiceOperationResponse></soapenv:Body></soapenv:Envelope>`
)

type xioRouter struct {
	SourceSysID string `xml:"SourceSysID"`
	ServiceID   string `xml:"ServiceID"`
	SerialNO    string `xml:"SerialNO"`
	ServiceTime string `xml:"ServiceTime"`
}

type xioData struct {
	Control  string      `xml:"Control"`
	Request  interface{} `xml:"Request"`
	Response string      `xml:"Response"`
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
	Control string `xml:"Control"`
	Request struct {
		Value string `xml:",innerxml"`
	} `xml:"Request"`
	Response struct {
		Code    int    `xml:"code"`
		Message string `xml:"message"`
	} `xml:"Response"`
}

type xioRequest struct {
	XMLName xml.Name   `xml:"Service"`
	Router  *xioRouter `xml:"Route"`
	Data    *xioData   `xml:"Data"`
}

type xioResponse struct {
	XMLName xml.Name      `xml:"Envelope"`
	Router  xioRespRouter `xml:"Body>esbServiceOperationResponse>return>Service>Route"`
	Data    xioRespData   `xml:"Body>esbServiceOperationResponse>return>Service>Data"`
}

type CancelOrderStruct struct {
	OrderCode string `json:"orderCode" xml:"orderCode" binding:"required,numeric,max=40"`
	Remark    string `json:"remark" xml:"remark" binding:"max=100"`
}

type AddOrderStruct struct {
	OrderCode  string `json:"orderCode" xml:"orderCode" binding:"required,numeric,max=40"`
	OrderType  int    `json:"orderType" xml:"orderType" binding:"required,gte=1,lte=2"`
	AlarmCode  string `json:"alarmCode" xml:"alarmCode" binding:"required,numeric,max=40"`
	AlarmType  int    `json:"alarmType" xml:"alarmType" binding:"required,gte=9000001,lte=9000004"`
	HappenTime string `json:"happenTime" xml:"happenTime" binding:"required"`
	LiftCode   string `json:"liftCode" xml:"liftCode" binding:"required,alphanum,max=40"`
	StreetAddr string `json:"streetAddr" xml:"streetAddr" binding:"max=120"`
	AreaAddr   string `json:"areaAddr" xml:"areaAddr" binding:"max=32"`
	LiftAddr   string `json:"liftAddr" xml:"liftAddr" binding:"max=40"`
	Longitude  string `json:"longitude" xml:"longitude" binding:"omitempty,longitude,max=20"`
	Latitude   string `json:"latitude" xml:"latitude" binding:"omitempty,latitude,max=20"`
	Remark     string `json:"remark" xml:"remark" binding:"max=100"`
}

// 西奥天梯平台
type xioLift struct {
	sysID                string
	addOrderServiceID    string
	cancelOrderServiceID string
	count                uint
}

var XioLift = &xioLift{
	sysID:                "01012",
	addOrderServiceID:    "01009000000001",
	cancelOrderServiceID: "01009000000002",
	count:                1,
}

func (xl *xioLift) AddOrder(data []byte) (resbody []byte, err error) {
	var addOrder AddOrderStruct
	if err := json.Unmarshal(data, &addOrder); err != nil {
		return nil, err
	}

	route := xl.getRoute(xl.addOrderServiceID)
	d := &xioData{
		Request: addOrder,
	}
	output, err := xml.MarshalIndent(&xioRequest{
		Router: route,
		Data:   d,
	}, "", "   ")
	if err != nil {
		return nil, err
	}

	return output, nil

}

func (xl *xioLift) CancelOrder(data []byte) (resbody []byte, err error) {
	var cancelOrder CancelOrderStruct
	if err := json.Unmarshal(data, &cancelOrder); err != nil {
		return nil, err
	}

	route := xl.getRoute(xl.cancelOrderServiceID)
	d := &xioData{
		Request: cancelOrder,
	}
	output, err := xml.MarshalIndent(&xioRequest{
		Router: route,
		Data:   d,
	}, "", "   ")
	if err != nil {
		return nil, err
	}

	return output, nil

}

func (xl *xioLift) UnmarshalResp(data []byte) (*xioResponse, error) {
	var res xioResponse
	if err := xml.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (xl *xioLift) serviceTime() (dateTime string, timeTime string) {
	t := time.Now()
	bs := make([]byte, 0, 14)
	bs = append(bs, []byte(strconv.Itoa(t.Year()))...)
	bs = append(bs, []byte(strconv.Itoa(int(t.Month())))...)
	bs = append(bs, []byte(strconv.Itoa(t.Day()))...)
	bs = append(bs, []byte(strconv.Itoa(t.Hour()))...)
	bs = append(bs, []byte(strconv.Itoa(t.Minute()))...)
	bs = append(bs, []byte(strconv.Itoa(t.Second()))...)

	return string(bs[:8]), string(bs[8:])
}

// countString ...
func (xl *xioLift) countUpAndFmt() string {
	xl.count++
	xl.count = xl.count % 999999
	return fmt.Sprintf("%06d", xl.count)
}

// getRoute ...
func (xl *xioLift) getRoute(serviceID string) *xioRouter {
	dt, tt := xl.serviceTime()
	return &xioRouter{
		SourceSysID: xl.sysID,
		ServiceID:   serviceID,
		SerialNO:    dt + goiaNum + xl.countUpAndFmt(),
		ServiceTime: dt + tt,
	}
}

type xioTrans struct {
	id  string
	url string
}

func NewXioTrans(id string, url string) *xioTrans {
	return &xioTrans{
		id:  id,
		url: url,
	}
}

func (dt *xioTrans) ID() string {
	return dt.id
}

func (dt *xioTrans) TransTo(APICODE string, Data []byte) (*bsdk.TransToResult, error) {
	if APICODE == "00000001" {
		log.Println(color.Green("收到业务请求数据："))
		log.Println(string(Data))
		v, err := XioLift.AddOrder(Data)
		if err != nil {
			return nil, err
		}
		head := map[string]string{
			"Content-Type": "application/xml",
		}
		log.Println(color.Yellow("将其转化为:"))
		log.Println(string(fmt.Sprintf(xioFormat, string(v))))

		log.Println(string(v))
		return &bsdk.TransToResult{
			Method: bmsg.EnumMethod_HttpPost,
			Head:   head,
			Body:   []byte(fmt.Sprintf(xioFormat, string(v))),
			URL:    dt.url,
		}, nil
	} else if APICODE == "00000002" {
		log.Println(color.Green("收到业务请求数据："))
		log.Println(string(Data))
		v, err := XioLift.CancelOrder(Data)
		if err != nil {
			return nil, err
		}
		head := map[string]string{
			"Content-Type": "application/xml",
		}
		log.Println(color.Yellow("将其转化为:"))
		log.Println(string(fmt.Sprintf(xioFormat, string(v))))

		log.Println(string(v))
		return &bsdk.TransToResult{
			Method: bmsg.EnumMethod_HttpPost,
			Head:   head,
			Body:   []byte(fmt.Sprintf(xioFormat, string(v))),
			URL:    dt.url,
		}, nil
	}
	return &bsdk.TransToResult{
		Method: bmsg.EnumMethod_HttpPost,
		Head:   make(map[string]string),
		Body:   Data,
		URL:    dt.url,
	}, nil
}

func (dt *xioTrans) TransFrom(APICODE string, Data []byte) (*bsdk.TransFromResult, error) {
	if APICODE == "00000001" {
		log.Println(color.Green("收到业务响应数据："))
		log.Println(string(Data))
		res, err := XioLift.UnmarshalResp(Data)
		if err != nil {
			return nil, err
		}

		v, err := json.MarshalIndent(res, "", "   ")
		if err != nil {
			return nil, err
		}
		log.Println(color.Yellow("将其转化为:"))
		log.Println(string(v))

		return &bsdk.TransFromResult{
			Head: make(map[string]string),
			Body: v,
		}, nil
	} else if APICODE == "00000002" {
		log.Println(color.Green("收到业务响应数据："))
		log.Println(string(Data))
		res, err := XioLift.UnmarshalResp(Data)
		if err != nil {
			return nil, err
		}

		v, err := json.MarshalIndent(res, "", "   ")
		if err != nil {
			return nil, err
		}
		log.Println(color.Yellow("将其转化为:"))
		log.Println(string(v))

		return &bsdk.TransFromResult{
			Head: make(map[string]string),
			Body: v,
		}, nil
	}
	return &bsdk.TransFromResult{
		Head: make(map[string]string),
		Body: Data,
	}, nil
}

var confMap = map[string]string{
	"ID":         "xioxioxi",
	"URL":        "http://127.0.0.1:9001/cancelOrder",
	"RegistAddr": "127.0.0.1:9101",
	"Hostport":   ":9101",
}

func init() {
	conf.InitConfMapFromEnv(confMap)
}

func main() {
	trans := NewXioTrans(confMap["ID"], confMap["URL"])
	server := s2t.NewServer(confMap["ID"], confMap["RegistAddr"], trans)

	if err := server.Run(confMap["Hostport"]); err != nil {
		log.Fatalln(err)
	}
}
