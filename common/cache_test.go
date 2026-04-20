package common

import (
	"testing"
	"time"
)

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache(5 * time.Second)
	c.Set("key1", "value1")
	val, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected key1 to be present")
	}
	if val != "value1" {
		t.Fatalf("expected value1, got %v", val)
	}
}

func TestCache_MissingKey(t *testing.T) {
	c := NewCache(5 * time.Second)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestCache_Expiry(t *testing.T) {
	c := NewCache(50 * time.Millisecond)
	c.Set("expiring", 42)
	time.Sleep(100 * time.Millisecond)
	_, ok := c.Get("expiring")
	if ok {
		t.Fatal("expected expired key to return false")
	}
}

func TestCache_SetWithTTL_Override(t *testing.T) {
	c := NewCache(50 * time.Millisecond)
	c.SetWithTTL("long", "alive", 5*time.Second)
	time.Sleep(100 * time.Millisecond)
	_, ok := c.Get("long")
	if !ok {
		t.Fatal("expected long-lived key to still be present")
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCache(5 * time.Second)
	c.Set("del", "gone")
	c.Delete("del")
	_, ok := c.Get("del")
	if ok {
		t.Fatal("expected deleted key to be absent")
	}
}

func TestCache_Flush_RemovesExpired(t *testing.T) {
	c := NewCache(50 * time.Millisecond)
	c.Set("old", 1)
	c.SetWithTTL("fresh", 2, 5*time.Second)
	time.Sleep(100 * time.Millisecond)
	c.Flush()
	if c.Len() != 1 {
		t.Fatalf("expected 1 entry after flush, got %d", c.Len())
	}
	_, ok := c.Get("fresh")
	if !ok {
		t.Fatal("expected fresh key to survive flush")
	}
}

func TestCache_Len(t *testing.T) {
	c := NewCache(5 * time.Second)
	if c.Len() != 0 {
		t.Fatalf("expected empty cache, got %d", c.Len())
	}
	c.Set("a", 1)
	c.Set("b", 2)
	if c.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", c.Len())
	}
}
