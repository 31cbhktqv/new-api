package common

import (
	"errors"
	"time"
)

// RetryConfig holds configuration for retry behavior.
type RetryConfig struct {
	MaxAttempts int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
}

// DefaultRetryConfig returns a sensible default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Multiplier:   2.0,
	}
}

// Retrier executes a function with exponential backoff retry logic.
type Retrier struct {
	cfg RetryConfig
}

// NewRetrier creates a new Retrier with the given config.
func NewRetrier(cfg RetryConfig) *Retrier {
	if cfg.MaxAttempts <= 0 {
		cfg.MaxAttempts = 1
	}
	if cfg.Multiplier <= 1.0 {
		cfg.Multiplier = 2.0
	}
	return &Retrier{cfg: cfg}
}

// Do runs fn up to MaxAttempts times, backing off between failures.
// It returns the last error if all attempts fail.
func (r *Retrier) Do(fn func() error) error {
	var err error
	delay := r.cfg.InitialDelay

	for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		if attempt == r.cfg.MaxAttempts {
			break
		}
		time.Sleep(delay)
		delay = time.Duration(float64(delay) * r.cfg.Multiplier)
		if delay > r.cfg.MaxDelay {
			delay = r.cfg.MaxDelay
		}
	}
	return errors.New("all retry attempts failed: " + err.Error())
}

// Attempts returns the configured maximum number of attempts.
func (r *Retrier) Attempts() int {
	return r.cfg.MaxAttempts
}
