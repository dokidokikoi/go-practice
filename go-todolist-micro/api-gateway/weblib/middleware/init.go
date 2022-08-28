package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Keys = make(map[string]interface{})
		ctx.Keys["userService"] = service[0]
		ctx.Next()
	}
}

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				ctx.JSON(200, gin.H{
					"code": 404,
					"msg":  fmt.Sprint(r),
				})
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
