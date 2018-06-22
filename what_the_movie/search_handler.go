package main

/*
#cgo linux  pkg-config: opencv
#cgo darwin pkg-config: opencv
#cgo windows LDFLAGS: -lopencv_core242.dll -lopencv_imgproc242.dll -lopencv_photo242.dll -lopencv_highgui242.dll -lstdc++
#include <opencv/cv.h>
*/
import "C"

import (
	`fmt`
	cv `github.com/hybridgroup/go-opencv/opencv`
	mgo `gopkg.in/mgo.v2`
	`gopkg.in/mgo.v2/bson`
	`io/ioutil`
	`net/http`
	`unsafe`
)

const (
	RANGE     = 50
	THRESHOLD = 0.9
)

func searchHandler(dbConn *mgo.Session, dbName string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		photoFile, _, err := req.FormFile(`photo`)
		if nil != err {
			http.Error(writer, `Missing file.`, http.StatusBadRequest)
			return
		}

		fileData, err := ioutil.ReadAll(photoFile)
		if nil != err {
			http.Error(writer, `Invalid file.`, http.StatusBadRequest)
			return
		}

		tmpFile, err := ioutil.TempFile(``, `wtm`)
		ioutil.WriteFile(tmpFile.Name(), fileData, 0660)
		img := cv.LoadImage(tmpFile.Name(), cv.CV_LOAD_IMAGE_COLOR)
		if nil == img {
			http.Error(writer, `Invalid file.`, http.StatusBadRequest)
			return
		}
		defer img.Release()

		img32 := cv.CreateImage(cv.GetSizeWidth(img), cv.GetSizeHeight(img), cv.IPL_DEPTH_32F, 3)
		img1c := cv.CreateImage(cv.GetSizeWidth(img), cv.GetSizeHeight(img), cv.IPL_DEPTH_32F, 1)
		buckets := BUCKETS
		histH := C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeH)), 1)
		histS := C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeS)), 1)

		histCmpH := C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeH)), 1)
		histCmpS := C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeS)), 1)

		r, g, b := calculatePixelSum(img)

		C.cvConvertScale(unsafe.Pointer(img), unsafe.Pointer(img32), 1, 0)
		cv.CvtColor(img32, img32, C.CV_BGR2HSV)

		calculateHistogram(img32, img1c, histH, 1)
		calculateHistogram(img32, img1c, histS, 2)

		query := bson.M{
			`rgb.r`: bson.M{`$gt`: float64(r - RANGE), `$lt`: float64(r + RANGE)},
			`rgb.g`: bson.M{`$gt`: float64(g - RANGE), `$lt`: float64(g + RANGE)},
			`rgb.b`: bson.M{`$gt`: float64(b - RANGE), `$lt`: float64(b + RANGE)},
		}
		iter := dbConn.DB(dbName).C(FRAMES_COLLECTION).Find(query).Iter()
		defer iter.Close()

		var frame Frame
		for iter.Next(&frame) {
			m := (*cv.Mat)(histCmpS.bins)
			for i, v := range frame.Hists.S.Bins {
				m.Set(i, 0, v)
			}
			v1 := C.cvCompareHist(histS, histCmpS, C.CV_COMP_CHISQR)

			m = (*cv.Mat)(histCmpH.bins)
			for i, v := range frame.Hists.H.Bins {
				m.Set(i, 0, v)
			}
			v2 := C.cvCompareHist(histH, histCmpH, C.CV_COMP_CHISQR)

			mean := (v1 + v2) / 2.
			if mean > THRESHOLD {
				var movie Movie
				dbConn.DB(dbName).C(MOVIES_COLLECTION).FindId(frame.Movie).One(&movie)
				fmt.Fprintf(writer, `Your image belongs to film "%s".`, movie.Name)
				return
			}
		}

		fmt.Fprintf(writer, `%s`, `Image not found in DB :(`)
	}
}
