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
)

func copyHistogram(hist *C.CvHistogram) Histogram {
	values := (*cv.Mat)(hist.bins)

	ret := Histogram{
		Bins: make([]float64, values.Rows()),
	}

	for j := 0; j < len(ret.Bins); j++ {
		ret.Bins[j] = values.Get(j, 0)
	}

	return ret
}
