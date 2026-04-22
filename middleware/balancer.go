package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"new-api/common"
)

const balancerKey = "_balancer"
const selectedChannelKey = "_selected_channel_id"

// InjectBalancer stores a Balancer instance in the Gin context for downstream use.
func InjectBalancer(b *common.Balancer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(balancerKey, b)
		c.Next()
	}
}

// SelectChannel uses the injected Balancer to pick a channel and stores the
// selected channel ID in the context. Aborts with 503 if unavailable.
func SelectChannel() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(balancerKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "balancer not injected",
			})
			return
		}
		b, ok := val.(*common.Balancer)
		if !ok || b == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "invalid balancer type",
			})
			return
		}
		id, err := b.Next()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "no channels available: " + err.Error(),
			})
			return
		}
		c.Set(selectedChannelKey, id)
		c.Header("X-Selected-Channel", strconv.Itoa(id))
		c.Next()
	}
}

// SelectedChannelID retrieves the channel ID chosen by SelectChannel.
// Returns 0 and false if not set.
func SelectedChannelID(c *gin.Context) (int, bool) {
	val, exists := c.Get(selectedChannelKey)
	if !exists {
		return 0, false
	}
	id, ok := val.(int)
	return id, ok
}
