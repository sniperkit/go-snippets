package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/datatogether/sql_datastore"
	"github.com/datatogether/sqlutil"
	"github.com/ipfs/go-datastore"
	"github.com/pborman/uuid"
	"github.com/streadway/amqp"
	"time"
)

// Task represents the storable state of a task. Note this is not the "task" itself
// (the function that will be called to do the actual work associated with a task)
// but the state associated with performing a task.
// Task holds the type of work to be done, parameters to configure the work
// to be done, and the status of that work. different types of "work" are done by
// implementing the Taskable interface specified in taskdef.go
// lots of the methods on Task overlap with Taskable, this is on purpose, as Task
// wraps phases of task completion to track the state of a task
type Task struct {
	// uuid identifier for task
	Id string `json:"id"`
	// created date rounded to secounds
	Created time.Time `json:"created"`
	// updated date rounded to secounds
	Updated time.Time `json:"updated"`
	// human-readable title for the task, meant to be descriptive & varied
	Title string `json:"title"`
	// id of user that submitted this task
	UserId string `json:"userId"`
	// Type of task to be executed
	Type string `json:"type"`
	// parameters supplied to the task, should be json bytes
	Params map[string]interface{} `json:"params"`
	// Status Message
	Status string `json:"status,omitempty"`
	// Error Message
	Error string `json:"error,omitempty"`
	// timstamp for when request was added to the tasks queue
	// nil if request hasn't been sent to the queue
	Enqueued *time.Time `json:"enqueued,omitempty"`
	// timestamp for when the task was removed from the queue
	// and started, nil if the request hasn't been started
	Started *time.Time `json:"started,omitempty"`
	// timestamp for when request succeeded
	// nil if task hasn't succeeded
	Succeeded *time.Time `json:"succeeded,omitempty"`
	// timestamp for when request failed
	// nil if task hasn't failed
	Failed *time.Time `json:"failed,omitempty"`
	// progress of this task's completion
	// progress may not be stored, but instead kept ephemerally
	Progress *Progress `json:"progress,omitempty"`
}

// DatastoreType is to fulfill the sql_datastore.Model interface
// It distinguishes "Task" as a storable type. "Task" is not (yet) intended for
// use outside of Datatogether servers.
func (t Task) DatastoreType() string {
	return "Task"
}

// GetId returns a task's cannoncial identifier
func (t Task) GetId() string {
	return t.Id
}

// Key is to fulfill the sql_datastore.Model interface
func (t Task) Key() datastore.Key {
	return datastore.NewKey(fmt.Sprintf("%s:%s", t.DatastoreType(), t.GetId()))
}

func (t *Task) PubSubChannelName() string {
	return fmt.Sprintf("tasks.%s", t.Id)
}

// QueueMsg formats the task as an amqp.Publishing message for adding to a queue
func (t *Task) QueueMsg() (amqp.Publishing, error) {
	body, err := json.Marshal(t.Params)
	if err != nil {
		return amqp.Publishing{}, fmt.Errorf("Error marshaling params to JSON: %s", err.Error())
	}

	return amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: t.Id,
		Type:          t.Type,
		UserId:        t.UserId,
		Body:          body,
	}, nil
}

// Enqueue adds a task to the queue located at ampqurl, writing creates/updates
// for the task to the given store
func (task *Task) Enqueue(store datastore.Datastore, amqpurl string) error {
	// Initial save to get an ID, prove we tried to submit
	if err := task.Save(store); err != nil {
		return err
	}

	// connect to queue server & submit task
	conn, err := amqp.Dial(amqpurl)
	if err != nil {
		return fmt.Errorf("Failed to connect to RabbitMQ: %s", err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to connect to open channel: %s", err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"tasks", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("Failed to declare a queue: %s", err.Error())
	}

	msg, err := task.QueueMsg()
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		msg,
	)

	if err != nil {
		return fmt.Errorf("Error publishing to queue: %s", err.Error())
	}

	now := time.Now()
	task.Enqueued = &now
	return task.Save(store)
}

// TaskFromDelivery reads a task from store based on an amqp.Delivery message
func TaskFromDelivery(store datastore.Datastore, msg amqp.Delivery) (*Task, error) {
	t := &Task{Id: msg.CorrelationId}
	if err := t.Read(store); err != nil {
		return nil, err
	}
	return t, nil
}

// Do performs the
func (task *Task) Do(store datastore.Datastore, tc chan *Task) error {
	newTask := taskdefs[task.Type]
	if newTask == nil {
		return fmt.Errorf("unknown task type: %s", task.Type)
	}

	tt := newTask()
	taskBytes, err := json.Marshal(task.Params)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(taskBytes, tt); err != nil {
		return fmt.Errorf("error decoding task body json: %s", err.Error())
	}

	// If the task supports the DatastoreTask interface,
	// pass in our host db connection
	if dsT, ok := tt.(DatastoreTaskable); ok {
		dsT.SetDatastore(store)
	}

	pc := make(chan Progress, 10)

	if err := task.Save(store); err != nil {
		return err
	}

	// execute the task in a goroutine
	go tt.Do(pc)

	for p := range pc {
		// TODO - log progress and pipe out of this func
		// so others can listen in for updates
		// fmt.Println(p.String())
		task.Progress = &p
		tc <- task

		if p.Error != nil {
			task.Error = p.Error.Error()
			now := time.Now()
			task.Failed = &now
			go task.Save(store)
			return p.Error
		}
		if p.Done {
			now := time.Now()
			task.Succeeded = &now
			go task.Save(store)
			return nil
		}
	}

	return nil
}

