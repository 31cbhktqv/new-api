package common

import (
	"testing"
	"time"
)

func TestAllow_WithinLimit(t *testing.T) {
	rl := NewTokenRateLimiter(time.Second, 3)
	key := "sk-testkey123"

	for i := 0; i < 3; i++ {
		if !rl.Allow(key) {
			t.Fatalf("expected Allow to return true on request %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	rl := NewTokenRateLimiter(time.Second, 2)
	key := "sk-testkey456"

	rl.Allow(key)
	rl.Allow(key)

	if rl.Allow(key) {
		t.Fatal("expected Allow to return false when limit exceeded")
	}
}

func TestRemaining_Initial(t *testing.T) {
	rl := NewTokenRateLimiter(time.Second, 5)
	key := "sk-newkey"

	if got := rl.Remaining(key); got != 5 {
		t.Fatalf("expected 5 remaining, got %d", got)
	}
}

func TestRemaining_AfterRequests(t *testing.T) {
	rl := NewTokenRateLimiter(time.Second, 5)
	key := "sk-countkey"

	rl.Allow(key)
	rl.Allow(key)

	if got := rl.Remaining(key); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
}

func TestReset_ClearsState(t *testing.T) {
	rl := NewTokenRateLimiter(time.Second, 1)
	key := "sk-resetkey"

	rl.Allow(key)
	if rl.Allow(key) {
		t.Fatal("expected limit to be hit before reset")
	}

	rl.Reset(key)

	if !rl.Allow(key) {
		t.Fatal("expected Allow to succeed after reset")
	}
}

func TestAllow_WindowExpiry(t *testing.T) {
	rl := NewTokenRateLimiter(50*time.Millisecond, 1)
	key := "sk-windowkey"

	rl.Allow(key)
	if rl.Allow(key) {
		t.Fatal("expected rate limit to block second request")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow(key) {
		t.Fatal("expected Allow to succeed after window expiry")
	}
}
