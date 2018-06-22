package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Campaign struct {
	Id           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	AdvertiserId bson.ObjectId `json:"_" bson:"advertiser_id,omitempty"`
	URLRoot      string        `json:"url_root" bson:"url_root"`
	Name         string        `json:"campaign_name" bson:"campaign_name,omitempty"`
}

func (camp *Campaign) ValidateURLRoot() error {

	if !URLRootRegexp.MatchString(camp.URLRoot) {
		return ErrInvalidURLRoot
	}
	return nil
}

// unique index
var URIIndex = mgo.Index{
	Key:    []string{"url_root"},
	Unique: true,
	Sparse: false,
}
