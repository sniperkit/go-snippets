package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gophergala/heatingeffect/chillingeffects"
	"github.com/gophergala/heatingeffect/common"
	mgo "gopkg.in/mgo.v2"
)

func main() {
	log.Infof("os.Args: %+v", os.Args)
	app := cli.NewApp()
	app.Name = "harvester"
	app.Usage = "Harvests notices from chillingeffects.org"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "payload",
			Value: "",
			Usage: "Config to run this tool.",
		},
		cli.StringFlag{
			Name:  "config, c",
			Value: "",
			Usage: "Config file to run this tool.",
		},
		cli.StringFlag{
			Name:  "e",
			Value: "",
			Usage: "Enviroment.",
		},
		cli.StringFlag{
			Name:  "id",
			Value: "",
			Usage: "Task ID.",
		},
		cli.StringFlag{
			Name:  "d",
			Value: "",
			Usage: "User writeable directory for temporary storage.",
		},
	}
	app.Action = func(c *cli.Context) {
		// Check parameters
		var (
			payloadFileName = c.String("payload")
		)
		log.Info("Check commandline parameters.")
		if len(payloadFileName) == 0 {
			log.Error("The payload parameter is empty.")
			cli.ShowAppHelp(c)
			return
		}

		// Get config
		log.Infof("Loading config file \"%s\".", payloadFileName)
		config, err := common.LoadConfig(payloadFileName)
		if err != nil {
			log.Error(err)
			return
		}

		// Set logging
		if config.RunMode == "debug" {
			log.SetLevel(log.DebugLevel)
		} else {
			log.SetLevel(log.WarnLevel)
		}
		log.Debugf("Config: %+v", config)

		// Connect to MongoDB
		log.Info("Connecting to MongoDB.")
		session, err := initMongoDB(config)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Harvesting notices and upserting to database.")
		harvestNotices(config, session)

		// ShutDown
		log.Info("Shutting down.")
		session.Close()
	}

	app.Run(os.Args)
}

func harvestNotices(config *common.Config, session *mgo.Session) {
	reports := make(chan *report)
	for i := config.IDRange.Low; i <= config.IDRange.High; i++ {
		go harvestNotice(i, config, session, reports)
	}
	var report *report
	for i := config.IDRange.Low; i <= config.IDRange.High; i++ {
		report = <-reports
		if report.err != nil {
			log.Errorf("HarvestNotice %d failed: %s", report.id, report.err)
		}
	}
}

func harvestNotice(id int, config *common.Config, session *mgo.Session, reports chan<- *report) {
	notice, err := chillingeffects.RequestNotice(id)
	if err != nil {
		reports <- &report{
			id:  id,
			err: err,
		}
		return
	}

	c := session.DB(config.MongoDB.Database).C(config.MongoDB.NoticesCollectionName)
	err = c.Insert(notice)
	reports <- &report{
		id:  id,
		err: err,
	}
}

type report struct {
	id  int
	err error
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
