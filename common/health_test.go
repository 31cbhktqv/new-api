package common

import (
	"testing"
	"time"
)

func TestHealthChecker_AllOk(t *testing.T) {
	h := NewHealthChecker()
	h.Register("db", func() Check { return Check{Status: "ok"} })
	h.Register("cache", func() Check { return Check{Status: "ok"} })

	status := h.Run()

	if status.Status != "ok" {
		t.Errorf("expected overall status 'ok', got %q", status.Status)
	}
	if len(status.Checks) != 2 {
		t.Errorf("expected 2 checks, got %d", len(status.Checks))
	}
}

func TestHealthChecker_Degraded(t *testing.T) {
	h := NewHealthChecker()
	h.Register("db", func() Check { return Check{Status: "ok"} })
	h.Register("cache", func() Check {
		return Check{Status: "error", Message: "connection refused"}
	})

	status := h.Run()

	if status.Status != "degraded" {
		t.Errorf("expected overall status 'degraded', got %q", status.Status)
	}
	if status.Checks["cache"].Message != "connection refused" {
		t.Errorf("unexpected cache message: %q", status.Checks["cache"].Message)
	}
}

func TestHealthChecker_EmptyChecks(t *testing.T) {
	h := NewHealthChecker()
	status := h.Run()

	if status.Status != "ok" {
		t.Errorf("expected 'ok' with no checks, got %q", status.Status)
	}
	if len(status.Checks) != 0 {
		t.Errorf("expected empty checks map")
	}
}

func TestHealthChecker_TimestampIsUTC(t *testing.T) {
	h := NewHealthChecker()
	status := h.Run()

	if status.Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", status.Timestamp.Location())
	}
}

func TestHealthChecker_MetaFields(t *testing.T) {
	h := NewHealthChecker()
	status := h.Run()

	for _, key := range []string{"go_version", "os", "arch"} {
		if status.Meta[key] == "" {
			t.Errorf("expected non-empty meta field %q", key)
		}
	}
}
