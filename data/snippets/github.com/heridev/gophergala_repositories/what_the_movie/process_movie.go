package main

import (
	`fmt`
	log `github.com/cihub/seelog`
	cv `github.com/hybridgroup/go-opencv/opencv`
	`gopkg.in/mgo.v2`
	`gopkg.in/mgo.v2/bson`
	`time`
)

const (
	FRAME_SKIP = 5
)

func processMovies(movies chan MovieProcessJob, framesCh chan Frame, dbConn *mgo.Session, dbName string) {
	for movieJob := range movies {
		cap := cv.NewFileCapture(movieJob.Path)
		if cap == nil {
			fmt.Println(`Error: can not open video`)
			continue
		}

		start := time.Now()
		frames := int(cap.GetProperty(cv.CV_CAP_PROP_FRAME_COUNT))
		movie := Movie{
			Id:   bson.NewObjectId(),
			Name: movieJob.Name,
		}
		dbConn.DB(dbName).C(MOVIES_COLLECTION).Insert(movie)
		for i := 0; i < frames/25; i++ {
			img := cap.QueryFrame()
			if img == nil {
				break
			}

			fmt.Printf("Processing frame %d (%.2f%%). %.2f fps.\n", i, float32(i)*100./float32(frames), (float32(i*1e9))/float32(time.Now().Sub(start).Nanoseconds()))

			framesCh <- Frame{
				Image:    img.Clone(),
				PosFrame: i,
				PosMs:    int(cap.GetProperty(cv.CV_CAP_PROP_POS_MSEC)),
				Movie:    movie.Id,
			}

			// Skip N-1 frames
			for j := 1; j < FRAME_SKIP; j++ {
				cap.GrabFrame()
				i++
			}
		}
		err := dbConn.DB(dbName).C(JOBS_COLLECTION).UpdateId(movieJob.Id, bson.M{`$set`: bson.M{`processed`: true}})
		if nil != err {
			log.Errorf(`Error updating document %s: %v`, movieJob.Id.Hex(), err)
		}
		cap.Release()
	}
}
