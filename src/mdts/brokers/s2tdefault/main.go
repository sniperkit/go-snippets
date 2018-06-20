package main

import (
	"color"
	"conf"
	"log"
	bsdk "mdts/brokerSDK/base"
	"mdts/brokerSDK/s2t"
	bmsg "mdts/protocols/brokermsg"
)

type defaultTrans struct {
	id  string
	url string
}

func NewDefaultTrans(id string, url string) *defaultTrans {
	return &defaultTrans{
		id:  id,
		url: url,
	}
}

func (dt *defaultTrans) ID() string {
	return dt.id
}

func (dt *defaultTrans) TransTo(APICODE string, Data []byte) (*bsdk.TransToResult, error) {
	log.Printf("%s : %s .", color.Green("收到业务请求数据"), string(Data))
	log.Println(color.Red("不经任何转换"))
	return &bsdk.TransToResult{
		Method: bmsg.EnumMethod_HttpPost,
		Head:   make(map[string]string),
		Body:   Data,
		URL:    dt.url,
	}, nil
}

func (dt *defaultTrans) TransFrom(APICODE string, Data []byte) (*bsdk.TransFromResult, error) {
	log.Printf("%s : %s .", color.Yellow("收到业务响应数据"), string(Data))
	log.Println(color.Red("不经任何转换"))
	return &bsdk.TransFromResult{
		Head: make(map[string]string),
		Body: Data,
	}, nil
}

var confMap = map[string]string{
	"ID":         "a3x77n02",
	"URL":        "http://127.0.0.1:9000/cancelOrder",
	"RegistAddr": "127.0.0.1:9100",
	"Hostport":   ":9100",
}

func init() {
	conf.InitConfMapFromEnv(confMap)
}

func main() {
	trans := NewDefaultTrans(confMap["ID"], confMap["URL"])
	server := s2t.NewServer(confMap["ID"], confMap["RegistAddr"], trans)

	if err := server.Run(confMap["Hostport"]); err != nil {
		log.Fatalln(err)
	}
}
