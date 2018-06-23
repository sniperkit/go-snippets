package main

import (
	"fmt"
	"github.com/gophergala/heatingeffect/common"
	mgo "gopkg.in/mgo.v2"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	topTenSenderTpl *template.Template
	config          *common.Config
	session         *mgo.Session
)

func main() {
	var err error
	if len(os.Args) < 2 {
		log.Fatal("Missing config path argument")
	}
	configPath := os.Args[1]

	config, err = common.LoadConfig(configPath)
	if err != nil {
		log.Fatal("LoadConfig: ", err)
	}

	session, err = initMongoDB(config)
	if err != nil {
		log.Fatal("initMongoDB: ", err)
	}

	topTenSenderTpl, err = template.ParseFiles("topTenSender.html")
	if err != nil {
		log.Fatal("template.ParseFiles: ", err)
	}

	http.HandleFunc("/", topTenSender)
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func topTenSender(w http.ResponseWriter, req *http.Request) {
	c := session.DB("").C(config.MongoDB.NoticesSendToStatCollectionName)
	stat, err := common.GetNoticesSendStats(c, 10, false)
	if err != nil {
		log.Print("GetNoticesSendStats: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = topTenSenderTpl.Execute(w, stat)
	if err != nil {
		log.Print("topTenSender: ", err)
	}
}

func initMongoDB(config *common.Config) (*mgo.Session, error) {
	if config.MongoDB == nil {
		return nil, fmt.Errorf("Config.MongoDB is nil")
	}
	dialInfo := &mgo.DialInfo{
		Addrs:    config.MongoDB.Addrs,
		Timeout:  config.MongoDB.Timeout,
		Database: config.MongoDB.Database,
		Username: config.MongoDB.Username,
		Password: config.MongoDB.Password,
	}
	return mgo.DialWithInfo(dialInfo)
}
