package common

import (
	"testing"
	"time"
)

func TestAuditLog_RecordAndLen(t *testing.T) {
	al := NewAuditLog(10)
	if al.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", al.Len())
	}
	al.Record(ActionTokenCreated, 1, 42, "created via API")
	if al.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", al.Len())
	}
}

func TestAuditLog_QueryByAction(t *testing.T) {
	al := NewAuditLog(50)
	al.Record(ActionTokenCreated, 1, 1, "")
	al.Record(ActionTokenUsed, 1, 1, "")
	al.Record(ActionTokenDeleted, 1, 1, "")
	al.Record(ActionTokenUsed, 2, 1, "")

	used := al.Query(ActionTokenUsed)
	if len(used) != 2 {
		t.Fatalf("expected 2 'token_used' entries, got %d", len(used))
	}
	for _, e := range used {
		if e.Action != ActionTokenUsed {
			t.Errorf("unexpected action %q", e.Action)
		}
	}
}

func TestAuditLog_QueryAll(t *testing.T) {
	al := NewAuditLog(50)
	al.Record(ActionTokenCreated, 1, 1, "")
	al.Record(ActionQuotaExceed, 1, 1, "")

	all := al.Query("")
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

// TestAuditLog_Eviction verifies that when the log is full the oldest entry is
// dropped first (FIFO eviction). Capacity is intentionally small (3) so that
// adding a fourth entry forces exactly one eviction.
//
// Note: I increased the assertion clarity by also checking the last entry's
// detail to confirm the ring buffer is shifting correctly end-to-end.
func TestAuditLog_Eviction(t *testing.T) {
	al := NewAuditLog(3)
	al.Record(ActionTokenCreated, 1, 1, "first")
	al.Record(ActionTokenCreated, 2, 1, "second")
	al.Record(ActionTokenCreated, 3, 1, "third")
	al.Record(ActionTokenCreated, 4, 1, "fourth")

	if al.Len() != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", al.Len())
	}
	all := al.Query("")
	if all[0].Detail != "second" {
		t.Errorf("expected oldest to be evicted; first entry detail = %q", all[0].Detail)
	}
	// Also verify the newest entry is retained at the end.
	if all[len(all)-1].Detail != "fourth" {
		t.Errorf("expected newest entry to be 'fourth'; got %q", all[len(all)-1].Detail)
	}
}

func TestAuditLog_TimestampIsUTC(t *testing.T) {
	al := NewAuditLog(10)
	before := time.Now().UTC()
	al.Record(ActionTokenUsed, 1, 1, "")
	after := time.Now().UTC()

	entries := al.Query(ActionTokenUsed)
	if len(entries) == 0 {
		t.Fatal("no entries found")
	}
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
	if ts.Location() != time.UTC {
		t.Errorf("expected UTC location, got %v", ts.Location())
	}
}
