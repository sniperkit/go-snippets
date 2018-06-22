package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	common "github.com/gophergala/heatingeffect/common"

	"github.com/PuerkitoBio/goquery"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gophergala/heatingeffect/chillingeffects"
	"github.com/iron-io/iron_go/worker"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	log.Infof("os.Args: %+v", os.Args)
	app := cli.NewApp()
	app.Name = "discovery"
	app.Usage = "Discovers the highest notice ID from chillingeffects.org and queues harvester worker tasks up on iron.io."
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

		// Connect to MongoDB
		log.Info("Connecting to MongoDB.")
		session, err := initMongoDB(config)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Info("Get latest notice ID from chillingeffects.org.")
		latestID, err := getLatestNoticeID()
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Infof("Latest ID: %d", latestID)

		latestHarvestedID, err := getLatestHarvestedID(config, session)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Infof("Latest harvested ID: %d", latestHarvestedID)

		log.Infof("Queue iron.io tasks")
		err = queueTasks(config, latestHarvestedID, latestID)
		if err != nil {
			log.Error(err)
		}

		// ShutDown
		log.Info("Shutting down.")
		session.Close()
	}

	app.Run(os.Args)
}

func queueTasks(config *common.Config, latestHarvestedID, latestID int) error {
	if latestHarvestedID >= latestID {
		return nil
	}
	numSchedules := (latestID - latestHarvestedID) / config.RequestsPerWorker
	if (latestID-latestHarvestedID)%config.RequestsPerWorker > 0 {
		numSchedules++
	}
	if numSchedules > config.IronIO.ScheduleTasksLimit {
		numSchedules = config.IronIO.ScheduleTasksLimit
	}
	log.Debugf("Scheduling %d tasks", numSchedules)

	cc := *config
	cc.IronIO = nil
	schedules := make([]worker.Schedule, numSchedules)
	n := 0
	for i := latestHarvestedID; i <= latestID && n < numSchedules; i += config.RequestsPerWorker {
		low := i
		high := i + config.RequestsPerWorker
		if high > latestID {
			high = latestID
		}
		cc.IDRange.Low = low
		cc.IDRange.High = high
		payload, err := json.Marshal(&cc)
		if err != nil {
			return err
		}
		schedules[n] = worker.Schedule{
			CodeName: config.IronIO.CodeName,
			Name:     config.IronIO.Name,
			Payload:  string(payload),
			Label:    config.IronIO.Label,
			Cluster:  config.IronIO.Label,
		}
		n++
	}
	w := worker.New()
	_, err := w.Schedule(schedules...)
	return err
}

func getLatestHarvestedID(config *common.Config, session *mgo.Session) (int, error) {
	c := session.DB(config.MongoDB.Database).C(config.MongoDB.NoticesCollectionName)
	query := c.Find(bson.M{}).Sort("-_id").Limit(1)
	var notice chillingeffects.Notice
	err := query.One(&notice)
	if err == mgo.ErrNotFound {
		return 1, nil
	} else if err != nil {
		return 1, err
	}
	return notice.ID, nil
}

func getLatestNoticeID() (int, error) {
	doc, err := goquery.NewDocument("https://chillingeffects.org/")
	if err != nil {
		return -1, err
	}
	id := -1
	notices := doc.Find("#recent-notices").Children()
	notices.Each(func(i int, s *goquery.Selection) {
		val, exists := s.Attr("id")
		if exists {
			idStr := strings.TrimLeft(val, "notice_")
			idTmp, err := strconv.Atoi(idStr)
			if err != nil {
				return
			}
			if idTmp > id {
				id = idTmp
			}
		}
	})
	if id == -1 {
		return id, fmt.Errorf("Could not find latest ID")
	}
	return id, nil
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
