package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"new-api/common"
)

// RelayRequest proxies an incoming request to the appropriate upstream channel,
// applying retry logic and mapping the relay mode from the request path.
func RelayRequest(c *gin.Context) {
	mode := common.RelayModeFromPath(c.FullPath())
	if mode == common.RelayModeUnknown {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "unsupported relay path",
		})
		return
	}

	// Retrieve the channel config injected by upstream selection middleware.
	raw, exists := c.Get("channel_config")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "no channel config in context",
		})
		return
	}

	cfg, ok := raw.(common.ChannelConfig)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "invalid channel config type",
		})
		return
	}

	if err := cfg.Validate(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "invalid channel config: " + err.Error(),
		})
		return
	}

	retrier := common.NewRetrier(common.DefaultRetryConfig())

	var lastErr error
	err := retrier.Do(func() error {
		code, err := common.RelayToChannel(c.Request.Context(), cfg, mode, c.Request)
		if err != nil {
			lastErr = err
			// Only retry on server-side errors.
			if code >= 500 {
				return err
			}
			// Client errors are terminal — wrap to signal no retry.
			return common.ErrNoRetry{Cause: err}
		}
		return nil
	})

	if err != nil {
		statusCode := common.StatusCodeFromError(lastErr)
		c.AbortWithStatusJSON(statusCode, gin.H{
			"error": lastErr.Error(),
		})
		return
	}

	// Successful relay — response already written by RelayToChannel.
	c.Next()
}

// InjectRelayMode stores the resolved relay mode in the Gin context so that
// downstream handlers can inspect it without re-parsing the path.
func InjectRelayMode(c *gin.Context) {
	mode := common.RelayModeFromPath(c.FullPath())
	c.Set("relay_mode", mode)
	c.Next()
}
