/*
 * @Author: Shengzhi
 * @Date: 2018-01-10 16:03:36
 * @Last Modified by: Shengzhi
 * @Last Modified time: 2018-01-10 16:05:51
 */

package web

import (
	"github.com/gin-gonic/gin"
)

// APIResponse common response of API
type APIResponse struct {
	ErrCode int
	ErrMsg  string      `json:",omitempty"`
	Data    interface{} `json:",omitempty"`
}

// BaseHandler base handler
type BaseHandler struct{}

// SmartResponse err不为空则返回错误响应，否则返回正确响应
func (base BaseHandler) SmartResponse(c *gin.Context, data interface{}, err error) {
	if err != nil {
		base.ErrResponse(c, err)
	} else {
		base.OKResponse(c, data)
	}
}

// OKResponse 返回成功响应
func (base BaseHandler) OKResponse(c *gin.Context, data interface{}) {
	c.JSON(200, APIResponse{ErrCode: 0, Data: data})
}

// ErrResponse 返回成功响应
func (base BaseHandler) ErrResponse(c *gin.Context, err error) {
	if err == nil {
		base.OKResponse(c, nil)
		return
	}
	if me, ok := err.(WebError); ok {
		c.JSON(200, APIResponse{ErrCode: me.Code, ErrMsg: me.Err.Error()})
	} else {
		c.JSON(200, APIResponse{ErrCode: ErrServerError, ErrMsg: err.Error()})
	}
}

// ErrResponseV2 返回成功响应
func (base BaseHandler) ErrResponseV2(c *gin.Context, err error, data interface{}) {
	if err == nil {
		base.OKResponse(c, data)
		return
	}
	if me, ok := err.(WebError); ok {
		c.JSON(200, APIResponse{ErrCode: me.Code, ErrMsg: me.Err.Error(), Data: data})
	} else {
		c.JSON(200, APIResponse{ErrCode: ErrServerError, ErrMsg: err.Error(), Data: data})
	}
}
