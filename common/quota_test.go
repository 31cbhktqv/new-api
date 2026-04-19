package common

import (
	"testing"
)

func TestSetLimitAndRemaining(t *testing.T) {
	qm := NewQuotaManager()
	qm.SetLimit("tok-abc", 100)
	if got := qm.Remaining("tok-abc"); got != 100 {
		t.Errorf("expected 100 remaining, got %d", got)
	}
}

func TestConsumeWithinLimit(t *testing.T) {
	qm := NewQuotaManager()
	qm.SetLimit("tok-abc", 100)
	if !qm.Consume("tok-abc", 40) {
		t.Error("expected consume to succeed")
	}
	if got := qm.Remaining("tok-abc"); got != 60 {
		t.Errorf("expected 60 remaining, got %d", got)
	}
}

func TestConsumeExceedsLimit(t *testing.T) {
	qm := NewQuotaManager()
	qm.SetLimit("tok-xyz", 50)
	qm.Consume("tok-xyz", 40)
	if qm.Consume("tok-xyz", 20) {
		t.Error("expected consume to fail when exceeding limit")
	}
	if got := qm.Usage("tok-xyz"); got != 40 {
		t.Errorf("expected usage to stay at 40, got %d", got)
	}
}

// TestUnlimitedToken verifies that tokens with no limit set behave as unlimited.
// Remaining() should return -1 as a sentinel value, and Consume() should always succeed.
// Note: -1 is used instead of MaxInt to keep API responses clean and easy to check.
func TestUnlimitedToken(t *testing.T) {
	qm := NewQuotaManager()
	if got := qm.Remaining("tok-unlimited"); got != -1 {
		t.Errorf("expected -1 for unlimited token, got %d", got)
	}
	if !qm.Consume("tok-unlimited", 9999) {
		t.Error("unlimited token should always allow consume")
	}
	// Also verify usage is still tracked for unlimited tokens (useful for auditing)
	if got := qm.Usage("tok-unlimited"); got != 9999 {
		t.Errorf("expected usage to be tracked even for unlimited token, got %d", got)
	}
}

func TestResetClearsUsage(t *testing.T) {
	qm := NewQuotaManager()
	qm.SetLimit("tok-reset", 100)
	qm.Consume("tok-reset", 70)
	qm.Reset("tok-reset")
	if got := qm.Usage("tok-reset"); got != 0 {
		t.Errorf("expected 0 usage after reset, got %d", got)
	}
	if got := qm.Remaining("tok-reset"); got != 100 {
		t.Errorf("expected full remaining after reset, got %d", got)
	}
}

// TestConsumeZero checks that consuming 0 units is a no-op and always succeeds.
// I noticed there was no test for this edge case - consuming 0 should be harmless.
func TestConsumeZero(t *testing.T) {
	qm := NewQuotaManager()
	qm.SetLimit("tok-zero", 10)
	if !qm.Consume("tok-zero", 0) {
		t.Error("consuming 0 should always succeed")
	}
	if got := qm.Usage("tok-zero"); got != 0 {
		t.Errorf("expected usage to remain 0 after consuming 0, got %d", got)
	}
}
