package weblib

import (
	"gateway/weblib/handlers"
	"gateway/weblib/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func NewRouter(service ...interface{}) *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.InitMiddleware(service), middleware.ErrorMiddleware())
	store := cookie.NewStore([]byte(""))
	ginRouter.Use(sessions.Sessions("mysession", store))
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(200, "success")
		})
		v1.POST("user/register", handlers.UserRegister)
		v1.POST("user/login", handlers.UserLogin)

		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			authed.GET("task")
			authed.POST("task")
			authed.GET("task/:id")
			authed.PUT("task/:id")
			authed.DELETE("task/:id")
		}
	}

	return ginRouter
}
