package main

import (
	"errors"
	"fmt"
	"github.com/guotie/gogb2312"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
>公司代码<>公司简称<>公司全称<>英文名称<>注册地址<>A股代码<>A股简称<>A股上市日期<>A股总股本<>A股流通股本<>B股代码<>B股 简 称<>B股上市日期<>B股总股本<>B股 流通股本<>地      区<>省    份<>城     市<>所属行业<>公司网址<
*/

type Stock struct {
	Code            string `json:"Code"`
	ShortName       string `json:"ShortName"`
	FullName        string `json:"FullName"`
	EngName         string `json:"EngName"`
	RegistryAddress string `json:"RegistryAddress"`
	TypeACode       string `json:"TypeACode"`
	TypeAShortName  string `json:"TypeAShortName"`
	TypeAOnSaleDate string `json:"TypeAOnSaleDate"`
	TypeATotalValue uint64 `json:"TypeATotalValue"`
	TypeACurrency   uint64 `json:"TypeACurrency"`
	TypeBCode       string `json:"TypeBCode"`
	TypeBShortName  string `json:"TypeBShortName"`
	TypeBOnSaleDate string `json:"TypeBOnSaleDate"`
	TypeBTotalValue uint64 `json:"TypeBTotalValue"`
	TypeBCurrency   uint64 `json:"TypeBCurrency"`
	Province        string `json:"Province"`
	City            string `json:"City"`
	Area            string `json:"Area"`
	Web             string `json:"Web"`
	Catagory        string `json:"Catagory"`
}

var shenzhen = "http://www.szse.cn/szseWeb/ShowReport.szse?SHOWTYPE=EXCEL&CATALOGID=1110&tab1PAGENUM=1&ENCODE=1&TABKEY=tab1"

func main() {
	client := http.Client{}
	req, err := http.NewRequest("GET", shenzhen, nil)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	//Need convert the encode to UTF-8
	body, err, _, _ = gogb2312.ConvertGB2312(body)
	if err != nil {
		log.Println(err)
		return
	}

	str := ShenZhenPageModify(string(body))
	stocks := strings.Split(str, "<tr")
	for _, s := range stocks {
		fmt.Println(s)
		stock, _ := ShenZhenStockInfoGet(s)
		fmt.Printf("%#v\n", stock)
	}
}

func ShenZhenPageModify(body string) string {
	str := strings.Replace(body, "align", "", -1)
	str = strings.Replace(str, "/td>", "", -1)
	str = strings.Replace(str, "<td", "", -1)
	str = strings.Replace(str, "center", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "left", "", -1)
	str = strings.Replace(str, "right", "", -1)
	str = strings.Replace(str, "=", "", -1)
	str = strings.Replace(str, "'", "", -1)
	str = strings.Replace(str, "{", "", -1)
	str = strings.Replace(str, "}", "", -1)
	str = strings.Replace(str, ">", "", -1)
	str = strings.Replace(str, "bgcolor#ffffff", "", -1)
	str = strings.Replace(str, "bgcolor#F8F8F8", "", -1)
	return str
}

func ShenZhenStockInfoFilter(stock string) string {
	s := strings.Replace(stock, "<tr", "", -1)
	s = strings.Replace(s, "</tr", "", -1)
	return s
}

func ShenZhenStockInfoGet(stock string) (*Stock, error) {
	var s = &Stock{}
	infos := strings.Split(stock, "<")
	if len(infos) < 20 {
		return nil, errors.New("Invalid input")
	}

	s.Code = infos[0]
	s.ShortName = infos[1]
	s.FullName = infos[2]
	s.EngName = infos[3]
	s.RegistryAddress = infos[4]
	s.TypeACode = infos[5]
	s.TypeAShortName = infos[6]
	s.TypeAOnSaleDate = infos[7]
	aTotal := strings.Replace(infos[8], ",", "", -1)
	s.TypeATotalValue, _ = strconv.ParseUint(aTotal, 10, 64)
	aCurrency := strings.Replace(infos[9], ",", "", -1)
	s.TypeACurrency, _ = strconv.ParseUint(aCurrency, 10, 64)
	s.TypeBCode = infos[10]
	s.TypeBShortName = infos[11]
	s.TypeBOnSaleDate = infos[12]
	bTotal := strings.Replace(infos[13], ",", "", -1)
	s.TypeBTotalValue, _ = strconv.ParseUint(bTotal, 10, 64)
	bCurrency := strings.Replace(infos[14], ",", "", -1)
	s.TypeBCurrency, _ = strconv.ParseUint(bCurrency, 10, 64)
	s.Area = infos[15]
	s.Province = infos[16]
	s.City = infos[17]
	s.Catagory = infos[18]
	s.Web = infos[19]

	return s, nil
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
