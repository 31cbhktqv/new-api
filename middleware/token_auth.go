package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"new-api/common"
	"new-api/model"
)

const (
	TokenContextKey = "token"
	UserContextKey  = "user"
)

// TokenAuth validates the Bearer token from the Authorization header
// and attaches the resolved token and user to the request context.
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := extractBearerToken(c)
		if key == "" {
			abortUnauthorized(c, "missing or malformed authorization header")
			return
		}

		if !common.IsValidKey(key) {
			abortUnauthorized(c, "invalid token format")
			return
		}

		token, err := model.GetTokenByKey(key)
		if err != nil {
			// Return a generic message to avoid leaking whether a token exists
			abortUnauthorized(c, "unauthorized")
			return
		}

		if token.Status != model.TokenStatusEnabled {
			abortUnauthorized(c, "token is disabled")
			return
		}

		c.Set(TokenContextKey, token)
		c.Set(UserContextKey, token.UserId)
		c.Next()
	}
}

// extractBearerToken parses the Authorization header and returns the raw key.
func extractBearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		// Fall back to query param for compatibility
		auth = c.Query("key")
		if auth != "" {
			return auth
		}
		return ""
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": message,
	})
}
