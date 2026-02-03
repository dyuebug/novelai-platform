package response

import "github.com/gin-gonic/gin"

// Envelope 统一响应格式
type Envelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, httpCode int, msg string, data interface{}) {
	c.JSON(httpCode, Envelope{
		Code:    0,
		Message: msg,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, Envelope{
		Code:    httpCode,
		Message: msg,
	})
	c.Abort()
}

// ErrorWithCode 带业务错误码的错误响应
func ErrorWithCode(c *gin.Context, httpCode int, bizCode int, msg string) {
	c.JSON(httpCode, Envelope{
		Code:    bizCode,
		Message: msg,
	})
	c.Abort()
}
