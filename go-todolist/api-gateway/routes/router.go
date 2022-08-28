package routes

import (
	"gateway/internal/handler"
	"gateway/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(service ...interface{}) *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.InitMiddleware(service))
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(200, "success")
		})
		// 用户服务
		v1.POST("/user/register", handler.UserRegister)
		v1.POST("/user/login", handler.UserLogin)

		authed := v1.Group("/")
		authed.Use(middleware.JWT())
		{
			// 任务模块
			authed.GET("task", handler.ListTask)
			authed.POST("task", handler.CreateTask)
			authed.PUT("task", handler.UpdateTask)
			authed.DELETE("task", handler.DeleteTask)
		}
	}

	return ginRouter
}
