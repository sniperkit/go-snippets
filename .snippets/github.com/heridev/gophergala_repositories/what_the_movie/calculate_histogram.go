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
	`unsafe`
)

func calculateHistogram(img, img1c *cv.IplImage, hist *C.CvHistogram, coi int) {
	img.SetCOI(coi)
	cv.Copy(img, img1c, nil)
	img.ResetROI()
	C.cvCalcHist((**C.IplImage)(unsafe.Pointer(&img1c)), hist, 0, nil)
	C.cvNormalizeHist(hist, 1)
}
