package model

import ()

type Envelope struct {
	Action string `json:"action"`
	Info   *Info  `json:"info,omitempty"`
	Tweet  *Tweet `json:"tweet,omitempty"`
}

// Info action
var INFO = "INFO"

// Info payload
type Info struct {
	UserCount int `json:"user_count"`
}

// Tweet action
var TWEET = "TWEET"

// Tweet payload
type Tweet struct {
	Picture     string  `json:"picture"`
	DisplayName string  `json:"display_name"`
	ScreenName  string  `json:"screen_name"`
	Message     string  `json:"message"`
	Lat         float64 `json:"lat"`
	Long        float64 `json:"long"`
}
