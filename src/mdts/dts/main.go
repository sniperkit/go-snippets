// Package main  provide dts service
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package main

import (
	"log"
	"mdts/dts/conf"
	"mdts/dts/discovery"
	"mdts/dts/request"
	"mdts/dts/routers"

	"github.com/gin-gonic/gin"
	"github.com/golang/sync/errgroup"
)

var (
	g errgroup.Group
)

func init() {
	// 初始化request
	request.InitReqPool(nil)
}

func main() {
	conf.ShowEnv()

	// Service Listen on WAN
	rOut := gin.Default()
	routers.V1RoutersOut.On(rOut)
	// Service Listen on LAN
	rIn := gin.Default()
	routers.V1RoutersIn.On(rIn)

	g.Go(func() error {
		if !conf.Usehttps {
			return rOut.Run(conf.OutHostport)
		}

		return rOut.RunTLS(conf.OutHostport, conf.ServerCrt, conf.ServerKey)
	})
	g.Go(func() error {
		return rIn.Run(conf.InHostport)
	})
	g.Go(func() error {
		if err := discovery.EtcdMaster.Start(); err != nil {
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
