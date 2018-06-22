package aws

import (
	"bufio"
	"fmt"
	"launchpad.net/goamz/aws"
	// "io"
	"launchpad.net/goamz/s3"
	"net/http"
	"os"
	"strings"
)

//TODO: use this from config file.
const (
	accesskey  = "AKIAIK7FZYOKWK5JXDEA"
	secretKey  = "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"
	region     = ""
	bucketname = "mazibucket"
	s3Path     = "imageData"
)

type FileUpload struct {
}

func (this *FileUpload) FileUploadFromPath(localPath string, campaign string, badgeGrp string, badge string) (err error, url string) {

	AWSAuth := aws.Auth{
		AccessKey: accesskey,
		SecretKey: secretKey,
	}

	//Create connection
	connection := s3.New(AWSAuth, aws.USEast)

	//Initilize bucket
	bucket := connection.Bucket(bucketname)

	//Path for images
	s := []string{s3Path, campaign, badgeGrp, badge}
	// path := append(s3Path, "/", campaign, "/", badge)
	path := strings.Join(s, "/")
	fmt.Println(path)

	file, err := os.Open(localPath)

	if err != nil {
		fmt.Println(err)
		return err, ""
	}

	defer file.Close()

	//get file size
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	bytes := make([]byte, size)

	//Prepare buffer to post
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	//get file type
	filetype := http.DetectContentType(bytes)

	//upload file to S3
	err = bucket.Put(path, bytes, filetype, s3.ACL("public-read"))

	if err != nil {
		fmt.Println(err)
		return err, ""
	}
	s = []string{"https://s3.amazonaws.com", bucketname, s3Path, campaign, badgeGrp, badge}
	return nil, strings.Join(s, "/")
}

func (this *FileUpload) UploadToS3(data []byte, campaignName string, badge string) (err error, url string) {
	AWSAuth := aws.Auth{
		AccessKey: accesskey,
		SecretKey: secretKey,
	}
	fmt.Println("upload to s3:" + campaignName + badge)
	//Create connection
	connection := s3.New(AWSAuth, aws.USEast)

	//Initilize bucket
	bucket := connection.Bucket(bucketname)

	//Path for images
	s := []string{s3Path, campaignName, badge}
	// path := append(s3Path, "/", campaign, "/", badge)
	path := strings.Join(s, "/")

	//get file type
	filetype := http.DetectContentType(data)

	//upload file to S3
	err = bucket.Put(path, data, filetype, s3.ACL("public-read"))

	if err != nil {
		fmt.Println(err)
		return err, ""
	}

	s = []string{"https://s3.amazonaws.com", bucketname, s3Path, campaignName, badge}
	return nil, strings.Join(s, "/")
}
