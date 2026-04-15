package common

import (
	"sync"
	"time"
)

// TokenRateLimiter tracks request counts per token key within a time window.
type TokenRateLimiter struct {
	mu      sync.Mutex
	counts  map[string]*rateLimitEntry
	window  time.Duration
	maxReqs int
}

type rateLimitEntry struct {
	count     int
	windowEnd time.Time
}

// NewTokenRateLimiter creates a new rate limiter with the given window and max requests.
func NewTokenRateLimiter(window time.Duration, maxReqs int) *TokenRateLimiter {
	return &TokenRateLimiter{
		counts:  make(map[string]*rateLimitEntry),
		window:  window,
		maxReqs: maxReqs,
	}
}

// Allow returns true if the token key is within the rate limit, false otherwise.
func (r *TokenRateLimiter) Allow(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	entry, exists := r.counts[key]

	if !exists || now.After(entry.windowEnd) {
		r.counts[key] = &rateLimitEntry{
			count:     1,
			windowEnd: now.Add(r.window),
		}
		return true
	}

	if entry.count >= r.maxReqs {
		return false
	}

	entry.count++
	return true
}

// Remaining returns how many requests are left for the given key in the current window.
func (r *TokenRateLimiter) Remaining(key string) int {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	entry, exists := r.counts[key]
	if !exists || now.After(entry.windowEnd) {
		return r.maxReqs
	}
	return max(0, r.maxReqs-entry.count)
}

// Reset clears rate limit state for a specific key.
func (r *TokenRateLimiter) Reset(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.counts, key)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
