/*
* @Author: souravray
* @Date:   2015-01-25 04:12:39
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 05:53:40
 */

package stacker

import (
	"sync"
)

type Queue struct {
	jobs  []*Job
	mutex sync.Mutex
}

func (q *Queue) Push(job *Job) {
	q.mutex.Lock()
	q.jobs = append(q.jobs, job)
	q.mutex.Unlock()
}

func (q *Queue) Pop() (job *Job) {
	q.mutex.Lock()
	old := q.jobs
	n := len(old)
	if n > 1 {
		job = old[0]
		q.jobs = old[1:n]
	} else if n == 1 {
		job = old[0]
		q.jobs = make([]*Job, 0)
	}
	q.mutex.Unlock()
	return
}

func NewQueue() (q Queue) {
	q = Queue{jobs: make([]*Job, 0)}
	return
}
