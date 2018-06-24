package main

import (
	"encoding/json"
	"fmt"
	"github.com/datatogether/api/apiutil"
	"github.com/datatogether/task_mgmt/tasks"
	"io"
	"net/http"
	"strconv"
	"time"
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Infoln("tasks req:", r.Method, r.URL.Path)
	switch r.Method {
	case "POST":
		EnqueueTaskHandler(w, r)
	case "GET":
		ListTasksHandler(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func EnqueueTaskHandler(w http.ResponseWriter, r *http.Request) {
	t := &tasks.Task{}
	if err := json.NewDecoder(r.Body).Decode(t); err != nil {
		log.Infoln(err)
		apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	// perform the task raw if no amqp url is specified
	if cfg.AmqpUrl == "" {
		now := time.Now()
		t.Enqueued = &now
		if err := t.Save(store); err != nil {
			apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
			return
		}

		task := tasks.Task{Id: t.Id}
		if err := task.Read(store); err != nil {
			apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
			return
		}

		go func() {
			tc := make(chan *tasks.Task, 10)
			go func() {
				if err := task.Do(store, tc); err != nil {
					log.Println(err.Error())
				}
			}()
			for t := range tc {
				fmt.Println(t.Progress.String())
			}
		}()

		apiutil.WriteMessageResponse(w, "task is running", nil)
		return
	}

	if err := t.Enqueue(store, cfg.AmqpUrl); err != nil {
		log.Infoln(err)
		apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	apiutil.WriteMessageResponse(w, "successfully enqueued task", t)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ReadTaskHandler(w, r)
	case "POST":
		EnqueueTaskHandler(w, r)
	default:
		NotFoundHandler(w, r)
	}
}

func ReadTaskHandler(w http.ResponseWriter, r *http.Request) {
	t := &tasks.Task{
		Id: r.URL.Path[len("/tasks/"):],
	}
	if err := t.Read(store); err != nil {
		apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	apiutil.WriteResponse(w, t)
}

func EnqueueIpfsAddHandler(w http.ResponseWriter, r *http.Request) {
	t := &tasks.Task{
		Type: "ipfs.add",
		Params: map[string]interface{}{
			"url":              r.FormValue("url"),
			"ipfsApiServerUrl": cfg.IpfsApiUrl,
		},
	}

	if err := t.Enqueue(store, cfg.AmqpUrl); err != nil {
		log.Infoln(err)
		apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
		return
	}

	apiutil.WriteMessageResponse(w, "task successfully enqueued", t)
}

func reqParamInt(key string, r *http.Request) (int, error) {
	i, err := strconv.ParseInt(r.FormValue(key), 10, 0)
	return int(i), err
}

func reqParamBool(key string, r *http.Request) (bool, error) {
	return strconv.ParseBool(r.FormValue(key))
}

func ListTasksHandler(w http.ResponseWriter, r *http.Request) {
	p := apiutil.PageFromRequest(r)
	ts, err := tasks.ReadTasks(store, "created DESC", p.Limit(), p.Offset())
	if err != nil {
		log.Infoln(err.Error())
		apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
		return
	}

	apiutil.WritePageResponse(w, ts, r, p)
}

// TODO - restore
func CancelTaskHandler(w http.ResponseWriter, r *http.Request) {
	// t := &tasks.Task{
	// 	Id: r.URL.Path[len("/tasks/cancel/"):],
	// }
	// if err := t.Read(store); err != nil {
	// 	renderError(w, err)
	// 	return
	// }

	// if err := t.Cancel(store); err != nil {
	// 	renderError(w, err)
	// 	return
	// }

	// renderMessage(w, "Task Cancelled", "You've cancelled this task")
}

// HealthCheckHandler is a basic "hey I'm fine" for load balancers & co
// TODO - add Database connection & proper configuration checks here for more accurate
// health reporting
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status" : 200 }`))
}

// EmptyOkHandler is an empty 200 response, often used
// for OPTIONS requests that responds with headers set in addCorsHeaders
func EmptyOkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// CertbotHandler pipes the certbot response for manual certificate generation
func CertbotHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, cfg.CertbotResponse)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{ "status" :  "not found" }`))
}
