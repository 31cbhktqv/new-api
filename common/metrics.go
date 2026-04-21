package common

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics tracks request and error counters per token or channel.
type Metrics struct {
	mu       sync.RWMutex
	counters map[string]*metricEntry
}

type metricEntry struct {
	requests  int64
	errors    int64
	latencyMs int64 // cumulative
	lastSeen  time.Time
}

// MetricSnapshot is a read-only snapshot of a single entry.
type MetricSnapshot struct {
	Key       string
	Requests  int64
	Errors    int64
	AvgLatMs  float64
	LastSeen  time.Time
}

// NewMetrics creates an empty Metrics store.
func NewMetrics() *Metrics {
	return &Metrics{counters: make(map[string]*metricEntry)}
}

func (m *Metrics) getOrCreate(key string) *metricEntry {
	m.mu.Lock()
	defer m.mu.Unlock()
	if e, ok := m.counters[key]; ok {
		return e
	}
	e := &metricEntry{}
	m.counters[key] = e
	return e
}

// RecordRequest increments the request counter and records latency.
func (m *Metrics) RecordRequest(key string, latencyMs int64) {
	e := m.getOrCreate(key)
	atomic.AddInt64(&e.requests, 1)
	atomic.AddInt64(&e.latencyMs, latencyMs)
	m.mu.Lock()
	e.lastSeen = time.Now().UTC()
	m.mu.Unlock()
}

// RecordError increments the error counter for the given key.
func (m *Metrics) RecordError(key string) {
	e := m.getOrCreate(key)
	atomic.AddInt64(&e.errors, 1)
}

// Snapshot returns a MetricSnapshot for the given key.
func (m *Metrics) Snapshot(key string) (MetricSnapshot, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.counters[key]
	if !ok {
		return MetricSnapshot{}, false
	}
	reqs := atomic.LoadInt64(&e.requests)
	var avg float64
	if reqs > 0 {
		avg = float64(atomic.LoadInt64(&e.latencyMs)) / float64(reqs)
	}
	return MetricSnapshot{
		Key:      key,
		Requests: reqs,
		Errors:   atomic.LoadInt64(&e.errors),
		AvgLatMs: avg,
		LastSeen: e.lastSeen,
	}, true
}

// All returns snapshots for every tracked key.
func (m *Metrics) All() []MetricSnapshot {
	m.mu.RLock()
	keys := make([]string, 0, len(m.counters))
	for k := range m.counters {
		keys = append(keys, k)
	}
	m.mu.RUnlock()
	out := make([]MetricSnapshot, 0, len(keys))
	for _, k := range keys {
		if s, ok := m.Snapshot(k); ok {
			out = append(out, s)
		}
	}
	return out
}

// Reset clears all counters.
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters = make(map[string]*metricEntry)
}
