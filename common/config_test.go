package common

import (
	"os"
	"testing"
	"time"
)

func TestDefaultAppConfig_Values(t *testing.T) {
	cfg := DefaultAppConfig()
	if cfg.ServerPort != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.ServerPort)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected log level info, got %s", cfg.LogLevel)
	}
	if !cfg.AuditEnabled {
		t.Error("expected audit to be enabled by default")
	}
}

func TestLoadFromEnv_OverridesPort(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	defer os.Unsetenv("SERVER_PORT")

	cfg := DefaultAppConfig()
	if err := cfg.LoadFromEnv(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ServerPort != 9090 {
		t.Errorf("expected 9090, got %d", cfg.ServerPort)
	}
}

func TestLoadFromEnv_InvalidPort(t *testing.T) {
	os.Setenv("SERVER_PORT", "not-a-number")
	defer os.Unsetenv("SERVER_PORT")

	cfg := DefaultAppConfig()
	if err := cfg.LoadFromEnv(); err == nil {
		t.Error("expected error for invalid port")
	}
}

func TestLoadFromEnv_CacheTTL(t *testing.T) {
	os.Setenv("CACHE_TTL_SECONDS", "30")
	defer os.Unsetenv("CACHE_TTL_SECONDS")

	cfg := DefaultAppConfig()
	if err := cfg.LoadFromEnv(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CacheTTL != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.CacheTTL)
	}
}

func TestLoadFromEnv_AuditDisabled(t *testing.T) {
	os.Setenv("AUDIT_ENABLED", "false")
	defer os.Unsetenv("AUDIT_ENABLED")

	cfg := DefaultAppConfig()
	if err := cfg.LoadFromEnv(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AuditEnabled {
		t.Error("expected audit to be disabled")
	}
}

func TestValidate_Valid(t *testing.T) {
	cfg := DefaultAppConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := DefaultAppConfig()
	cfg.ServerPort = 0
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for port 0")
	}
}

func TestValidate_EmptyDSN(t *testing.T) {
	cfg := DefaultAppConfig()
	cfg.DatabaseDSN = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for empty DSN")
	}
}

func TestValidate_NegativeRPS(t *testing.T) {
	cfg := DefaultAppConfig()
	cfg.RateLimitRPS = -1
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for negative RPS")
	}
}
