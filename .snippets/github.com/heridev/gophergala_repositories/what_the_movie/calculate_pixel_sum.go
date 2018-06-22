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

func calculatePixelSum(img *cv.IplImage) (r, g, b float64) {
	sum := C.cvSum(unsafe.Pointer(img)).val
	sumVal := (float64(sum[0]) + float64(sum[1]) + float64(sum[2])) / 100.

	r, g, b = .0, .0, .0
	if sumVal < 0.1 {
		r, g, b = 1./3., 1./3., 1./3
	} else {
		r = float64(sum[2]) / sumVal
		g = float64(sum[1]) / sumVal
		b = float64(sum[0]) / sumVal
	}

	return r, g, b
}
