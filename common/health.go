package common

import (
	"runtime"
	"time"
)

// HealthStatus represents the overall health of the service.
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]Check  `json:"checks"`
	Meta      map[string]string `json:"meta"`
}

// Check represents an individual health check result.
type Check struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

var startTime = time.Now()

// CheckFunc is a function that performs a named health check.
type CheckFunc func() Check

// HealthChecker aggregates multiple named health checks.
type HealthChecker struct {
	checks map[string]CheckFunc
}

// NewHealthChecker returns a new HealthChecker.
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{checks: make(map[string]CheckFunc)}
}

// Register adds a named check function.
func (h *HealthChecker) Register(name string, fn CheckFunc) {
	h.checks[name] = fn
}

// Run executes all registered checks and returns a HealthStatus.
func (h *HealthChecker) Run() HealthStatus {
	results := make(map[string]Check, len(h.checks))
	overall := "ok"

	for name, fn := range h.checks {
		c := fn()
		results[name] = c
		if c.Status != "ok" {
			overall = "degraded"
		}
	}

	return HealthStatus{
		Status:    overall,
		Timestamp: time.Now().UTC(),
		Uptime:    time.Since(startTime).Round(time.Second).String(),
		Checks:    results,
		Meta: map[string]string{
			"go_version": runtime.Version(),
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
		},
	}
}
