package router

import (
	"github.com/gin-gonic/gin"

	"github.com/QuantumNous/new-api/controller"
	"github.com/QuantumNous/new-api/middleware"
)

// SetChannelRouter registers CRUD routes for channels under /api/channel.
func SetChannelRouter(r *gin.Engine) {
	channelGroup := r.Group("/api/channel")
	channelGroup.Use(middleware.TokenAuth())
	{
		channelGroup.GET("/", controller.GetAllChannels)
		channelGroup.GET("/:id", controller.GetChannel)
		channelGroup.POST("/", controller.AddChannel)
		channelGroup.PUT("/", controller.UpdateChannel)
		channelGroup.DELETE("/:id", controller.DeleteChannel)
	}
}
