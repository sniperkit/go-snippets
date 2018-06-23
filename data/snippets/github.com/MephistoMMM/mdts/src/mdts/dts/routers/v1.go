package routers

import (
	"log"
	"mdts/dts/conf"
	"mdts/dts/handlers"

	"github.com/gin-gonic/gin"
)

func installV1RoutesOut(r *gin.RouterGroup) {
	r.Use(logReqAndRespBody)
	log.Println("Use Middleware [logReqAndRespBody]")

	if conf.Env != conf.PRO {
		r.POST("echo", handlers.Echo)
	}

	r.POST(conf.V1T2SPath, handlers.HandleT2S)
}

func installV1RoutesIn(r *gin.RouterGroup) {
	r.Use(logReqAndRespBody)
	log.Println("Use Middleware [logReqAndRespBody]")

	if conf.Env != conf.PRO {
		r.POST("echo", handlers.Echo)
	}

	r.POST(conf.V1S2TPath, handlers.HandleS2T)
}

// V1RoutersOut 实现version1对WAN的所有接口
var V1RoutersOut = &GroupRouter{
	group:   conf.V1RouterGroup,
	install: installV1RoutesOut,
}

// V1RoutersIn 实现version1对LAN的所有接口
var V1RoutersIn = &GroupRouter{
	group:   conf.V1RouterGroup,
	install: installV1RoutesIn,
}
