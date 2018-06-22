/*
* @Author: souravray
* @Date:   2015-01-25 03:28:47
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 09:29:43
 */

package stacker

import (
	"github.com/gophergala/tinyembassy/webservice/stacker/worker"
	"math"
)

type Stacker struct {
	Queue
	bucket    *Bucket
	churnRate int32
	stop      chan bool
}

func GetStacker(maxConcurrentWorker int32, maxDequeueRate int32) (*Stacker, error) {
	bucket, err := NewBucket(maxConcurrentWorker, maxDequeueRate)
	if err != nil {
		return nil, err
	}
	churnRate := int32(math.Ceil(float64(maxDequeueRate/3)) * 2)

	stacker := &Stacker{
		NewQueue(),
		bucket,
		churnRate,
		make(chan bool)}

	return stacker, nil
}

func (s *Stacker) Start() {
	go s.bucket.Fill()
	for {
		select {
		case <-s.stop:
			s.bucket.Close()
			return
		default:
			n := <-s.bucket.Take(s.churnRate)
			for i := int32(0); i < n; i++ {
				job := s.Pop()
				if job != nil {
					go Shoveler(s, job)
					// go func(s *Stacker, w W.Interface, payload url.Values) {
					// 	defer s.bucket.Spend()
					// 	err := w.Perform(payload)
					// 	if err != nil {

					// 	} else {

					// 	}
					// }(s, w.Interface, item.Payload)
				}
				s.bucket.Spend()
			}
		}
	}
}

func (s *Stacker) Delete() bool {
	s.bucket.Close()
	s = nil
	return true
}

func (s *Stacker) AddJob(payload worker.Payload) {
	job := &Job{payload, 0}
	s.Push(job)
}

func (s *Stacker) retryJob(job *Job) {
	job.RetryCount++
	s.Push(job)
}

func Shoveler(s *Stacker, job *Job) {
	worker.Z(job.Payload)
}
