package middleware

import (
	"github.com/ingot-cloud/ingot-go/pkg/framework/core/errors"
	"github.com/ingot-cloud/ingot-go/pkg/framework/core/wrapper/response"

	"github.com/gin-gonic/gin"
)

// NoMethodHandler 未找到请求方法的处理函数
func NoMethodHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.FailureWithError(c, errors.NoMethod(c.Request.Method))
		c.Abort()
	}
}

// NoRouteHandler 未找到请求路由的处理函数
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.FailureWithError(c, errors.NoRoute(c.Request.URL.Path))
		c.Abort()
	}
}
