// +build local

package aws

import (
	"os"

	"github.com/smartystreets/go-aws-auth"
)

var (
	c = awsauth.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
	Region = ""
)

func init() {
	c.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	c.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
}

func Credentials() awsauth.Credentials {
	return c
}
