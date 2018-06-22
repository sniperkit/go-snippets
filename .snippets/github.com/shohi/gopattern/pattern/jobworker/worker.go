package jobworker

import (
	"io/ioutil"
	"net/http"
)

// NewWorker - create a new worker
func NewWorker(workerPool chan chan job, resultCh chan *result) Worker {
	return Worker{
		workerPool:    workerPool,
		jobChannel:    make(chan job),
		resultChannel: resultCh,
		quit:          make(chan bool)}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) Start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				// we have received a work request.
				// download
				var res = result{resourceURL: job.resourceURL}

				res.data, res.err = w.download(job.resourceURL)
				w.resultChannel <- &res

			case <-w.quit:
				// we have received a signal to stop
				return
			}
		}
	}()
}

func (w Worker) download(resourceURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", resourceURL, nil)
	if err != nil {
		return nil, err
	}

	//
	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
