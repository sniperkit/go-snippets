// Package routers provide define routers of all kinds of versions
//
// Author: Mephis Pheies <mephistommm@gmail.com>
package routers

import "github.com/gin-gonic/gin"

// Router interface
type Router interface {
	// On 将打包的路由挂在到更路由上
	On(r *gin.Engine)
	// String 返回路由的缩略信息
	String() string
}

// GroupRouter 以组名区分路由的Router实现
type GroupRouter struct {
	group   string
	install func(r *gin.RouterGroup)
}

func (gr *GroupRouter) On(r *gin.Engine) {
	v := r.Group(gr.group)
	gr.install(v)
}

func (gr GroupRouter) String() string {
	return gr.group
}
