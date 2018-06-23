// Copyright 2015 Gavin Bong. All rights reserved.
// Use of this source code is governed by the Hugsware
// license that can be found in the LICENSE file.

// Package charlieswithcrayons provides access
// to random.org's JSON RPC api.
package charlieswithcrayons

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
    "errors"
)

const Version = "1.0.0"
const StdUrl = "https://api.random.org/json-rpc/1/invoke"
const methodBlob = "generateBlobs"
const methodBlobSigned = "generateSignedBlobs"


type Arguments struct {
	ApiKey string
	Url    string
}

type RandomOrgClient struct {
	apiKey string
	url    string
}

func New(args Arguments) *RandomOrgClient {
	client := new(RandomOrgClient)
	client.apiKey = args.ApiKey
	client.url = args.Url
	return client
}

type BlobParams struct {
	ApiKey string `json:"apiKey"`
	N      int    `json:"n"`
	Size   int    `json:"size"`
}

type BlobArguments struct {
	JsonRpc    string `json:"jsonrpc"`
	Method     string `json:"method"`
	BlobParams `json:"params"`
	Id         int `json:"id"`
}

type Blob struct {
	Data           []string `json:"data"`
	CompletionTime string   `json:"completionTime"`
}

type BlobWrapper struct {
	Blob          `json:"random"`
	BitsUsed      int `json:"bitsUsed"`
	BitsLeft      int `json:"bitsLeft"`
	RequestsLeft  int `json:"requestsLeft"`
	AdvisoryDelay int `json:"advisoryDelay"`
}

type BlobError struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type BlobResponse struct {
	JsonRpc     string `json:"jsonrpc"`
	BlobWrapper `json:"result"`
	BlobError   `json:"error"`
	Id          int `json:"id"`
}

func (c *RandomOrgClient) newBlobArguments(method string, bitSize int, txId int) *BlobArguments {
	retval := new(BlobArguments)
	retval.JsonRpc = "2.0"
	retval.Method = method
	retval.Id = txId
	retval.BlobParams.ApiKey = c.apiKey
	retval.BlobParams.N = 1
	retval.BlobParams.Size = bitSize
	return retval
}

// You must have a valid Api key to invoke this method.
// Argument `txId` can be used to differentiate calls.
func (c *RandomOrgClient) GetRandomBits(bitSize int, count int, txId int) ([]string, error) {
	log.Println("invoke generateBlobs")
	blobArguments := c.newBlobArguments(methodBlob, bitSize, txId)
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	postBody, err := json.Marshal(blobArguments)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

    var blob BlobResponse
    if err = json.Unmarshal(body, &blob); err != nil {
       log.Println("\terror:", err)
       return nil, err
    }
    if len(blob.BlobError.Message) > 0 {
       return nil, errors.New(blob.BlobError.Message)
    } else {
       return blob.BlobWrapper.Blob.Data, nil
    }
}

// @TODO
func (c *RandomOrgClient) GetSignedRandom(bitSize int, count int) error {
	log.Println("invoke generateSignedBlobs")
	return nil
}
