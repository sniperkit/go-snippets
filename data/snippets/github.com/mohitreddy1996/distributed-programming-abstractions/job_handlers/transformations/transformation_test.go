package transformation_test

import (
	"fmt"
	"testing"
	"time"
	"distributed-programming-abstractions/job_handlers"
	"distributed-programming-abstractions/job_handlers/transformations"
	"distributed-programming-abstractions/job_handlers/sync"
)

func TestGuaranteedResponse(t *testing.T) {
	confirm_channel := make(chan job.Job)
	job_handler := sync.New()
	transformation_handler := transformation.New(job_handler, 1)

	transformation_handler.Confirm(func(j job.Job) {
		fmt.Println("Confirmed job: ", j)
		confirm_channel <- j
	})

	transformation_handler.Submit("job1")
	select {
		case c := <-confirm_channel:
			content, ok := c.(string) 
			if !ok || content != "job1" {
				t.Errorf("Job was never submitted.")
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Timeout after 100 milliseconds.")
	}
}

func TestRunProcess(t *testing.T) {
	process_channel := make(chan job.Job)
	job_handler := sync.New()
	transformation_handler := transformation.New(job_handler, 1)
	transformation_handler.Process(func(j job.Job) {
		fmt.Println("Running process: ", j)
		process_channel <- j
	})
	transformation_handler.Submit("job1")
	select {
		case c := <-process_channel:
			content, ok := c.(string) 
			if !ok || content != "job1" {
				t.Errorf("job was never submitted.")
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("timeout after 100 milliseconds.")
	}
}

func TestSoundness(t *testing.T) {
	job_handler := sync.New()
	transformation_handler := transformation.New(job_handler, 1)

	process_channel := make(chan job.Job)
	error_channel := make(chan job.Job)
	wait_channel := make(chan job.Job)

	transformation_handler.Process(func(j job.Job) {
		process_channel <- j
		fmt.Println("Processing job: ", j)
		<-wait_channel
	})

	transformation_handler.Error(func(j job.Job) {
		error_channel <- j
	})

	transformation_handler.Confirm(func(j job.Job) {})

	transformation_handler.Submit("job1")
	<-process_channel
	transformation_handler.Submit("job2")
	transformation_handler.Submit("job3")

	select {
		case j := <-error_channel:
			content, ok := j.(string)
			if !ok || content != "job3" {
				t.Errorf("Error in the Submit(job.Job) code flow.")
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Timeout in the desired operation.")
	}
	close(wait_channel)
}