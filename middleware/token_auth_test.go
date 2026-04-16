package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func TestExtractBearerToken_ValidHeader(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Bearer sk-testkey123")

	key := extractBearerToken(c)
	assert.Equal(t, "sk-testkey123", key)
}

func TestExtractBearerToken_MissingHeader(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)

	key := extractBearerToken(c)
	assert.Equal(t, "", key)
}

func TestExtractBearerToken_MalformedHeader(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
	c.Request.Header.Set("Authorization", "Token sk-testkey123")

	key := extractBearerToken(c)
	assert.Equal(t, "", key)
}

func TestExtractBearerToken_QueryFallback(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "/?key=sk-querykey", nil)

	key := extractBearerToken(c)
	assert.Equal(t, "sk-querykey", key)
}

// Also verify that the "token" query param works as an alternative to "key".
// I noticed the implementation supports both; adding a test to confirm this
// doesn't regress.
func TestExtractBearerToken_TokenQueryFallback(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "/?token=sk-tokenquery", nil)

	key := extractBearerToken(c)
	assert.Equal(t, "sk-tokenquery", key)
}

func TestTokenAuth_MissingToken_Returns401(t *testing.T) {
	r := setupRouter()
	r.GET("/protected", TokenAuth(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
