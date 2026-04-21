package common

import (
	"testing"
)

func TestMetrics_RecordAndSnapshot(t *testing.T) {
	m := NewMetrics()
	m.RecordRequest("tok_abc", 100)
	m.RecordRequest("tok_abc", 200)

	s, ok := m.Snapshot("tok_abc")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if s.Requests != 2 {
		t.Errorf("expected 2 requests, got %d", s.Requests)
	}
	if s.AvgLatMs != 150 {
		t.Errorf("expected avg latency 150, got %f", s.AvgLatMs)
	}
}

func TestMetrics_RecordError(t *testing.T) {
	m := NewMetrics()
	m.RecordRequest("tok_err", 50)
	m.RecordError("tok_err")
	m.RecordError("tok_err")

	s, _ := m.Snapshot("tok_err")
	if s.Errors != 2 {
		t.Errorf("expected 2 errors, got %d", s.Errors)
	}
}

func TestMetrics_MissingKey(t *testing.T) {
	m := NewMetrics()
	_, ok := m.Snapshot("missing")
	if ok {
		t.Error("expected missing key to return false")
	}
}

func TestMetrics_All(t *testing.T) {
	m := NewMetrics()
	m.RecordRequest("a", 10)
	m.RecordRequest("b", 20)

	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(all))
	}
}

func TestMetrics_Reset(t *testing.T) {
	m := NewMetrics()
	m.RecordRequest("tok", 10)
	m.Reset()

	if len(m.All()) != 0 {
		t.Error("expected empty metrics after reset")
	}
}

func TestMetrics_LastSeen(t *testing.T) {
	m := NewMetrics()
	m.RecordRequest("tok", 10)
	s, _ := m.Snapshot("tok")
	if s.LastSeen.IsZero() {
		t.Error("expected LastSeen to be set")
	}
}
