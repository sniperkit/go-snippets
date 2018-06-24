package main

import (
	"encoding/json"
	"github.com/datatogether/task_mgmt/tasks"
	"net"
	"net/rpc"
)

type TasksRequestAct struct {
	ReqAction
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

func (TasksRequestAct) Type() string        { return "TASKS_FETCH_REQUEST" }
func (TasksRequestAct) SuccessType() string { return "TASKS_FETCH_SUCCESS" }
func (TasksRequestAct) FailureType() string { return "TASKS_FETCH_FAILURE" }

func (TasksRequestAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &TasksRequestAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *TasksRequestAct) Exec() (res *ClientResponse) {
	conn, err := net.Dial("tcp", cfg.TasksServiceUrl)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
		return
	}
	cli := rpc.NewClient(conn)
	p := &tasks.TasksListParams{
		Limit:  a.PageSize,
		Offset: (a.Page - 1) * a.PageSize,
	}
	reply := []*tasks.Task{}
	if err := cli.Call("TaskRequests.List", p, &reply); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
		return
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "TASK_ARRAY",
		Data:      reply,
	}
}

type TaskEnqueueAct struct {
	ReqAction
	Title    string                 `json:"title"`
	TaskType string                 `json:"taskType"`
	UserId   string                 `json:"userId"`
	Params   map[string]interface{} `json:"params"`
}

func (TaskEnqueueAct) Type() string        { return "TASK_ENQUEUE_REQUEST" }
func (TaskEnqueueAct) SuccessType() string { return "TASK_ENQUEUE_SUCCESS" }
func (TaskEnqueueAct) FailureType() string { return "TASK_ENQUEUE_FAILURE" }

func (TaskEnqueueAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &TaskEnqueueAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *TaskEnqueueAct) Exec() (res *ClientResponse) {
	log.Infof("adding task %s: %s", a.TaskType, a.Title)
	conn, err := net.Dial("tcp", cfg.TasksServiceUrl)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
		return
	}

	log.Infof("successfully dialed tasks server")
	cli := rpc.NewClient(conn)
	p := &tasks.TasksEnqueueParams{
		Title:  a.Title,
		Type:   a.TaskType,
		UserId: a.UserId,
		Params: a.Params,
	}

	reply := &tasks.Task{}
	log.Infof("enquing task")
	if err := cli.Call("TaskRequests.Enqueue", p, reply); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
		return
	}
	log.Infof("task enqued")
	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "TASK",
		Data:      reply,
	}
}
