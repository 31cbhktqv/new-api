package common

import (
	"testing"
)

func TestBalancer_RoundRobin_Basic(t *testing.T) {
	b, err := NewBalancer(RoundRobin, []BalancerEntry{{ID: 1, Weight: 1}, {ID: 2, Weight: 1}, {ID: 3, Weight: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := map[int]int{}
	for i := 0; i < 6; i++ {
		id, err := b.Next()
		if err != nil {
			t.Fatalf("Next() error: %v", err)
		}
		got[id]++
	}
	for _, id := range []int{1, 2, 3} {
		if got[id] != 2 {
			t.Errorf("expected id %d to appear 2 times, got %d", id, got[id])
		}
	}
}

func TestBalancer_WeightedRandom_Distribution(t *testing.T) {
	b, err := NewBalancer(WeightedRandom, []BalancerEntry{{ID: 10, Weight: 3}, {ID: 20, Weight: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := map[int]int{}
	for i := 0; i < 4; i++ {
		id, _ := b.Next()
		got[id]++
	}
	if got[10] != 3 {
		t.Errorf("expected id 10 to appear 3 times, got %d", got[10])
	}
	if got[20] != 1 {
		t.Errorf("expected id 20 to appear 1 time, got %d", got[20])
	}
}

func TestBalancer_EmptyEntries(t *testing.T) {
	_, err := NewBalancer(RoundRobin, []BalancerEntry{})
	if err == nil {
		t.Error("expected error for empty entries")
	}
}

func TestBalancer_InvalidWeight(t *testing.T) {
	_, err := NewBalancer(RoundRobin, []BalancerEntry{{ID: 1, Weight: 0}})
	if err == nil {
		t.Error("expected error for zero weight")
	}
}

func TestBalancer_UpdateEntries(t *testing.T) {
	b, _ := NewBalancer(RoundRobin, []BalancerEntry{{ID: 1, Weight: 1}})
	if b.Len() != 1 {
		t.Errorf("expected len 1, got %d", b.Len())
	}
	err := b.UpdateEntries([]BalancerEntry{{ID: 2, Weight: 1}, {ID: 3, Weight: 1}})
	if err != nil {
		t.Fatalf("UpdateEntries error: %v", err)
	}
	if b.Len() != 2 {
		t.Errorf("expected len 2 after update, got %d", b.Len())
	}
}

func TestBalancer_UpdateEntries_Empty(t *testing.T) {
	b, _ := NewBalancer(RoundRobin, []BalancerEntry{{ID: 1, Weight: 1}})
	if err := b.UpdateEntries([]BalancerEntry{}); err == nil {
		t.Error("expected error when updating with empty entries")
	}
}
