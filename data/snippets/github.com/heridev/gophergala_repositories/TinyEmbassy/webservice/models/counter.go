/*
* @Author: souravray
* @Date:   2015-01-24 17:08:56
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 11:26:55
 */

package models

import (
	//	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Counter struct {
	U int `json:"unique", bson:"u"`
	I int `json:"impression", bosn:"i"`
	C int `json:"click", bosn:"c"`
}

type CampaignCounter struct {
	Id  bson.ObjectId               `json:"id",bson:"_id,omitempty"`
	Ct  Counter                     `json:"counter", bson:"ct"`
	Tct map[string]TimedCounter     `json:"time", bson:"tct"`
	Pct map[string]PublisherCounter `json:"publisher", bson:"pct"`
}

type TimedCounter struct {
	H map[string]BadgeCounter `json:"hour", bson:"h"`
	D map[string]BadgeCounter `json:"day", bson:"d"`
	W map[string]BadgeCounter `json:"week", bson:"w"`
}

type BadgeCounter struct {
	Ct  Counter            `json:"counter", bson:"ct"`
	Bct map[string]Counter `json:"badges", bson:"bct"`
}

type PublisherCounter struct {
	Ct  Counter            `json:"counter", bson:"ct"`
	Hct map[string]Counter `json:"hour", bson:"hct"`
	Dct map[string]Counter `json:"day", bson:"dct"`
	Wct map[string]Counter `json:"week", bson:"wct"`
}
