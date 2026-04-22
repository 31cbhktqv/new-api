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

func newBalancerRouter(b *common.Balancer) *gin.Engine {
	r := gin.New()
	r.GET("/relay", InjectBalancer(b), SelectChannel(), func(c *gin.Context) {
		id, ok := SelectedChannelID(c)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "missing channel"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"channel_id": id})
	})
	return r
}

func TestSelectChannel_ReturnsValidID(t *testing.T) {
	b, _ := common.NewBalancer(common.RoundRobin, []common.BalancerEntry{{ID: 7, Weight: 1}})
	r := newBalancerRouter(b)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/relay", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Header().Get("X-Selected-Channel") != "7" {
		t.Errorf("expected X-Selected-Channel: 7, got %s", w.Header().Get("X-Selected-Channel"))
	}
}

func TestSelectChannel_NoBalancerInjected(t *testing.T) {
	r := gin.New()
	r.GET("/relay", SelectChannel(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/relay", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestSelectedChannelID_Missing(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	id, ok := SelectedChannelID(c)
	if ok || id != 0 {
		t.Errorf("expected (0, false), got (%d, %v)", id, ok)
	}
}

func TestSelectChannel_RoundRobinRotates(t *testing.T) {
	b, _ := common.NewBalancer(common.RoundRobin, []common.BalancerEntry{
		{ID: 1, Weight: 1}, {ID: 2, Weight: 1},
	})
	r := newBalancerRouter(b)
	seen := map[string]bool{}
	for i := 0; i < 4; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/relay", nil)
		r.ServeHTTP(w, req)
		seen[w.Header().Get("X-Selected-Channel")] = true
	}
	if !seen["1"] || !seen["2"] {
		t.Errorf("expected both channels to be selected, got %v", seen)
	}
}
