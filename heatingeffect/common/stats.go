package common

import (
	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

type NoticesSendStats struct {
	ID struct {
		Sender    string `bson:"sender"`
		Recipient string `bson:"recipient"`
	} `bson:"_id"`
	Notices uint `bson:"notices"`
}

func GetNoticesSendStats(c *mgo.Collection, amount int, sortAsc bool) ([]NoticesSendStats, error) {
	sort := "notices"
	if !sortAsc {
		sort = "-" + sort
	}
	query := c.Find(bson.M{}).Sort(sort).Limit(amount)
	stats := []NoticesSendStats{}
	err := query.All(&stats)
	return stats, err

}
