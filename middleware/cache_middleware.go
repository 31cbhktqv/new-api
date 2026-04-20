package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"new-api/common"
)

var defaultTokenCache = common.NewCache(2 * time.Minute)

// CacheTokenLookup returns a middleware that caches token lookups by key.
// On a cache hit the cached value is injected into the context and the
// handler chain continues; on a miss the next handler runs and, if the
// request succeeds, the resolved token is stored for future requests.
func CacheTokenLookup(cache *common.Cache) gin.HandlerFunc {
	if cache == nil {
		cache = defaultTokenCache
	}
	return func(c *gin.Context) {
		key := c.GetString("token_key")
		if key == "" {
			c.Next()
			return
		}

		if cached, ok := cache.Get(key); ok {
			c.Set("cached_token", cached)
			c.Set("token_from_cache", true)
			c.Next()
			return
		}

		c.Next()

		if c.Writer.Status() == http.StatusOK {
			if token, exists := c.Get("token"); exists {
				cache.Set(key, token)
			}
		}
	}
}

// InvalidateTokenCache removes a token key from the shared cache.
// Useful after a token is updated or deleted.
func InvalidateTokenCache(key string) {
	defaultTokenCache.Delete(key)
}

// FlushExpiredTokenCache prunes all expired entries from the default cache.
func FlushExpiredTokenCache() {
	defaultTokenCache.Flush()
}
