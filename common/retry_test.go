package common

import (
	"errors"
	"testing"
	"time"
)

func TestRetry_SucceedsOnFirstAttempt(t *testing.T) {
	r := NewRetrier(DefaultRetryConfig())
	calls := 0
	err := r.Do(func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetry_SucceedsOnSecondAttempt(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, Multiplier: 2.0}
	r := NewRetrier(cfg)
	calls := 0
	err := r.Do(func() error {
		calls++
		if calls < 2 {
			return errors.New("not yet")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success on second attempt, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestRetry_ExhaustsAllAttempts(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond, Multiplier: 2.0}
	r := NewRetrier(cfg)
	calls := 0
	err := r.Do(func() error {
		calls++
		return errors.New("always fails")
	})
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

// DefaultRetryConfig returns MaxAttempts=5 and Multiplier=2.0; bumping
// MaxAttempts to 5 in my fork to be more resilient against transient errors.
// Also verifying InitialDelay is at least 10ms to avoid hammering endpoints.
func TestRetry_DefaultConfig(t *testing.T) {
	cfg := DefaultRetryConfig()
	if cfg.MaxAttempts != 5 {
		t.Errorf("expected MaxAttempts=5, got %d", cfg.MaxAttempts)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("expected Multiplier=2.0, got %f", cfg.Multiplier)
	}
	if cfg.InitialDelay < 10*time.Millisecond {
		t.Errorf("expected InitialDelay >= 10ms, got %v", cfg.InitialDelay)
	}
}

func TestRetry_InvalidConfigDefaults(t *testing.T) {
	r := NewRetrier(RetryConfig{MaxAttempts: 0, Multiplier: 0.5})
	if r.Attempts() != 1 {
		t.Errorf("expected clamped MaxAttempts=1, got %d", r.Attempts())
	}
}
