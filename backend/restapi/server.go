package restapi

import (
	"fmt"
	"net"
	"net/http"

	"github.com/garethr/kubeval/log"
	"github.com/gorilla/mux"
)

type Config struct {
	NoUi *bool
}

func StartRestApi(port int, deps *Config) func() error {
	r := mux.NewRouter()
	log.Info(fmt.Sprintf("Sherlock REST API reachable :%d", port))
	if *deps.NoUi == false {
		log.Info(fmt.Sprintf("UI exposed at :%d/", port))
		r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("ui"))))
	} else {
		log.Info("UI not exposed.")
	}

	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

	r.HandleFunc("/api/trace/{trace_id}", func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `{
  "spans": [
    {
      "trace": "1234",
      "id": "1",
      "name": "nginx",
      "start_time": "2011-08-30T09:30:16.768-04:00",
      "duration_ns": 58000000000
    },
    {
      "trace": "1234",
      "id": "2",
      "name": "gateway",
      "start_time": "2011-08-30T09:31:16.768-04:00",
      "duration_ns": 10000000000
    },
    {
      "trace": "1234",
      "id": "3",
      "name": "kafka",
      "start_time": "2011-08-30T09:31:50.768-04:00",
      "duration_ns": 2000000000
    },
    {
      "trace": "1234",
      "id": "4",
      "name": "etcd",
      "start_time": "2011-08-30T09:32:00.768-04:00",
      "duration_ns": 5563463456
    }
  ],
  "relations": [
    {
      "trace": "1234",
      "parent_id": "1",
      "id": "2"
    },
    {
      "trace": "1234",
      "parent_id": "2",
      "id": "3"
    },
    {
      "trace": "1234",
      "parent_id": "1",
      "id": "4"
    }
  ]
}`)
	}).Methods("GET")

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Info(err)
	}
	http.Serve(l, r)
	return l.Close
}
