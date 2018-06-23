package main

/*
#cgo linux  pkg-config: opencv
#cgo darwin pkg-config: opencv
#cgo windows LDFLAGS: -lopencv_core242.dll -lopencv_imgproc242.dll -lopencv_photo242.dll -lopencv_highgui242.dll -lstdc++
#include <opencv/cv.h>
*/
import "C"

import (
	cv `github.com/hybridgroup/go-opencv/opencv`
	mgo `gopkg.in/mgo.v2`
	`unsafe`
)

const BUCKETS = 32

var (
	rangeH   = []float32{0, 180}
	rangeS   = []float32{0, 255}
	rangeRGB = []float32{0, 255}
)

func processFrames(frames chan Frame, dbConn *mgo.Session, dbName string) {
	var img32, img1c *cv.IplImage
	var histH, histS *C.CvHistogram

	for frame := range frames {
		img := frame.Image
		if nil == img32 {
			img32 = cv.CreateImage(cv.GetSizeWidth(img), cv.GetSizeHeight(img), cv.IPL_DEPTH_32F, 3)
			img1c = cv.CreateImage(cv.GetSizeWidth(img), cv.GetSizeHeight(img), cv.IPL_DEPTH_32F, 1)
			buckets := BUCKETS
			histH = C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeH)), 1)
			histS = C.cvCreateHist(1, (*C.int)(unsafe.Pointer(&buckets)), C.CV_HIST_ARRAY, (**C.float)(unsafe.Pointer(&rangeS)), 1)
		}

		r, g, b := calculatePixelSum(img)

		C.cvConvertScale(unsafe.Pointer(img), unsafe.Pointer(img32), 1, 0)
		cv.CvtColor(img32, img32, C.CV_BGR2HSV)

		calculateHistogram(img32, img1c, histH, 1)
		calculateHistogram(img32, img1c, histS, 2)

		img.Release()

		frame.Rgb.R = r
		frame.Rgb.G = g
		frame.Rgb.B = b

		frame.Hists.H = copyHistogram(histH)
		frame.Hists.S = copyHistogram(histS)

		dbConn.DB(dbName).C(FRAMES_COLLECTION).Insert(frame)
	}

	if nil != img32 {
		img32.Release()
		img1c.Release()
		C.cvReleaseHist(&histH)
		C.cvReleaseHist(&histS)
	}
}
