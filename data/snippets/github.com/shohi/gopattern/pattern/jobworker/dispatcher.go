package jobworker

import (
	"log"
	"net/http"

	"github.com/allegro/bigcache"
)

// HTTP handlers send requests to `Dispatcher`, `Dispatcher` wrappers request to `Job`
// `Dispatcher` then assign it to a worker

// RequestQueue - a buffered channel that hold
var RequestQueue chan *RequestContext
var defaultRequestQueueSize = 1024

var defaultCache *bigcache.BigCache

// requestGroupMap - group requests by its URL
var requestGroupMap map[string][]*RequestContext

// RequestContext - request context wrappers requests and response
type RequestContext struct {
	ReqesutURL string
	Response   chan *Response
}

// Response - response that request wants
type Response struct {
	Err  error
	Data []byte
}

// jobQueue - A buffered channel that we can send work jobs on.
var jobQueue chan job
var defaultJobQueueSize = 512

// job represents the job to be run
type job struct {
	resourceURL string
}

type result struct {
	resourceURL string
	err         error
	data        []byte
}

// Worker represents the worker that executes the job
type Worker struct {
	client        http.Client
	workerPool    chan chan job
	jobChannel    chan job
	resultChannel chan *result
	quit          chan bool
}

// Dispatcher dispatch all requests to workers
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	workerPool    chan chan job
	maxWorkers    int
	resultChannel chan *result
	jobCount      int // created jobs count
	cache         *bigcache.BigCache
}

func init() {
	// init ReqeustQueue
	RequestQueue = make(chan *RequestContext, defaultRequestQueueSize)
	requestGroupMap = make(map[string][]*RequestContext)

	// init jobQueue
	jobQueue = make(chan job, defaultJobQueueSize)

	cacheConfig := bigcache.DefaultConfig(0)
	cacheConfig.MaxEntriesInWindow = 1024
	cache, err := bigcache.NewBigCache(cacheConfig)
	if err != nil {
		log.Printf("failed to init local cache, err: %v\n", err)
	} else {
		defaultCache = cache
	}

	log.Println("defaultCache ==> ", defaultCache != nil)

}

// NewDispatcher - create a new dispatcher
func NewDispatcher(maxWorkers int, cache *bigcache.BigCache) *Dispatcher {
	pool := make(chan chan job, maxWorkers)
	dispatcher := Dispatcher{
		workerPool:    pool,
		maxWorkers:    maxWorkers,
		resultChannel: make(chan *result, maxWorkers),
		cache:         cache,
	}
	return &dispatcher
}

// Run - run dispatch
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.workerPool, d.resultChannel)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		// handle requests
		case reqCtx := <-RequestQueue:

			reqURL := reqCtx.ReqesutURL

			// Check Cache first
			if d.cache != nil {
				if data, err := d.cache.Get(reqURL); err == nil {
					reqCtx.Response <- &Response{Data: data, Err: nil}
					continue
				}
			}
			// check if there is a job running
			group, ok := requestGroupMap[reqURL]
			if ok {
				requestGroupMap[reqURL] = append(group, reqCtx)
				continue
			}

			// add a new job
			requestGroupMap[reqURL] = []*RequestContext{reqCtx}
			d.jobCount++

			// a job request has been received
			go func(j job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.workerPool

				// dispatch the job to the worker job channel
				jobChannel <- j
			}(job{resourceURL: reqURL})

		// handle response
		case result := <-d.resultChannel:
			group, ok := requestGroupMap[result.resourceURL]
			if !ok {
				continue
			}

			response := &Response{Data: result.data, Err: result.err}
			for _, reqCtx := range group {
				reqCtx.Response <- response
			}

			delete(requestGroupMap, result.resourceURL)

			// save to cache if no error occurs
			if result.err == nil && d.cache != nil {
				d.cache.Set(result.resourceURL, result.data)
			}
		}
	}
}

// TotalJobs - number of jobs created for downloading
func (d *Dispatcher) TotalJobs() int {
	return d.jobCount
}
