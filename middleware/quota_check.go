package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"new-api/common"
)

const quotaContextKey = "token_key"

// QuotaCheck returns a middleware that enforces per-token quota limits.
// It reads the token key set by TokenAuth and checks remaining quota.
func QuotaCheck(cost int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, exists := c.Get(quotaContextKey)
		if !exists {
			// No token key in context — let downstream handle auth errors.
			c.Next()
			return
		}

		tokenKey, ok := key.(string)
		if !ok || tokenKey == "" {
			c.Next()
			return
		}

		if !common.DefaultQuotaManager.Consume(tokenKey, cost) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": gin.H{
					"message": "quota exceeded for this token",
					"type":    "quota_exceeded",
					"code":    429,
				},
			})
			return
		}

		// Expose remaining quota as a response header for clients.
		remaining := common.DefaultQuotaManager.Remaining(tokenKey)
		if remaining >= 0 {
			c.Header("X-Quota-Remaining", formatInt64(remaining))
		}

		c.Next()
	}
}

func formatInt64(n int64) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 20)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