// StatusString returns a string representation of the status
// of a task based on the state of it's date stamps
func (t *Task) StatusString() string {
	if t.Enqueued == nil {
		return "enquing"
	} else if t.Enqueued != nil {
		return "queued"
	} else if t.Succeeded != nil {
		return "finished"
	} else if t.Failed != nil {
		return "failed"
	} else {
		return "running"
	}
}

func (t *Task) valid() error {
	if taskdefs[t.Type] == nil {
		return fmt.Errorf("unrecognized task type: '%s'", t.Type)
	}

	body, err := json.Marshal(t.Params)
	if err != nil {
		return fmt.Errorf("Error marshaling params to JSON: %s", err.Error())
	}

	// create the task locally to check validity
	// TODO - this should be moved into tasks package?
	tt := taskdefs[t.Type]()
	if err := json.Unmarshal(body, tt); err != nil {
		return fmt.Errorf("Error creating task from JSON: %s", err.Error())
	}

	if err := tt.Valid(); err != nil {
		return fmt.Errorf("Invalid task: %s", err.Error())
	}

	return nil
}

func (t *Task) Read(store datastore.Datastore) error {
	if t.Id == "" {
		return datastore.ErrNotFound
	}

	ti, err := store.Get(t.Key())
	if err != nil {
		return err
	}

	got, ok := ti.(*Task)
	if !ok {
		return fmt.Errorf("Invalid Response")
	}
	*t = *got
	return nil
}

func (t *Task) Save(store datastore.Datastore) (err error) {
	if err := t.valid(); err != nil {
		return err
	}

	var exists bool
	if t.Id != "" {
		exists, err = store.Has(t.Key())
		if err != nil {
			return err
		}
	}

	if !exists {
		t.Id = uuid.New()
		t.Created = time.Now().Round(time.Second).In(time.UTC)
		t.Updated = t.Created
	} else {
		t.Updated = time.Now().Round(time.Second).In(time.UTC)
	}

	return store.Put(t.Key(), t)
}

func (t *Task) Delete(store datastore.Datastore) error {
	return store.Delete(t.Key())
}

func (t *Task) NewSQLModel(key datastore.Key) sql_datastore.Model {
	return &Task{Id: key.Name()}
}

func (t *Task) SQLQuery(cmd sql_datastore.Cmd) string {
	switch cmd {
	case sql_datastore.CmdCreateTable:
		return qTaskCreateTable
	case sql_datastore.CmdExistsOne:
		return qTaskExists
	case sql_datastore.CmdSelectOne:
		return qTaskReadById
	case sql_datastore.CmdInsertOne:
		return qTaskInsert
	case sql_datastore.CmdUpdateOne:
		return qTaskUpdate
	case sql_datastore.CmdDeleteOne:
		return qTaskDelete
	case sql_datastore.CmdList:
		return qTasks
	default:
		return ""
	}
}

func (t *Task) UnmarshalSQL(row sqlutil.Scannable) error {
	var (
		id, title, userId, typ, status, e    string
		paramBytes                           []byte
		params                               map[string]interface{}
		created, updated                     time.Time
		enqueued, started, succeeded, failed *time.Time
	)
	err := row.Scan(
		&id, &created, &updated, &title, &userId, &typ, &paramBytes, &status, &e,
		&enqueued, &started, &succeeded, &failed,
	)
	if err == sql.ErrNoRows {
		return datastore.ErrNotFound
	}

	if paramBytes != nil {
		params = map[string]interface{}{}
		if err := json.Unmarshal(paramBytes, &params); err != nil {
			return err
		}
	}

	*t = Task{
		Id:        id,
		Created:   created,
		Updated:   updated,
		Title:     title,
		UserId:    userId,
		Type:      typ,
		Params:    params,
		Status:    status,
		Error:     e,
		Enqueued:  enqueued,
		Started:   started,
		Succeeded: succeeded,
		Failed:    failed,
	}

	return nil
}

func (t *Task) SQLParams(cmd sql_datastore.Cmd) []interface{} {
	switch cmd {
	case sql_datastore.CmdSelectOne, sql_datastore.CmdExistsOne, sql_datastore.CmdDeleteOne:
		return []interface{}{t.Id}
	case sql_datastore.CmdList:
		return []interface{}{}
	default:
		var params []byte
		if t.Params != nil {
			params, _ = json.Marshal(t.Params)
		}
		return []interface{}{
			t.Id,
			t.Created,
			t.Updated,
			t.Title,
			t.UserId,
			t.Type,
			params,
			t.Status,
			t.Error,
			t.Enqueued,
			t.Started,
			t.Succeeded,
			t.Failed,
			// t.Progress,
		}
	}
}
