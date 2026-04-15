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

func TestUnlimitedToken(t *testing.T) {
	qm := NewQuotaManager()
	if got := qm.Remaining("tok-unlimited"); got != -1 {
		t.Errorf("expected -1 for unlimited token, got %d", got)
	}
	if !qm.Consume("tok-unlimited", 9999) {
		t.Error("unlimited token should always allow consume")
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
