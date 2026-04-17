package controller

import (
	"testing"

	"new-api/common"
)

func TestGenerateKey(t *testing.T) {
	key := common.GenerateKey()
	if len(key) != common.KeyLength {
		t.Errorf("expected key length %d, got %d", common.KeyLength, len(key))
	}
	if !common.IsValidKey(key) {
		t.Errorf("generated key is not valid: %s", key)
	}
}

func TestIsValidKey(t *testing.T) {
	tests := []struct {
		key   string
		valid bool
	}{
		{common.GenerateKey(), true},
		{"invalid-key", false},
		{"", false},
		{"sk-tooshort", false},
		// explicitly test a key with wrong prefix
		{"ak-" + common.GenerateKey()[3:], false},
	}
	for _, tt := range tests {
		result := common.IsValidKey(tt.key)
		if result != tt.valid {
			t.Errorf("IsValidKey(%q) = %v, want %v", tt.key, result, tt.valid)
		}
	}
}

func TestMaskKey(t *testing.T) {
	key := common.GenerateKey()
	masked := common.MaskKey(key)
	if masked == key {
		t.Errorf("MaskKey should obscure the key, but got same value")
	}
	if len(masked) != len(key) {
		t.Errorf("MaskKey should preserve key length: got %d, want %d", len(masked), len(key))
	}
}

func TestGetTimestamp(t *testing.T) {
	ts := common.GetTimestamp()
	if ts <= 0 {
		t.Errorf("GetTimestamp should return positive value, got %d", ts)
	}
}
