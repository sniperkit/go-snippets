package models

import (
	"labix.org/v2/mgo/bson"
)

type BadgeGroup struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	CampaignUUId string        `json:"campaign_uuid" bson:"campaign_uuid,omitempty"`
	Title        string        `json:"title" bson:"title"`
	TargetURL    string        `json:"target_url" bson:"targeturl"`
	Badges       []Badge
}
