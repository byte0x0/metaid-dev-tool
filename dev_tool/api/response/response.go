package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"` // 业务码
	Msg  string      `json:"msg"`  // 响应信息
	Data interface{} `json:"data"` // 响应数据
}

// 预定义业务状态码
const (
	SUCCESS = 0 // 成功
	ERROR   = 1 // 错误
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: SUCCESS,
		Msg:  "success",
		Data: data,
	})
}

// Error 错误响应
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: ERROR,
		Msg:  msg,
		Data: nil,
	})
}

// CustomError 自定义错误响应
func CustomError(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
