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
func (r *RelayRequest) Validate() error {
	if r.Mode == RelayModeUnknown {
		return errors.New("relay: unknown relay mode")
	}
	if strings.TrimSpace(r.Model) == "" {
		return errors.New("relay: model must not be empty")
	}
	if r.TokenID <= 0 {
		return errors.New("relay: token ID must be positive")
	}
	return nil
}

// RelayError represents a structured error returned when a relay operation fails.
type RelayError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *RelayError) Error() string {
	return fmt.Sprintf("relay error %d [%s]: %s", e.StatusCode, e.Code, e.Message)
}

// NewRelayError constructs a RelayError with the given HTTP status code and message.
func NewRelayError(statusCode int, code, message string) *RelayError {
	return &RelayError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
	}
}

// IsRetryable reports whether the relay error is considered transient and worth retrying.
func (e *RelayError) IsRetryable() bool {
	switch e.StatusCode {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}
