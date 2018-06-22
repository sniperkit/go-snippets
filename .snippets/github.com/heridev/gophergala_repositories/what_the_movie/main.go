package main

import (
	log `github.com/cihub/seelog`
	mgo `gopkg.in/mgo.v2`
	`os`
)

const (
	JOBS_COLLECTION   = `jobs`
	FRAMES_COLLECTION = `frames`
	MOVIES_COLLECTION = `movies`
)

func main() {
	defer log.Flush()
	if 3 != len(os.Args) {
		log.Errorf(`Usage: %s dbhost dbname`, os.Args[0])
		return
	}

	// Database connection.
	log.Infof(`Connecting to database %s at address %s.`, os.Args[2], os.Args[1])
	mgoSession, err := mgo.Dial(os.Args[1])
	if err != nil {
		log.Errorf(`Error connecting to MongoDB: %v.`, err)
		return
	}
	log.Infof(`Connected!`)
	runProcesses(mgoSession, os.Args[2])

	<-(chan struct{})(nil)
}
