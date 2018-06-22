package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"strings"
)

var accessKey, secretKey string
var auth aws.Auth

func init() {
	key, err := ioutil.ReadFile("accessKey")
	if err != nil {
		fmt.Println("Lol have fun uploading images without an AWS key", err)
	}
	accessKey = strings.TrimSpace(string(key))

	key, err = ioutil.ReadFile("secretKey")
	if err != nil {
		fmt.Println("Lol have fun uploading images without an AWS key", err)
	}
	secretKey = strings.TrimSpace(string(key))

	auth = aws.Auth{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
}

func AddWallImage(wallName, fileName string, buf *bytes.Buffer) error {
	connection := s3.New(auth, aws.USEast)
	picBucket := connection.Bucket("gollage/" + wallName + "")
	return picBucket.PutReader(fileName+".png", buf, int64(buf.Len()), "image/png", s3.PublicRead)
}

func NewWallBucket(name string) error {
	//connection := s3.New(auth, aws.USEast)
	//picBucket := connection.Bucket("gollage/" + name)
	//return picBucket.PutBucket(s3.PublicRead)
	return nil
}
