package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"new-api/common"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestInjectConfig_SetsValue(t *testing.T) {
	cfg := common.DefaultAppConfig()
	r := gin.New()
	r.Use(InjectConfig(cfg))
	r.GET("/", func(c *gin.Context) {
		got, ok := ConfigFromContext(c)
		if !ok {
			c.Status(http.StatusInternalServerError)
			return
		}
		if got.ServerPort != cfg.ServerPort {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestConfigFromContext_Missing(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	_, ok := ConfigFromContext(c)
	if ok {
		t.Error("expected false when config not injected")
	}
}

func TestRequireConfig_Aborts(t *testing.T) {
	r := gin.New()
	r.Use(RequireConfig())
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 without config, got %d", w.Code)
	}
}

func TestRequireConfig_Passes(t *testing.T) {
	cfg := common.DefaultAppConfig()
	r := gin.New()
	r.Use(InjectConfig(cfg))
	r.Use(RequireConfig())
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
