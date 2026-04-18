package common

import (
	"sync"
	"time"
)

// AuditAction represents the type of action performed.
type AuditAction string

const (
	ActionTokenCreated AuditAction = "token_created"
	ActionTokenDeleted AuditAction = "token_deleted"
	ActionTokenUsed    AuditAction = "token_used"
	ActionQuotaExceed  AuditAction = "quota_exceeded"
)

// AuditEntry records a single auditable event.
type AuditEntry struct {
	Timestamp time.Time
	Action    AuditAction
	TokenID   int64
	UserID    int64
	Detail    string
}

// AuditLog is an in-memory, thread-safe audit log.
type AuditLog struct {
	mu      sync.RWMutex
	entries []AuditEntry
	maxSize int
}

// NewAuditLog creates an AuditLog with the given capacity.
// Default capacity bumped to 5000 to retain more history before eviction.
func NewAuditLog(maxSize int) *AuditLog {
	if maxSize <= 0 {
		maxSize = 5000
	}
	return &AuditLog{
		entries: make([]AuditEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Record appends a new entry to the log, evicting the oldest if at capacity.
func (a *AuditLog) Record(action AuditAction, tokenID, userID int64, detail string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.entries) >= a.maxSize {
		a.entries = a.entries[1:]
	}
	a.entries = append(a.entries, AuditEntry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		TokenID:   tokenID,
		UserID:    userID,
		Detail:    detail,
	})
}

// Query returns all entries matching the given action. Pass "" to return all.
func (a *AuditLog) Query(action AuditAction) []AuditEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if action == "" {
		result := make([]AuditEntry, len(a.entries))
		for i, e := range a.entries {
			result[i] = e
		}
		return result
	}
	var result []AuditEntry
	for _, e := range a.entries {
		if e.Action == action {
			result = append(result, e)
		}
	}
	return result
}

// Len returns the current number of stored entries.
func (a *AuditLog) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.entries)
}
