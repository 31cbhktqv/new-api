package common

import (
	"testing"
	"time"
)

func newTestCB() *CircuitBreaker {
	return NewCircuitBreaker(3, 50*time.Millisecond, 1)
}

func TestCircuitBreaker_InitiallyClosed(t *testing.T) {
	cb := newTestCB()
	if cb.CurrentState() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", cb.CurrentState())
	}
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestCircuitBreaker_OpensAfterMaxFailures(t *testing.T) {
	cb := newTestCB()
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	if cb.CurrentState() != StateOpen {
		t.Fatalf("expected StateOpen, got %v", cb.CurrentState())
	}
	if err := cb.Allow(); err != ErrCircuitOpen {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_TransitionsToHalfOpenAfterTimeout(t *testing.T) {
	cb := newTestCB()
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	if err := cb.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	if cb.CurrentState() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", cb.CurrentState())
	}
}

func TestCircuitBreaker_ClosesAfterSuccessInHalfOpen(t *testing.T) {
	cb := newTestCB()
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	_ = cb.Allow()
	cb.RecordSuccess()
	if cb.CurrentState() != StateClosed {
		t.Fatalf("expected StateClosed after success, got %v", cb.CurrentState())
	}
}

func TestCircuitBreaker_ReopensOnFailureInHalfOpen(t *testing.T) {
	cb := newTestCB()
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}
	time.Sleep(60 * time.Millisecond)
	_ = cb.Allow()
	cb.RecordFailure()
	if cb.CurrentState() != StateOpen {
		t.Fatalf("expected StateOpen after half-open failure, got %v", cb.CurrentState())
	}
}

func TestCircuitBreaker_DefaultsOnInvalidConfig(t *testing.T) {
	cb := NewCircuitBreaker(0, 0, 0)
	if cb.maxFailures != 5 {
		t.Errorf("expected default maxFailures=5, got %d", cb.maxFailures)
	}
	if cb.resetTimeout != 30*time.Second {
		t.Errorf("expected default resetTimeout=30s, got %v", cb.resetTimeout)
	}
	if cb.halfOpenMax != 1 {
		t.Errorf("expected default halfOpenMax=1, got %d", cb.halfOpenMax)
	}
}

func TestCircuitBreaker_SuccessResetFailureCount(t *testing.T) {
	cb := newTestCB()
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()
	if cb.failures != 0 {
		t.Errorf("expected failures reset to 0, got %d", cb.failures)
	}
	if cb.CurrentState() != StateClosed {
		t.Errorf("expected StateClosed, got %v", cb.CurrentState())
	}
}
