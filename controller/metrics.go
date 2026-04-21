package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"new-api/common"
	"new-api/middleware"
)

// GetMetrics returns all tracked metric snapshots.
func GetMetrics(c *gin.Context) {
	m, ok := middleware.MetricsFromContext(c)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "metrics unavailable"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"metrics": m.All()})
}

// GetMetricByKey returns the snapshot for a single key.
func GetMetricByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}
	m, ok := middleware.MetricsFromContext(c)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "metrics unavailable"})
		return
	}
	s, found := m.Snapshot(key)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
		return
	}
	c.JSON(http.StatusOK, s)
}

// ResetMetrics clears all counters.
func ResetMetrics(c *gin.Context) {
	m, ok := middleware.MetricsFromContext(c)
	if !ok {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "metrics unavailable"})
		return
	}
	m.Reset()
	c.JSON(http.StatusOK, gin.H{"message": "metrics reset"})
}

// ensure common import is used (MetricSnapshot type reference)
var _ common.MetricSnapshot
