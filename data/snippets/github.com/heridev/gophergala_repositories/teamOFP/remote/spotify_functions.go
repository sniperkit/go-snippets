package main

import (
	"log"
	"strconv"
)

var nextTrack string

func getTimeLeft() float32 {
	duration, derr := strconv.ParseFloat(callSpotify("duration", ""), 64)
	position, perr := strconv.ParseFloat(callSpotify("position", ""), 64)

	if derr != nil {
		log.Panic(derr)
	}

	if perr != nil {
		log.Panic(perr)
	}

	return float32(duration - position)
}

func getPlayerState() string {
	return callSpotify("state", "")
}

func getCurrentTrack() string {
	return callSpotify("name", "")
}

func getCurrentTrackID() string {
	return callSpotify("id", "")
}

func setCurrentTrack(id string) {
	callSpotify("play_track", "\""+id+"\"")
}

func getNextTrack() string {
	return nextTrack
}

func setNextTrack(track string) {
	nextTrack = track
}
