package main

import (
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/errors"
	"log"
)

var (
	server *machinery.Server
	worker *machinery.Worker

	err error
)

func init() {
	cnf := config.Config{
		Broker:        "redis://127.0.0.1:6379",
		ResultBackend: "redis://127.0.0.1:6379",
		Exchange:      "machinery_exchange",
		ExchangeType:  "direct",
		DefaultQueue:  "machinery_tasks",
		BindingKey:    "machinery_task",
	}

	server, err = machinery.NewServer(&cnf)
	if err != nil {
		errors.Fail(err, "Could not initialize server")
	}

	tasks := map[string]interface{}{
		"add": Add,
	}
	server.RegisterTasks(tasks)

	worker = server.NewWorker("machinery_worker")
}

func main() {
	err := worker.Launch()
	if err != nil {
		errors.Fail(err, "Could not launch worker")
	}
}

func Add(args ...int64) (int64, error) {
	sum := int64(0)
	for _, arg := range args {
		sum += arg
	}

	log.Println("lalalala")

	return sum, nil
}
