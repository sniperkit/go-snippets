// Package backend implements the main for the backend requests
package backend

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"appengine"

	"github.com/gophergala/yps/core"
	"github.com/gophergala/yps/queue/aetq"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/processUserInput", queueHandler).Methods("GET")
	http.Handle("/", r)
}

func queueHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	msgMq := aetq.NewQueue(c, core.UserInputQueue, 5)
	playlistMq := aetq.NewQueue(c, core.PlaylistQueue, 5)
	videoMq := aetq.NewQueue(c, core.VideoQueue, 5)

	processChan := make(chan error, 1)

	for i := 0; i < 12; i++ {
		start := time.Now()
		go core.ProcessUserInputTasks(&msgMq, &playlistMq, &videoMq, processChan)

		select {
		case <-time.After(time.Duration(5) * time.Second):
			{
				log.Printf("[info] processing tasks took too long")
			}
		case resp := <-processChan:
			{
				if resp != nil {
					log.Printf("[error] error while processing user input message: %q", resp)
				}

				// We do want to process stuff every 5 seconds for now
				time.Sleep(time.Duration(5)*time.Second - time.Since(start))
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "task finished")
}
