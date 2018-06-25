//
// storage package wraps both local disc and S3 storage with CRUD operations and common os.*, ioutil.* functions.
//
// @author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2017, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package storage

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// Version of package
	Version = `v0.0.2`

	// UNSUPPORTED is used if Init fails the and a non-nil Store struck gets returned.
	UNSUPPORTED = iota
	// FS local file system
	FS
	// S3 remote storage via AWS S3
	S3
	// GS remote storage via Google Cloud Storage
	GS
	// Other constants will be create as other storage systems are implemented
)

// Store wrapps the given system interface normalizing to simple Create, Read, Update, Delete operations
type Store struct {
	// Attributes holds any data needed for managing the remote session for desired operations
	Config map[string]interface{}
	Type   int

	// Basic CRUD Operations
	Create func(string, io.Reader) error
	Read   func(string) ([]byte, error)
	Update func(string, io.Reader) error
	Delete func(string) error

	// Additional operations for compatibility with os.* and ioutil.*
	Stat      func(string) (os.FileInfo, error)
	Mkdir     func(string, os.FileMode) error
	MkdirAll  func(string, os.FileMode) error
	Remove    func(string) error
	RemoveAll func(string) error
	ReadFile  func(string) ([]byte, error)
	WriteFile func(string, []byte, os.FileMode) error

	// Extended operations for datatools and dataset
	// Writefilter takes a final path and a processing function which accepts the temp pointer
	WriteFilter func(string, func(*os.File) error) error
}

// EnvToOptions given an environment map envvars to their option.
func EnvToOptions(env []string) map[string]interface{} {
	opts := map[string]interface{}{}
	for _, stmt := range env {
		switch {
		case strings.HasPrefix(stmt, "AWS_") && strings.Contains(stmt, "="):
			kv := strings.SplitN(stmt, "=", 2)
			switch kv[0] {
			case "AWS_SDK_LOAD_CONFIG":
				if kv[0] == "1" || strings.ToLower(kv[1]) == "true" {
					opts["AwsSDKLoadConfig"] = true
				} else {
					opts["AwsSDKLoadConfig"] = false
				}
			case "AWS_PROFILE":
				if kv[0] != "" {
					opts["AwsProfile"] = kv[1]
				} else {
					opts["AwsProfile"] = "default"
				}
			case "AWS_SHARED_CONFIG_ENABLED":
				if kv[0] == "1" || strings.ToLower(kv[1]) == "true" {
					opts["AwsSharedConfigEnabled"] = true
				} else {
					opts["AwsSharedConfigEnabled"] = false
				}
			case "AWS_BUCKET":
				opts["AwsBucket"] = kv[1]
			}
		case (strings.HasPrefix(stmt, "GOOGLE_")) && strings.Contains(stmt, "="):
			kv := strings.SplitN(stmt, "=", 2)
			switch kv[0] {
			case "GOOGLE_BUCKET":
				opts["GoogleBucket"] = kv[1]
			}
		}
	}
	return opts
}

// Init returns a *Store structure that points to configuration info (e.g. S3 credentials)
// and basic CRUD functions associated with the Store's storage type.
func Init(storeType int, options map[string]interface{}) (*Store, error) {
	store := new(Store)
	store.Type = UNSUPPORTED
	// Copy the values into store.Config so they can travel with the struct
	store.Config = map[string]interface{}{}
	for key, val := range options {
		store.Config[key] = val
	}
	switch storeType {
	case FS:
		return fsConfigure(store)
	case S3:
		return s3Configure(store)
	case GS:
		return gsConfigure(store)
	default:
		return store, fmt.Errorf("storeType not supported")
	}
}

// GetDefaultStore tries to guess the storage type based on environment settings
// if it can't is assumes storage.FS type.
func GetDefaultStore() *Store {
	opts := map[string]interface{}{}
	sType := FS
	for _, env := range os.Environ() {
		if strings.Contains(env, "=") {
			kv := strings.SplitN(env, "=", 2)
			if len(kv) == 2 {
				k, v := kv[0], kv[1]
				opts[k] = v
				if strings.HasPrefix(k, "AWS_") == true {
					sType = S3
				}
			}
		}
	}
	//FIXME: Shouldn't we be valling individual typed default functions per sType? (e.g. in fs.go, s3.go, gs.go)
	switch sType {
	case S3:
		if s := os.Getenv("AWS_BUCKET"); s != "" {
			opts["AwsBucket"] = s
		}
		if s := os.Getenv("AWS_SDK_LOAD_CONFIG"); s == "1" || strings.ToLower(s) == "true" {
			opts["AwsSDKLoadConfig"] = true
			opts["AwsSharedConfigEnabled"] = true
			if t := os.Getenv("AWS_PROFILE"); t != "" {
				opts["AwsProfile"] = t
			} else {
				opts["AwsProfile"] = "default"
			}
		}
	case GS:
		if s := os.Getenv("GOOGLE_PROJECT_ID"); s != "" {
			opts["GoogleProjectID"] = s
		}
		if s := os.Getenv("GOOGLE_BUCKECT"); s != "" {
			opts["GoogleBucket"] = s
		}
	}
	store, _ := Init(sType, opts)
	return store
}
