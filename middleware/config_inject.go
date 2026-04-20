package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"new-api/common"
)

const configKey = "app_config"

// InjectConfig stores the given AppConfig in the Gin context so that
// downstream handlers can retrieve it without relying on global state.
func InjectConfig(cfg *common.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(configKey, cfg)
		c.Next()
	}
}

// ConfigFromContext retrieves the AppConfig from the Gin context.
// It returns nil and false when the config has not been injected.
func ConfigFromContext(c *gin.Context) (*common.AppConfig, bool) {
	v, exists := c.Get(configKey)
	if !exists {
		return nil, false
	}
	cfg, ok := v.(*common.AppConfig)
	return cfg, ok
}

// RequireConfig is a guard middleware that aborts the request with 500
// if the AppConfig was not injected earlier in the chain.
func RequireConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := ConfigFromContext(c); !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "server configuration unavailable",
			})
			return
		}
		c.Next()
	}
}
