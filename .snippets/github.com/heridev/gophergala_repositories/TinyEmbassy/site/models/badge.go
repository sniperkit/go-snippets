package models

import (
// "labix.org/v2/mgo/bson"
)

type Badge struct {
	Id        string `json:badge" bson:"id,omitempty"`
	IamgeURL  string `json:"image_url", bson:"imageurl"`
	S3BadgeId string `json:"_", bson:"s3_badge_id"`
	Size      ImageSize
}

type ImageSize struct {
	Height int `json:"height", bson:"height"`
	Width  int `json:"width", bosn:"width"`
}
