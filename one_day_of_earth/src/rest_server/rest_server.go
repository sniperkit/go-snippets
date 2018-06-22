package rest_server

import (
	"APIs"
	"APIs/twitter"
	"APIs/youtube"
	"config"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"mongodatabase"
	"net/http"
	"strconv"
)

type Message struct {
	Error   bool
	Message string
	Data    interface{}
}

const (
	PerPage = 50
)

func StartServer() {
	handler := rest.ResourceHandler{}
	err := handler.SetRoutes(
		&rest.Route{"GET", "/images/:locationhash/:start", func(w rest.ResponseWriter, req *rest.Request) {
			locationHash := req.PathParam("locationhash")
			start_index, err := strconv.Atoi(req.PathParam("start"))
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: "Wrong 'start' parameter",
				})
				return
			}
			m := mongodatabase.Mongo{}
			err = m.Connect()
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: err.Error(),
				})
				return
			}
			//Instagram Images
			images := make([]mongodatabase.InstagramCollection, 0)
			found := false
			found, err = m.FindAll(config.INSTAGRAM_DB_COLLECTION, map[string]interface{}{
				"locationhash": locationHash,
			}, &images)
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: err.Error(),
				})
				return
			}
			if !found {
				w.WriteJson(&Message{
					Error:   true,
					Message: "Images Not Found",
				})
				return
			}
			ret_images := make([]APIs.ApiImage, 0)
			for _, im := range images {
				ret_images = append(ret_images, im.Images...)
			}

			//Instagram Images
			images2 := make([]mongodatabase.FlickrCollection, 0)
			found = false
			found, err = m.FindAll(config.FLICKR_DB_COLLECTION, map[string]interface{}{
				"locationhash": locationHash,
			}, &images2)
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: err.Error(),
				})
				return
			}
			if !found {
				w.WriteJson(&Message{
					Error:   true,
					Message: "Images Not Found",
				})
				return
			}

			for _, im := range images2 {
				ret_images = append(ret_images, im.Images...)
			}

			leng := start_index + PerPage
			if leng >= len(ret_images) {
				ret_images = ret_images[start_index:]
			} else {
				ret_images = ret_images[start_index:leng]
			}
			w.WriteJson(&ret_images)
		}},
		&rest.Route{"GET", "/videos/:locationhash/:start", func(w rest.ResponseWriter, req *rest.Request) {
			locationHash := req.PathParam("locationhash")
			start_index, err := strconv.Atoi(req.PathParam("start"))
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: "Wrong 'start' parameter",
				})
				return
			}
			m := mongodatabase.Mongo{}
			err = m.Connect()
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: err.Error(),
				})
				return
			}
			videos := make([]mongodatabase.YoutubeCollection, 0)
			found := false
			found, err = m.FindAll(config.YOUTUBE_DB_COLLECTION, map[string]interface{}{
				"locationhash": locationHash,
			}, &videos)
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: err.Error(),
				})
				return
			}
			if !found {
				w.WriteJson(&Message{
					Error:   true,
					Message: "Videos Not Found",
				})
				return
			}
			ret_videos := make([]youtube.Video, 0)
			for _, v := range videos {
				ret_videos = append(ret_videos, v.Videos...)
			}

			leng := start_index + PerPage
			if leng >= len(ret_videos) {
				ret_videos = ret_videos[start_index:]
			} else {
				ret_videos = ret_videos[start_index:leng]
			}
			w.WriteJson(&ret_videos)
		}},
		&rest.Route{"GET", "/tweets/:locationhash/:start", func(w rest.ResponseWriter, req *rest.Request) {
			locationHash := req.PathParam("locationhash")
			start_index, err := strconv.Atoi(req.PathParam("start"))
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: "Wrong 'start' parameter",
				})
				return
			}
			m := mongodatabase.Mongo{}
			err = m.Connect()
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: err.Error(),
				})
				return
			}
			tweets := make([]mongodatabase.TwitterCollection, 0)
			found := false
			found, err = m.FindAll(config.TWITTER_DB_COLLECTION, map[string]interface{}{
				"locationhash": locationHash,
			}, &tweets)
			if err != nil {
				w.WriteJson(&Message{
					Error:   true,
					Message: err.Error(),
				})
				return
			}
			if !found {
				w.WriteJson(&Message{
					Error:   true,
					Message: "Tweets Not Found",
				})
				return
			}
			ret_tweets := make([]twitter.Tweet, 0)
			for _, t := range tweets {
				ret_tweets = append(ret_tweets, t.Tweets...)
			}
			leng := start_index + PerPage
			if leng >= len(ret_tweets) {
				ret_tweets = ret_tweets[start_index:]
			} else {
				ret_tweets = ret_tweets[start_index:leng]
			}
			w.WriteJson(&ret_tweets)
		}},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8080", &handler))
}
