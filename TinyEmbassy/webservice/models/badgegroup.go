/*
* @Author: souravray
* @Date:   2015-01-25 07:33:48
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 11:26:58
 */

package models

import (
	"gopkg.in/mgo.v2/bson"
)

type BadgeGroup struct {
	BadgeGroupId bson.ObjectId `json:"_" bson:"_id,omitempty"`
	CampaignId   bson.ObjectId `json:"_" bson:"_id,omitempty"`
	Title        string        `json:"title" bson:"title"`
	TargetURL    string        `json:"target_url" bson:"targeturl"`
	Badges       []Badge
}

type Badge struct {
	Id        string `json:badge_id" bson:"id,omitempty"`
	IamgeURL  string `json:"image_url", bson:"imageurl"`
	S3BadgeId string `json:"_", bson:"s3_badge_id"`
	Size      ImageSize
}

type ImageSize struct {
	Height int `json:"height", bson:"height"`
	Width  int `json:"width", bosn:"width"`
}
