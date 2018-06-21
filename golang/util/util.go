// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// Package util provides shared utilities for all gopacket examples.
package util

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime/pprof"
)

var cpuprofile = flag.String("cpuprofile", "", "Where to write CPU profile")

// Run starts up stuff at the beginning of a main function, and returns a
// function to defer until the function completes.  It should be used like this:
//
//   func main() {
//     defer util.Run()()
//     ... stuff ...
//   }
func Run() func() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("could not open cpu profile file %q", *cpuprofile)
		}
		pprof.StartCPUProfile(f)
		return func() {
			pprof.StopCPUProfile()
			f.Close()
		}
	}
	return func() {}
}

// Run a command and wait to be finished.
//func main() {
//	wait("ls", "-al")
//	wait("ps", "aelf")
//}
func Wait(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

//Save data into a file -- overwrite
func SaveToFile(name string, data []byte) {
	file, err := os.Create(name)
	if err != nil {
		log.Println("Cannot create file: ", name, " ", err.Error())
		return
	}

	file.Write(data)
	file.Sync()
	defer file.Close()
}

//Save data into a file. -- append
func AppendToFile(name string, data []byte) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(data)
}

//Download link and write to name
func Download(link, name string) error {
	if name == "" || link == "" {
		return errors.New("Invalid paramenters: " + link + ":" + name)
	}

	resp, err := http.Get(link)
	if err != nil {
		return errors.New("Do http request error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Do http request error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Read http request body error: " + err.Error())
	}

	if err := ioutil.WriteFile(name, body, 0644); err != nil {
		return errors.New("Write content to file error: " + err.Error())
	}

	return nil
}

//Generate a string of random number
func GenerateSessionID() string {
	var id = "0x"

	rand.Seed(time.Now().Unix())
	result := rand.Perm(128)
	for _, i := range result {
		id = id + strconv.Itoa(i)
	}

	return id
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func CloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}
