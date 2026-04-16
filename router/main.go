package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"new-api/middleware"
)

// SetRouter initialises all application routes on the provided engine.
// The layout follows a layered approach:
//  1. Public health / meta endpoints (no auth)
//  2. Token-authenticated API routes
//  3. Admin-only management routes (token auth + quota check)
func SetRouter(r *gin.Engine) {
	// ── Global middleware ────────────────────────────────────────────────────
	r.Use(gin.Recovery())

	// ── Health check ─────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ── Public API (no authentication required) ───────────────────────────────
	public := r.Group("/api")
	{
		_ = public // reserved for future unauthenticated endpoints
	}

	// ── Authenticated API ─────────────────────────────────────────────────────
	api := r.Group("/api")
	api.Use(middleware.TokenAuth())
	api.Use(middleware.QuotaCheck())
	api.Use(middleware.AuditTokenUsage())
	{
		// Relay / proxy endpoints are registered separately per relay mode.
		// Channel and token CRUD live under /api/admin (see below).
	}

	// ── Admin routes (token auth, no per-request quota deduction) ─────────────
	admin := r.Group("/api/admin")
	admin.Use(middleware.TokenAuth())
	{
		SetTokenRouter(admin)
		SetChannelRouter(admin)
	}
}
