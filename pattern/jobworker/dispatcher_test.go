package jobworker

import (
	"log"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDispatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dispatcher suite")
}

var _ = Describe("Dispatcher for requesting same resources multiple times", func() {
	var dispatcher = NewDispatcher(10, defaultCache)

	Context("when multiple requests for same resources arrive", func() {
	})

	It("there should be one real request sent", func() {
		var requestsCount = 10
		var resourceURL = "www.douban.com"
		var wg sync.WaitGroup

		dispatcher.Run()

		for k := 0; k < requestsCount; k++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				reqCtx := &RequestContext{
					ReqesutURL: resourceURL,
					Response:   make(chan *Response),
				}
				RequestQueue <- reqCtx
			}()
		}

		wg.Wait()
		Expect(dispatcher.TotalJobs()).To(Equal(1))
	})
})

func TestSameRequests(t *testing.T) {
	var dispatcher = NewDispatcher(10, defaultCache)
	var requestsCount = 10
	var resourceURL = "http://www.baidu.com"
	var wg sync.WaitGroup

	dispatcher.Run()

	for k := 0; k < requestsCount; k++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			if index == 5 {
				time.Sleep(3 * time.Second)
			}
			reqCtx := &RequestContext{
				ReqesutURL: resourceURL,
				Response:   make(chan *Response),
			}

			log.Println(index, " ===> ", reqCtx)

			// add to the queue
			RequestQueue <- reqCtx

			// wait result
			resp := <-reqCtx.Response
			log.Println("resp ==> ", resp.Err, " data length ==> ", len(resp.Data))
			// log.Println(string(resp.Data))

			log.Printf("request %d finished!", index)
		}(k)
	}

	wg.Wait()

	log.Println(dispatcher.TotalJobs())
}
