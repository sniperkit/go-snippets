package sync_test

import (
	"fmt"
	"testing"
	"time"
	"distributed-programming-abstractions/job_handlers"
	jh "distributed-programming-abstractions/job_handlers/sync"
)

// Fulfill the properties.
func TestGuaranteedResponse(t *testing.T) {
	// fmt.Println("Testing GuaranteedResponse")
	job_channel := make(chan job.Job)
	jh := jh.New()
	jh.Confirm(func(j job.Job) {
		fmt.Println("Confirmed job: ", j)
		job_channel <- j
	})
	// For the second error, i.e: The submitted job was never Confirmed, uncomment the following function.
	/*
	jh.Process(func(j job.Job) {
		time.Sleep(100*time.Millisecond)
	})*/
	jh.Submit("Job 1")
	select {
		case job_ := <-job_channel:
			value, status := job_.(string)
			if !status || value != "Job 1" {
				t.Errorf("The Confirmed job was never submitted.")
			}
		case <-time.After(100 * time.Millisecond): // Timer functionality.
			t.Errorf("The submitted job was never Confirmed.")
	}
}

func TestRun(t *testing.T) {
	// fmt.Println("Testing testrun")
	var msg = make(chan string)
	jh := jh.New()
	jh.Confirm(func(j job.Job) {
		msg <- fmt.Sprintln("Confirmed job: ", j)
	})

	// Setting process to test whether processing takes place or not!.
	jh.Process(func(j job.Job) {
		fmt.Println("Processing job: ", j)
	})
	jh.Submit("Job 1")
	jh.Submit("Job 2")
	jh.Submit("Job 3")

	fmt.Println(<-msg)
	fmt.Println(<-msg)
	fmt.Println(<-msg)
}