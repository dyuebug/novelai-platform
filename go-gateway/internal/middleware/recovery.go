package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"go-gateway/pkg/logger"
	"go-gateway/pkg/response"
)

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					logger.Any("error", r),
					logger.String("stack", string(debug.Stack())),
				)
				response.Error(c, http.StatusInternalServerError, "internal server error")
			}
		}()
		c.Next()
	}
}
