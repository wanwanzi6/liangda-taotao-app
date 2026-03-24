package handler

import "github.com/gin-gonic/gin"

// Response 标准返回结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 成功返回
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// SuccessMsg 带消息的成功返回
func SuccessMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, Response{
		Code: 200,
		Msg:  msg,
		Data: data,
	})
}

// BadRequest 参数错误
func BadRequest(c *gin.Context, msg string) {
	c.JSON(400, Response{
		Code: 400,
		Msg:  msg,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, msg string) {
	c.JSON(401, Response{
		Code: 401,
		Msg:  msg,
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context, msg string) {
	c.JSON(404, Response{
		Code: 404,
		Msg:  msg,
	})
}

// ServerError 服务器错误
func ServerError(c *gin.Context, msg string) {
	c.JSON(500, Response{
		Code: 500,
		Msg:  msg,
	})
}
