// +build ec2

package tron

import (
	"time"

	"github.com/golang/glog"

	"github.com/gophergala/tron/aws"
)

var (
	AppID string
)

func init() {
	AppID = getAppID()
}

func getAppID() string {
	var id string
	var err error
	retries := 0
	backoff := 1
	for {
		id, err = aws.EbEnvID()
		if err == nil || retries > 10 {
			break
		}

		<-time.After(time.Duration(backoff) * time.Second)
		if backoff < 60 {
			backoff *= 2
		}
		retries += 1
	}

	if err != nil {
		glog.Fatalf("error getting elasticbeanstalk environment ID: %v", err)
	}
	return id
}
