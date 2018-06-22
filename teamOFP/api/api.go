// Spotify Remote API
//
// TeamOFP - GopherGala 2015
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/crowdmob/goamz/sqs"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"strings"
)

// App context
type Context struct {
	db *sqlx.DB
	//airbrake *gobrake.Notifier
	sqs   *sqs.Queue
	rsqs  *sqs.Queue
	tq    *TrackQueue
	np    *nowPlaying
	oauth *oauth2.Config
}

var context = &Context{}
var store = sessions.NewCookieStore([]byte("Groupify.go FTW!"))

// GetInfo - Info Endpoint. Returns versioning info.
func GetInfo(w http.ResponseWriter, r *http.Request) {

	getTrackDetails("0eGsygTp906u18L0Oimnem")
	fmt.Fprintf(w, "Spotify Remote API v0.1.0")
}

func SearchSpotify(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(context.tq.list())
	w.Write(response)
}

//var tracks []Track

func init() {
	//s1 := Track{Time: "3:07", Name: "The Ride", Artist: "David Allan Cole", Album: "16 Biggest Hits"}
	//s2 := Track{Time: "1:20", Name: "Bookends", Artist: "Simon and Garfunkel", Album: "Greatest Hits"}
	//s3 := Track{Time: "3:28", Name: "A Woman Left Lonely", Artist: "Janis Joplin", Album: "The Pearl Sessions"}
	//tracks = []Track{s1, s2, s3}
}

func main() {
	log.Println("Starting Spotify Remote API...")

	// Load .env
	err := godotenv.Load()
	if err != nil {
		// Can't load .env, so setenv defaults
		os.Setenv("SQL_HOST", "localhost:8091")
		os.Setenv("SQL_USER", "root")
		os.Setenv("SQL_PASSWORD", "")
		os.Setenv("SQL_DB", "spotify_remote")
	}

	// Setup App Context
	// Setup DB
	db, err := sqlx.Open("sqlite3", "./spotify-remote.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	context.db = db

	// Setup SQS
	s, err := sqs.NewFrom(os.Getenv("AWS_ACCESS"), os.Getenv("AWS_SECRET"), "us-east-1")
	if err != nil {
		log.Panic(err)
	}

	// Send Queue: API -> Remote
	q, err := s.GetQueue("spotify-ofp")
	if err != nil {
		log.Panic(err)
	}
	context.sqs = q

	// Receive Queue: Remote -> API
	qq, err := s.GetQueue("spotify-ofp-notification")
	if err != nil {
		log.Panic(err)
	}
	context.rsqs = qq

	// Queue Processing Logic
	messages := make(chan *sqs.Message)
	go listenOnQueue(context.rsqs, messages)
	go processQueue(messages)

	// Track Queue
	tq := &TrackQueue{}
	context.tq = tq

	// Now Playing
	context.np = &nowPlaying{}

	router := mux.NewRouter()
	r := router.PathPrefix("/api/v1").Subrouter() // Prepend API Version

	// Setup Negroni
	n := negroni.Classic()

	// Info
	r.HandleFunc("/", GetInfo).Methods("GET")

	// TrackQueue
	r.HandleFunc("/queue/add", PostAddTrack).Methods("POST")
	r.HandleFunc("/queue/list", GetListTracks).Methods("GET")
	r.HandleFunc("/queue/delete", PostDeleteTrack).Methods("POST")
	r.HandleFunc("/queue/next", PostSkipTrack).Methods("POST")
	//r.HandleFunc("/queue/upvote", AddTrack).Methods("POST")
	//r.HandleFunc("/queue/downvote", AddTrack).Methods("POST")

	r.HandleFunc("/search", SearchSpotify).Methods("GET")

	r.HandleFunc("/auth", Auth).Methods("GET")
	r.HandleFunc("/callback", Callback).Methods("GET")

	//r.HandleFunc("/push", QueueTrackRemote).Methods("GET")

	//tq := &TrackQueue{}

	//trackList := tq.list()
	//log.Println("Track Queue: ", trackList)

	//tq.push("song1")
	//tq.push("song2")

	//trackList = tq.list()
	//log.Println("Track Queue: ", trackList)

	//track, _ := tq.pop()
	//log.Println("Track: ", track)

	//trackList = tq.list()
	//log.Println("Track Queue: ", trackList)

	//log.Println("Track Queue Length: ", tq.length())

	// Setup router
	n.UseHandler(r)

	// Start Serve
	if os.Getenv("PORT") != "" {
		n.Run(strings.Join([]string{":", os.Getenv("PORT")}, ""))
	} else {
		n.Run(":8080")
	}

}
