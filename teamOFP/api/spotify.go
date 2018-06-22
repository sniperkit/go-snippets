package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const spotifyUrl = "https://api.spotify.com/v1"

func getTrackDetails(ID string) *Track {
	t := &Track{}

	resp, err := http.Get(spotifyUrl + "/tracks/" + ID)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	var trackDetails map[string]interface{}
	err = json.Unmarshal(body, &trackDetails)
	if err != nil {
		log.Panic(err)
	}

	//log.Println("Spotify Track Details: ", t)

	//log.Println("Track Details:", trackDetails)
	if trackDetails != nil {
		t.Album = trackDetails["album"].(map[string]interface{})["name"].(string)
		t.Artist = trackDetails["artists"].([]interface{})[0].(map[string]interface{})["name"].(string)
		t.Id = ID
		t.AlbumArt = trackDetails["album"].(map[string]interface{})["images"].([]interface{})[0].(map[string]interface{})["url"].(string)
		t.Name = trackDetails["name"].(string)
		duration := trackDetails["duration_ms"].(float64)
		if err != nil {
			log.Panic(err)
		}

		t.Time = strconv.Itoa(int(duration) / 1000)
	}
	return t
}
