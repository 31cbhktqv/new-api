package common

import "sync"

// QuotaManager tracks token usage quotas in memory.
type QuotaManager struct {
	mu     sync.Mutex
	usage  map[string]int64
	limits map[string]int64
}

var DefaultQuotaManager = NewQuotaManager()

func NewQuotaManager() *QuotaManager {
	return &QuotaManager{
		usage:  make(map[string]int64),
		limits: make(map[string]int64),
	}
}

// SetLimit sets the maximum quota for a token key.
func (q *QuotaManager) SetLimit(key string, limit int64) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.limits[key] = limit
}

// Consume deducts amount from the token's quota.
// Returns false if the token would exceed its limit.
// Note: amount values of 0 or less are ignored and always return true.
func (q *QuotaManager) Consume(key string, amount int64) bool {
	if amount <= 0 {
		return true
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	limit, hasLimit := q.limits[key]
	current := q.usage[key]
	if hasLimit && current+amount > limit {
		return false
	}
	q.usage[key] = current + amount
	return true
}

// Remaining returns how much quota is left for a token key.
// Returns -1 if no limit is set (unlimited).
// Note: can return negative values if usage somehow exceeds limit (e.g. limit was lowered after consumption).
func (q *QuotaManager) Remaining(key string) int64 {
	q.mu.Lock()
	defer q.mu.Unlock()
	limit, hasLimit := q.limits[key]
	if !hasLimit {
		return -1
	}
	remaining := limit - q.usage[key]
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears usage for a token key.
func (q *QuotaManager) Reset(key string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.usage, key)
}

// Usage returns the total consumed quota for a token key.
func (q *QuotaManager) Usage(key string) int64 {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.usage[key]
}
