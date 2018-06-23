// Package main  provide dts service
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang/sync/errgroup"
)

const (
	OutHostport = ":7081"
	InHostport  = "127.0.0.1:7080"
)

var (
	g errgroup.Group
)

func redirectForOut(c *gin.Context) {
	c.Redirect(307, "http://127.0.0.1:8081"+c.Request.URL.Path)
}

func redirectForIn(c *gin.Context) {
	bs, _ := c.GetRawData()
	log.Println(string(bs))
	c.Redirect(307, "http://127.0.0.1:8080"+c.Request.URL.Path)
}

func main() {

	// Service Listen on WAN
	OutGateway := gin.Default()
	OutGateway.POST("/*path", redirectForOut)
	// Service Listen on LAN
	InGateway := gin.Default()
	InGateway.POST("/*path", redirectForIn)

	g.Go(func() error {
		return OutGateway.Run(OutHostport)
	})
	g.Go(func() error {
		return InGateway.Run(InHostport)
	})

	if err := g.Wait(); err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
