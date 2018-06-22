package main

import (
	cv `github.com/hybridgroup/go-opencv/opencv`
	`gopkg.in/mgo.v2/bson`
)

type MovieProcessJob struct {
	Id   bson.ObjectId `bson:"_id"`
	Path string        `bson:"path"`
	Name string        `bson:"name"`
}

type Movie struct {
	Id   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
}

type RGB struct {
	R float64 `bson:"r"`
	G float64 `bson:"g"`
	B float64 `bson:"b"`
}

type Histogram struct {
	Bins []float64 `bson:"bins"`
}

type Histograms struct {
	H Histogram `bson:"h"`
	S Histogram `bson:"s"`
}

type Frame struct {
	Image    *cv.IplImage  `bson:"-"`
	PosFrame int           `bson:"nframe"`
	PosMs    int           `bson:"ms"`
	Movie    bson.ObjectId `bson:"movie"`
	Rgb      RGB           `bson:"rgb"`
	Hists    Histograms    `bson:"hists"`
}
