// Package processing holds the logic for processing the requests
package processing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"runtime"
	"time"

	"appengine"

	"github.com/gophergala/yps/core"
	"github.com/gophergala/yps/queue/aetq"

	"github.com/gorilla/mux"
)

var (
	publicAPIKey string
)

func init() {
	loadAPIKey()

	r := mux.NewRouter()
	r.HandleFunc("/processPlaylist", queueHandler).Methods("GET")
	http.Handle("/", r)
}

func loadAPIKey() {
	if publicAPIKey != "" {
		return
	}

	_, currentFilename, _, ok := runtime.Caller(0)
	if !ok {
		panic(fmt.Errorf("Could not retrieve the caller for loading the config file"))
	}

	currentDir := path.Dir(currentFilename)
	config, err := ioutil.ReadFile(path.Join(currentDir, "config.json"))
	if err != nil {
		panic(err)
	}

	cfg := &struct {
		PublicAPIKey string `json:"public_api_key"`
	}{}

	if err := json.Unmarshal(config, cfg); err != nil {
		panic(err)
	}

	publicAPIKey = cfg.PublicAPIKey
}

func queueHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	playlistMq := aetq.NewQueue(c, core.PlaylistQueue, 5)
	videoMq := aetq.NewQueue(c, core.VideoQueue, 5)

	// TODO implement playlist fetching for items from youtube API before removing this
	fmt.Fprint(w, "skipped processig for now")
	return

	processChan := make(chan error, 1)

	for i := 0; i < 12; i++ {
		start := time.Now()
		go core.ProcessPlaylistTasks(&playlistMq, &videoMq, processChan)

		select {
		case <-time.After(time.Duration(5) * time.Second):
			{
				log.Printf("[info] processing tasks took too long")
			}
		case resp := <-processChan:
			{
				if resp != nil {
					log.Printf("[error] error while processing playlist message: %q", resp)
				}

				// We do want to process stuff every 5 seconds for now
				time.Sleep(time.Duration(5)*time.Second - time.Since(start))
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "task finished")
}
