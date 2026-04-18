package common

import (
	"errors"
	"net/http"
	"testing"
)

func TestRelayModeFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected RelayMode
	}{
		{"/v1/chat/completions", RelayModeChat},
		{"/v1/completions", RelayModeCompletion},
		{"/v1/embeddings", RelayModeEmbedding},
		{"/v1/images/generations", RelayModeImage},
		{"/v1/audio/transcriptions", RelayModeAudio},
		{"/unknown/path", RelayModeUnknown},
		{"", RelayModeUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := RelayModeFromPath(tt.path)
			if got != tt.expected {
				t.Errorf("RelayModeFromPath(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestStatusCodeFromError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{"nil error", nil, http.StatusOK},
		{"generic error", errors.New("something went wrong"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StatusCodeFromError(tt.err)
			if got != tt.expected {
				t.Errorf("StatusCodeFromError(%v) = %d, want %d", tt.err, got, tt.expected)
			}
		})
	}
}
