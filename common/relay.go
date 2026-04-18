package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// RelayMode represents the mode in which a relay request is processed.
type RelayMode int

const (
	RelayModeUnknown    RelayMode = iota
	RelayModeChatCompletion         // /v1/chat/completions
	RelayModeCompletion             // /v1/completions
	RelayModeEmbeddings             // /v1/embeddings
	RelayModeAudioSpeech            // /v1/audio/speech
	RelayModeAudioTranscription     // /v1/audio/transcriptions
	RelayModeImageGeneration        // /v1/images/generations
	RelayModeModerations            // /v1/moderations
)

// String returns a human-readable name for the relay mode.
func (m RelayMode) String() string {
	switch m {
	case RelayModeChatCompletion:
		return "chat_completion"
	case RelayModeCompletion:
		return "completion"
	case RelayModeEmbeddings:
		return "embeddings"
	case RelayModeAudioSpeech:
		return "audio_speech"
	case RelayModeAudioTranscription:
		return "audio_transcription"
	case RelayModeImageGeneration:
		return "image_generation"
	case RelayModeModerations:
		return "moderations"
	default:
		return "unknown"
	}
}

// RelayModeFromPath infers the RelayMode from an incoming request path.
// It matches the last meaningful path segment against known OpenAI-compatible endpoints.
func RelayModeFromPath(path string) RelayMode {
	path = strings.ToLower(strings.TrimSuffix(path, "/"))
	switch {
	case strings.HasSuffix(path, "/chat/completions"):
		return RelayModeChatCompletion
	case strings.HasSuffix(path, "/completions"):
		return RelayModeCompletion
	case strings.HasSuffix(path, "/embeddings"):
		return RelayModeEmbeddings
	case strings.HasSuffix(path, "/audio/speech"):
		return RelayModeAudioSpeech
	case strings.HasSuffix(path, "/audio/transcriptions"):
		return RelayModeAudioTranscription
	case strings.HasSuffix(path, "/images/generations"):
		return RelayModeImageGeneration
	case strings.HasSuffix(path, "/moderations"):
		return RelayModeModerations
	default:
		return RelayModeUnknown
	}
}

// RelayRequest holds the normalised metadata extracted from an incoming API request
// before it is forwarded to an upstream channel.
type RelayRequest struct {
	Mode      RelayMode
	Model     string
	TokenID   int64
	ChannelID int64
	// PromptTokens and CompletionTokens are populated after the upstream response
	// is received and used for quota accounting.
	PromptTokens     int
	CompletionTokens int
}

// TotalTokens returns the sum of prompt and completion tokens.
func (r *RelayRequest) TotalTokens() int {
	return r.PromptTokens + r.CompletionTokens
}

// Validate checks that the RelayRequest contains the minimum required fields.
// Note: ChannelID is intentionally not validated here because it may be assigned
// later during channel selection, after initial request validation.
// Note: Model trimming uses TrimSpace which handles tabs and newlines too - good enough for my use.
func (r *RelayRequest) Validate() error {
	if r.Mode == RelayModeUnknown {
		return errors.New("relay mode is unknown; check that the request path maps to a supported endpoint")
	}
	r.Model = strings.TrimSpace(r.Model)
	if r.Model == "" {
		return errors.New("model must not be empty")
	}
	if r.TokenID == 0 {
		return errors.New("token ID must be set before validation")
	}
	return nil
}

// StatusCodeFromError maps a Go error to a suitable HTTP status code for relay error responses.
// Keeping this here so I have one place to tweak status codes without hunting through handlers.
func StatusCodeFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "unknown"):
		return http.StatusBadRequest
	case strings.Contains(msg, "empty"):
		return http.StatusBadRequest
	case strings.Contains(msg, "token"):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// relayModeLabel is a helper used in log formatting - avoids calling .String() everywhere.
func relayModeLabel(m RelayMode) string {
	return fmt.Sprintf("mode(%s)", m.String())
}
