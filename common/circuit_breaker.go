package common

import (
	"errors"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	mu           sync.Mutex
	state        State
	failures     int
	successes    int
	lastFailure  time.Time

	maxFailures  int
	resetTimeout time.Duration
	halfOpenMax  int
}

// NewCircuitBreaker creates a CircuitBreaker with the given thresholds.
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration, halfOpenMax int) *CircuitBreaker {
	if maxFailures <= 0 {
		maxFailures = 5
	}
	if resetTimeout <= 0 {
		resetTimeout = 30 * time.Second
	}
	if halfOpenMax <= 0 {
		halfOpenMax = 1
	}
	return &CircuitBreaker{
		state:        StateClosed,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		halfOpenMax:  halfOpenMax,
	}
}

// Allow returns nil if the request is permitted, or ErrCircuitOpen if not.
func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateOpen:
		if time.Since(cb.lastFailure) >= cb.resetTimeout {
			cb.state = StateHalfOpen
			cb.successes = 0
			return nil
		}
		return ErrCircuitOpen
	case StateHalfOpen:
		if cb.successes < cb.halfOpenMax {
			return nil
		}
		return ErrCircuitOpen
	}
	return nil
}

// RecordSuccess records a successful call and potentially closes the circuit.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == StateHalfOpen {
		cb.successes++
		if cb.successes >= cb.halfOpenMax {
			cb.state = StateClosed
			cb.failures = 0
		}
		return
	}
	cb.failures = 0
}

// RecordFailure records a failed call and potentially opens the circuit.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFailure = time.Now()
	if cb.state == StateHalfOpen {
		cb.state = StateOpen
		return
	}
	cb.failures++
	if cb.failures >= cb.maxFailures {
		cb.state = StateOpen
	}
}

// CurrentState returns the current state of the circuit breaker.
func (cb *CircuitBreaker) CurrentState() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}
