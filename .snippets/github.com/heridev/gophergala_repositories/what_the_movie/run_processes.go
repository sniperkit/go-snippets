package main

import (
	log `github.com/cihub/seelog`
	`github.com/gorilla/mux`
	mgo `gopkg.in/mgo.v2`
	`gopkg.in/mgo.v2/bson`
	`net/http`
	`time`
)

const (
	N_FRAME_THREADS = 10
	N_MOVIE_THREADS = 5
	LISTEN_ADDRESS  = `:2302`
)

func runProcesses(dbConn *mgo.Session, dbName string) {
	runHttpServer(dbConn, dbName)
	frames := runFrameProcessors(dbConn, dbName)
	movies := runMovieProcessors(frames, dbConn, dbName)
	runMovieJobProcessor(movies, dbConn, dbName)
}

func runFrameProcessors(dbConn *mgo.Session, dbName string) chan Frame {
	frames := make(chan Frame, N_FRAME_THREADS)
	for i := 0; i < N_FRAME_THREADS; i++ {
		go processFrames(frames, dbConn.Copy(), dbName)
	}
	return frames
}

func runMovieProcessors(frames chan Frame, dbConn *mgo.Session, dbName string) chan MovieProcessJob {
	movies := make(chan MovieProcessJob, N_MOVIE_THREADS)
	for i := 0; i < N_MOVIE_THREADS; i++ {
		go processMovies(movies, frames, dbConn.Copy(), dbName)
	}
	return movies
}

func runMovieJobProcessor(movies chan MovieProcessJob, dbConn *mgo.Session, dbName string) {
	go func() {
		db := dbConn.DB(dbName)
		for {
			log.Infof(`Querying...`)
			jobsColl := db.C(JOBS_COLLECTION)
			jobsIter := jobsColl.Find(bson.M{`processed`: false}).Tail(-1)
			var movieJob MovieProcessJob
			for jobsIter.Next(&movieJob) {
				movies <- movieJob
			}
			if err := jobsIter.Err(); nil != err {
				log.Errorf(`Error: %v.`, err)
			}
			jobsIter.Close()
			time.Sleep(time.Second)
		}
	}()
}

func runHttpServer(dbConn *mgo.Session, dbName string) {
	func() {
		r := mux.NewRouter()
		r.HandleFunc(`/`, homeHandler).Methods(`GET`)
		r.HandleFunc(`/search`, searchHandler(dbConn, dbName)).Methods(`POST`)
		http.Handle(`/`, r)

		log.Infof(`Listening HTTP at address %s.`, LISTEN_ADDRESS)
		http.ListenAndServe(LISTEN_ADDRESS, nil)
	}()
}
