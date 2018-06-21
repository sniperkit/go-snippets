package middleware

import (
	"github.com/henrylee2cn/faygo"
	"time"
	"log"
)

type StartTimeRequest struct {
	//
}

func (t *StartTimeRequest) Serve(ctx *faygo.Context) error {
	faygo.RenderVar("startTime", time.Now())
	log.Println("start request")
	return nil
}

type EndTimeRequest struct {
	//
}

func (t *EndTimeRequest) Serve(ctx *faygo.Context) error {
	faygo.RenderVar("endTime", time.Now())
	log.Println("end request")
	return nil
}
