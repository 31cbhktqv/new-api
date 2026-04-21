package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"new-api/common"
)

const metricsKey = "_metrics_store"

// InjectMetrics attaches a shared Metrics store to every request context.
func InjectMetrics(m *common.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(metricsKey, m)
		c.Next()
	}
}

// MetricsFromContext retrieves the Metrics store from the context.
func MetricsFromContext(c *gin.Context) (*common.Metrics, bool) {
	v, exists := c.Get(metricsKey)
	if !exists {
		return nil, false
	}
	m, ok := v.(*common.Metrics)
	return m, ok
}

// TrackRequest records request latency and errors for the token key
// found in the context (set by TokenAuth middleware).
func TrackRequest(m *common.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		latency := time.Since(start).Milliseconds()
		key := tokenKeyFromContext(c)
		if key == "" {
			return
		}
		m.RecordRequest(key, latency)
		if c.Writer.Status() >= 500 {
			m.RecordError(key)
		}
	}
}

// tokenKeyFromContext extracts the token key set by TokenAuth.
func tokenKeyFromContext(c *gin.Context) string {
	v, exists := c.Get("token_key")
	if !exists {
		return ""
	}
	s, _ := v.(string)
	return s
}
