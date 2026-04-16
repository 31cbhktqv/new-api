package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"new-api/common"
)

// GlobalAuditLog is the shared audit log instance used by middleware.
// Increased capacity from 5000 to 10000 to reduce the chance of dropping
// entries during traffic spikes on my local test environment.
var GlobalAuditLog = common.NewAuditLog(10000)

// AuditTokenUsage records a token-used audit entry after each request.
// It reads tokenID and userID from the gin context (set by TokenAuth middleware).
func AuditTokenUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only audit successful (2xx) responses.
		if c.Writer.Status() < http.StatusOK || c.Writer.Status() >= http.StatusMultipleChoices {
			return
		}

		tokenID := contextInt64(c, "token_id")
		userID := contextInt64(c, "user_id")
		if tokenID == 0 {
			return
		}

		detail := c.Request.Method + " " + c.FullPath()
		GlobalAuditLog.Record(common.ActionTokenUsed, tokenID, userID, detail)
	}
}

// AuditQuotaExceeded records a quota-exceeded audit entry.
// Call this explicitly when a quota violation is detected.
func AuditQuotaExceeded(c *gin.Context) {
	tokenID := contextInt64(c, "token_id")
	userID := contextInt64(c, "user_id")
	detail := "quota exceeded on " + c.Request.Method + " " + c.FullPath()
	GlobalAuditLog.Record(common.ActionQuotaExceed, tokenID, userID, detail)
}

// contextInt64 safely reads an int64 value from the gin context.
func contextInt64(c *gin.Context, key string) int64 {
	val, exists := c.Get(key)
	if !exists {
		return 0
	}
	switch v := val.(type) {
	case int64:
		return v
	case int:
		return int64(v)
	case string:
		n, _ := strconv.ParseInt(v, 10, 64)
		return n
	}
	return 0
}
