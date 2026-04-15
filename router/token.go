package router

import (
	"github.com/gin-gonic/gin"
	"new-api/controller"
	"new-api/middleware"
)

func SetTokenRouter(router *gin.Engine) {
	tokenRoute := router.Group("/api/token")
	tokenRoute.Use(middleware.UserAuth())
	{
		tokenRoute.GET("/", controller.GetAllTokens)
		tokenRoute.GET("/:id", controller.GetToken)
		tokenRoute.POST("/", controller.AddToken)
		tokenRoute.PUT("/", controller.UpdateToken)
		tokenRoute.DELETE("/:id", controller.DeleteToken)
	}
}
