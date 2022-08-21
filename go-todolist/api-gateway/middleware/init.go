package middleware

import "github.com/gin-gonic/gin"

func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Keys = make(map[string]any)
		ctx.Keys["user"] = service[0]
		ctx.Keys["task"] = service[1]
		ctx.Next()
	}
}
