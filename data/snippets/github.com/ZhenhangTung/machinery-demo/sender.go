package main

import (
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/RichardKnop/machinery/v1/log"
	"fmt"
	"time"
)

func sendSms() error {
	server, err := startServer()
	if err != nil {
		return err
	}

	var (
		singleMsgTask tasks.Signature
		mutiMsgTask, mutiMsgTask1  tasks.Signature
		longRunningTask tasks.Signature
		delayedTask tasks.Signature
		onErrorTask tasks.Signature
	)

	var initTasks = func() {
		onErrorTask = tasks.Signature{
			Name: "on_error_task",
		}

		singleMsgTask = tasks.Signature{
			Name: "single_sms",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: "15021788888",
				},
				{
					Type:  "string",
					Value: "This is a message",
				},
			},
			OnError: []*tasks.Signature{&onErrorTask},
		}

		mutiMsgTask = tasks.Signature{
			Name: "multiple_sms",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: "15021788888",
				},
				{
					Type:  "string",
					Value: "This is a multi message",
				},
			},
		}

		mutiMsgTask1 = tasks.Signature{
			Name: "multiple_sms",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: "150217777777",
				},
				{
					Type:  "string",
					Value: "This is a multi message",
				},
			},
		}

		longRunningTask = tasks.Signature{
			Name: "long_running_task",
		}

		delayedTask = tasks.Signature{
			Name: "delayed_task",
		}
	}

	/*
	 * First, let's try sending a single task
	 */
	initTasks()
	log.INFO.Println("Single task:")

	asyncResult, err := server.SendTask(&singleMsgTask)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf("Sending result: %v\n", tasks.HumanReadableResults(results))

	log.INFO.Printf("==================")

	log.INFO.Println("Group of tasks (parallel execution):")
	var groupTasks []*tasks.Signature
	// 1000 tasks in parallel
	for i := 0; i < 1000; i++ {
		groupTasks = append(groupTasks, &mutiMsgTask)
	}
	group := tasks.NewGroup(groupTasks...)
	asyncResults, err := server.SendGroup(group, 10)
	if err != nil {
		return fmt.Errorf("Could not send group: %s", err.Error())
	}

	for _, asyncResult := range asyncResults {
		results, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
		if err != nil {
			return fmt.Errorf("Getting task result failed with error: %s", err.Error())
		}
		log.INFO.Printf(
			"Sending result: %v\n", tasks.HumanReadableResults(results),
		)
	}

	log.INFO.Printf("==================")
	// Let's try a long running task
	asyncResult, err = server.SendTask(&longRunningTask)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return fmt.Errorf("Getting long running task result failed with error: %s", err.Error())
	}
	log.INFO.Printf("Long running task returned = %v\n", tasks.HumanReadableResults(results))


	log.INFO.Printf("ddddddddddddddddddddddddd")
	// Let's try a delayed running task
	eta := time.Now().UTC().Add(time.Second * 10)
	delayedTask.ETA = &eta
	asyncResult, err = server.SendTask(&delayedTask)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err = asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return fmt.Errorf("Getting delayed running task result failed with error: %s", err.Error())
	}
	log.INFO.Printf("delayed running task returned = %v\n", tasks.HumanReadableResults(results))

	return nil
}
