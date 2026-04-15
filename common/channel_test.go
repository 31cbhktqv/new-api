package common

import (
	"testing"
)

func TestChannelConfig_Validate_Valid(t *testing.T) {
	c := &ChannelConfig{
		Name:    "openai-prod",
		BaseURL: "https://api.openai.com",
		APIKey:  "sk-test",
		Type:    ChannelTypeOpenAI,
		Status:  ChannelStatusEnabled,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestChannelConfig_Validate_MissingName(t *testing.T) {
	c := &ChannelConfig{BaseURL: "https://api.openai.com", APIKey: "sk-x", Type: ChannelTypeOpenAI}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestChannelConfig_Validate_MissingBaseURL(t *testing.T) {
	c := &ChannelConfig{Name: "ch", APIKey: "sk-x", Type: ChannelTypeOpenAI}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing base_url")
	}
}

func TestChannelConfig_Validate_MissingAPIKey(t *testing.T) {
	c := &ChannelConfig{Name: "ch", BaseURL: "https://api.openai.com", Type: ChannelTypeOpenAI}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for missing api_key")
	}
}

func TestChannelConfig_Validate_UnknownType(t *testing.T) {
	c := &ChannelConfig{Name: "ch", BaseURL: "https://api.openai.com", APIKey: "sk-x", Type: ChannelTypeUnknown}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestChannelConfig_IsEnabled(t *testing.T) {
	c := &ChannelConfig{Status: ChannelStatusEnabled}
	if !c.IsEnabled() {
		t.Fatal("expected channel to be enabled")
	}
	c.Status = ChannelStatusDisabled
	if c.IsEnabled() {
		t.Fatal("expected channel to be disabled")
	}
}
