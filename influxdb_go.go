package main

import (
	"github.com/influxdata/influxdb/client/v2"
	"time"
	"os"
	"bufio"
	"strings"
	"io"
	"github.com/go-simplejson"
	//"fmt"
)

const (
	DB = "density"
	timeLayOut = "2006-01-02 15:04:05"
)

var tags = map[string]string{"source": "client"}
var timemap = make(map[string]int)

func main() {
	WrReadLine("collector.json",printline)
	toinflux(timemap)
}

func WrReadLine(filename string, handler func(string)) error {
	// 逐行读取某个日志文件，并对每行文件执行handler 函数
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}

func printline(rawstr string) {
	js, err := simplejson.NewJson([]byte(rawstr))
	if err != nil {
		return
	}
	switch{
	case strings.Contains(rawstr, "SessionStat"):
		timestr, _ := js.Get("SessionStat").Get("time").String()
		timemap[timestr]++
	case strings.Contains(rawstr, "NetStat"):
		timestr, _ := js.Get("NetStat").Get("time").String()
		timemap[timestr]++
	case strings.Contains(rawstr, "BaseInfo"):
		timestr, _ := js.Get("BaseInfo").Get("time").String()
		timemap[timestr]++
	}
}

func toinflux(srcmap map[string]int) {
	// Create a new HTTPClient
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:	"http://192.168.1.32:8086",
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// Create a point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: DB,
		Precision: "s",
	})
	if err != nil {
		panic(err)
	}



	// 把从文本中获取到的字典数据按每秒生成数据点，然后放入batch points
	for k, v := range srcmap {
        fields := map[string]interface{}{
        	"density": v,
		}
        loc, _ := time.LoadLocation("Local")
        theTime, _ := time.ParseInLocation(timeLayOut, k, loc)
        sr := theTime.Local()

        pt, err := client.NewPoint("density", tags, fields, sr)
        if err != nil {
        	return
		}
		bp.AddPoint(pt)
	}

	c.Write(bp);
}



