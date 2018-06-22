/*
* @Author: souravray
* @Date:   2015-01-25 03:49:44
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 05:45:47
 */
package stacker

import (
	w "github.com/gophergala/tinyembassy/webservice/stacker/worker"
)

type Job struct {
	// Worker payload.
	Payload w.Payload

	// The number of times the task has been dispatched
	RetryCount int
}
