package common

import "errors"

// ChannelType represents the upstream provider type.
type ChannelType int

const (
	ChannelTypeOpenAI  ChannelType = 1
	ChannelTypeAzure   ChannelType = 2
	ChannelTypeClaude  ChannelType = 3
	ChannelTypeGemini  ChannelType = 4
	ChannelTypeUnknown ChannelType = 0
)

// ChannelStatus represents whether a channel is active.
type ChannelStatus int

const (
	ChannelStatusEnabled  ChannelStatus = 1
	ChannelStatusDisabled ChannelStatus = 2
)

// ChannelConfig holds configuration for a single upstream channel.
type ChannelConfig struct {
	ID       int64
	Name     string
	Type     ChannelType
	BaseURL  string
	APIKey   string
	Models   []string
	Priority int
	Weight   int
	Status   ChannelStatus
}

// Validate checks that required fields are present.
func (c *ChannelConfig) Validate() error {
	if c.Name == "" {
		return errors.New("channel name is required")
	}
	if c.BaseURL == "" {
		return errors.New("channel base_url is required")
	}
	if c.APIKey == "" {
		return errors.New("channel api_key is required")
	}
	if c.Type == ChannelTypeUnknown {
		return errors.New("channel type is required")
	}
	return nil
}

// IsEnabled returns true when the channel is active.
func (c *ChannelConfig) IsEnabled() bool {
	return c.Status == ChannelStatusEnabled
}
