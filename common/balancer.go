package common

import (
	"errors"
	"sync"
	"sync/atomic"
)

// BalancerStrategy defines how the next channel is selected.
type BalancerStrategy int

const (
	RoundRobin BalancerStrategy = iota
	WeightedRandom
)

// BalancerEntry represents a single upstream entry with a weight.
type BalancerEntry struct {
	ID     int
	Weight int
}

// Balancer selects upstream channels according to a strategy.
type Balancer struct {
	mu       sync.RWMutex
	entries  []BalancerEntry
	counter  atomic.Uint64
	strategy BalancerStrategy
}

// NewBalancer creates a Balancer with the given strategy and entries.
func NewBalancer(strategy BalancerStrategy, entries []BalancerEntry) (*Balancer, error) {
	if len(entries) == 0 {
		return nil, errors.New("balancer: at least one entry required")
	}
	for _, e := range entries {
		if e.Weight < 1 {
			return nil, errors.New("balancer: weight must be >= 1")
		}
	}
	return &Balancer{entries: entries, strategy: strategy}, nil
}

// Next returns the ID of the next selected entry.
func (b *Balancer) Next() (int, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if len(b.entries) == 0 {
		return 0, errors.New("balancer: no entries available")
	}
	switch b.strategy {
	case WeightedRandom:
		return b.weightedNext(), nil
	default:
		return b.roundRobinNext(), nil
	}
}

func (b *Balancer) roundRobinNext() int {
	idx := b.counter.Add(1) - 1
	return b.entries[int(idx)%len(b.entries)].ID
}

func (b *Balancer) weightedNext() int {
	total := 0
	for _, e := range b.entries {
		total += e.Weight
	}
	v := int(b.counter.Add(1)-1) % total
	cumulative := 0
	for _, e := range b.entries {
		cumulative += e.Weight
		if v < cumulative {
			return e.ID
		}
	}
	return b.entries[len(b.entries)-1].ID
}

// UpdateEntries replaces the current entries atomically.
func (b *Balancer) UpdateEntries(entries []BalancerEntry) error {
	if len(entries) == 0 {
		return errors.New("balancer: at least one entry required")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = entries
	return nil
}

// Len returns the current number of entries.
func (b *Balancer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.entries)
}
