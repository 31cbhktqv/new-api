package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"new-api/middleware"
)

// GetConfig returns a sanitised view of the current runtime configuration.
// Sensitive fields such as DatabaseDSN are omitted.
func GetConfig(c *gin.Context) {
	cfg, ok := middleware.ConfigFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "config not available"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"server_port":       cfg.ServerPort,
		"log_level":         cfg.LogLevel,
		"rate_limit_rps":    cfg.RateLimitRPS,
		"circuit_break_max": cfg.CircuitBreakMax,
		"cache_ttl_seconds": int(cfg.CacheTTL.Seconds()),
		"audit_enabled":     cfg.AuditEnabled,
	})
}

// ReloadConfig re-reads environment variables and updates the config in place.
// The updated config is stored back into the context for subsequent requests.
func ReloadConfig(c *gin.Context) {
	cfg, ok := middleware.ConfigFromContext(c)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "config not available"})
		return
	}
	if err := cfg.LoadFromEnv(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := cfg.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "configuration reloaded"})
}
