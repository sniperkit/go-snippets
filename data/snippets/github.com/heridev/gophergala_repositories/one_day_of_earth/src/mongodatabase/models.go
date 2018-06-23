package mongodatabase

import (
	"APIs"
	"APIs/twitter"
	"APIs/youtube"
	"time"
)

type TwitterCollection struct {
	Index        string //MD5 for combination lat, lng, date, distance
	LocationHash string //MD5 for lat, lng
	DateStr      string //Date for compiring on 1 day period
	Tweets       []twitter.Tweet
	CreateDate   time.Time
}

type YoutubeCollection struct {
	Index        string //MD5 for combination lat, lng, date, distance
	LocationHash string //MD5 for lat, lng
	DateStr      string //Date for compiring on 1 day period
	Videos       []youtube.Video
	CreateDate   time.Time
}

type InstagramCollection struct {
	Index        string //MD5 for combination lat, lng, date, distance
	LocationHash string //MD5 for lat, lng
	DateStr      string //Date for compiring on 1 day period
	Images       []APIs.ApiImage
	CreateDate   time.Time
}

type FlickrCollection struct {
	Index        string //MD5 for combination lat, lng, date, distance
	LocationHash string //MD5 for lat, lng
	DateStr      string //Date for compiring on 1 day period
	Images       []APIs.ApiImage
	CreateDate   time.Time
}

type CityCollection struct {
	LocationHash string //MD5 for lat, lng
	Lat          string
	Lng          string
	Name         string //With Country name
}
