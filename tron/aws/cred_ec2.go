// +build ec2

package aws

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/smartystreets/go-aws-auth"
)

const (
	credentialURL = "http://169.254.169.254/latest/meta-data/iam/security-credentials/"
)

var (
	cred = struct {
		sync.RWMutex
		c awsauth.Credentials
	}{}
	Region = detectRegionMust()
)

func init() {
	c := refreshCredentials()
	go func() {
		for {
			<-time.After(c.Expiration.Sub(time.Now()) - 15*time.Minute)
			c = refreshCredentials()
			glog.Infof("refreshed aws credentials: %+v", c)
		}
	}()
}

func Credentials() awsauth.Credentials {
	cred.RLock()
	c := cred.c
	cred.RUnlock()
	return c
}

func refreshCredentials() awsauth.Credentials {
	var c *awsauth.Credentials
	var err error
	backoff := 1
	for backoff < 20 {
		c, err = queryMetadata()
		if err == nil {
			break
		}
		<-time.After(time.Duration(backoff) * time.Second)
		backoff *= 2
	}

	if err != nil {
		glog.Fatalf("instance metadata error: %v", err)
	}
	cred.Lock()
	cred.c = *c
	cred.Unlock()
	return *c
}

func queryMetadata() (*awsauth.Credentials, error) {
	roleResp, err := http.Get(credentialURL)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(roleResp.Body)
	if !scanner.Scan() {
		return nil, fmt.Errorf("no role in instance metadata")
	}
	role := scanner.Text()
	roleResp.Body.Close()

	credResp, err := http.Get(credentialURL + role)
	if err != nil {
		return nil, err
	}
	defer credResp.Body.Close()
	b, err := ioutil.ReadAll(credResp.Body)
	if err != nil {
		return nil, err
	}
	c := &awsauth.Credentials{}
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, fmt.Errorf(`json error "%v" for body: %s`, err, string(b))
	}
	return c, nil
}
