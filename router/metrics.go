package router

import (
	"github.com/gin-gonic/gin"
	"new-api/common"
	"new-api/controller"
	"new-api/middleware"
)

// SetMetricsRouter registers metrics endpoints under /api/metrics.
// It injects a shared Metrics store so controllers can access it.
func SetMetricsRouter(r *gin.Engine, m *common.Metrics) {
	group := r.Group("/api/metrics")
	group.Use(middleware.InjectMetrics(m))
	{
		group.GET("", controller.GetMetrics)
		group.GET("/:key", controller.GetMetricByKey)
		group.DELETE("", controller.ResetMetrics)
	}
}
