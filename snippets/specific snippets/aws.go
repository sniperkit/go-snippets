/*
This is an example config for connecting to your AWS account

Your AWS credentials need to be saved in 'aws.text' (or adjust code)
It only needs to contain this (These are example credentials):
`
ADJKNSKWNNIDMSMSDKSD
chjKmsjeKSYNmalOQksmeJSUWJWMAl

`
Be sure to exclude your text file from your github repo, add `aws.txt` to .gitignore!
*/
package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func main() {
	// Read Credentials from file
	content, err := ioutil.ReadFile("aws.txt")
	if err != nil {
		fmt.Println("Couldn't access file")
	}

	lines := strings.Split(string(content), "\n")
	id := lines[0]
	key := lines[1]

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		//Credentials: credentials.NewSharedCredentials("AWS_PROFILE", "temp")
		Credentials: credentials.NewStaticCredentials(id, key, ""),
	})
	if err != nil {
		fmt.Println("Couldn't build new session")
	}

	_, err := sess.Config.Credentials.Get()
	if err != nil {
		fmt.Println("Couldn't find credentials")
	}
}
