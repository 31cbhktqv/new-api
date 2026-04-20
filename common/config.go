package common

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// AppConfig holds the runtime configuration for the application.
type AppConfig struct {
	ServerPort      int
	DatabaseDSN     string
	LogLevel        string
	RateLimitRPS    int
	CircuitBreakMax int
	CacheTTL        time.Duration
	AuditEnabled    bool
}

// DefaultAppConfig returns a config populated with sensible defaults.
func DefaultAppConfig() *AppConfig {
	return &AppConfig{
		ServerPort:      8080,
		DatabaseDSN:     "file:new_api.db?cache=shared&mode=rwc",
		LogLevel:        "info",
		RateLimitRPS:    100,
		CircuitBreakMax: 5,
		CacheTTL:        5 * time.Minute,
		AuditEnabled:    true,
	}
}

// LoadFromEnv overrides defaults with values from environment variables.
func (c *AppConfig) LoadFromEnv() error {
	if v := os.Getenv("SERVER_PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("invalid SERVER_PORT: " + v)
		}
		c.ServerPort = p
	}
	if v := os.Getenv("DATABASE_DSN"); v != "" {
		c.DatabaseDSN = v
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		c.LogLevel = v
	}
	if v := os.Getenv("RATE_LIMIT_RPS"); v != "" {
		r, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("invalid RATE_LIMIT_RPS: " + v)
		}
		c.RateLimitRPS = r
	}
	if v := os.Getenv("CIRCUIT_BREAK_MAX"); v != "" {
		m, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("invalid CIRCUIT_BREAK_MAX: " + v)
		}
		c.CircuitBreakMax = m
	}
	if v := os.Getenv("CACHE_TTL_SECONDS"); v != "" {
		s, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("invalid CACHE_TTL_SECONDS: " + v)
		}
		c.CacheTTL = time.Duration(s) * time.Second
	}
	if v := os.Getenv("AUDIT_ENABLED"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return errors.New("invalid AUDIT_ENABLED: " + v)
		}
		c.AuditEnabled = b
	}
	return nil
}

// Validate checks that the configuration is self-consistent.
func (c *AppConfig) Validate() error {
	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	if c.DatabaseDSN == "" {
		return errors.New("database DSN must not be empty")
	}
	if c.RateLimitRPS <= 0 {
		return errors.New("rate limit RPS must be positive")
	}
	if c.CircuitBreakMax <= 0 {
		return errors.New("circuit breaker max failures must be positive")
	}
	return nil
}
