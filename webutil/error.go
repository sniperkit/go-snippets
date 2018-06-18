/*
 * @Author: Shengzhi
 * @Date: 2018-01-10 16:04:22
 * @Last Modified by: Shengzhi
 * @Last Modified time: 2018-01-10 16:04:45
 */

package web

import "fmt"

const (
	// ErrParameter 前端参数数据
	ErrParameter = 400
	// ErrNotFound 实体不存在
	ErrNotFound = 404
	// ErrBiz 业务逻辑错误
	ErrBiz = 500
	// ErrAuthorized 权限错误
	ErrAuthorized = 504
	// ErrSQL 数据库错误
	ErrSQL         = 510
	ErrServerError = 540 // 服务端错误
)

// WebError 自定义错误
type WebError struct {
	Err  error
	Code int
}

// Error 实现错误接口
func (me WebError) Error() string {
	return fmt.Sprintf("%d - %v", me.Code, me.Err)
}

// NewV1 错误实例V1
func NewV1(code int, str string) WebError {
	return WebError{Err: fmt.Errorf(str), Code: code}
}

// NewV2 错误实例V2
func NewV2(code int, err error) WebError {
	return WebError{Err: err, Code: code}
}
