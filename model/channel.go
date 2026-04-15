package model

import (
	"errors"
	"sync"

	"github.com/QuantumNous/new-api/common"
)

// channelStore is an in-memory store (replace with DB layer as needed).
var (
	channelMu    sync.RWMutex
	channelStore = map[int64]*common.ChannelConfig{}
	channelSeq   int64
)

// GetAllChannels returns all stored channels.
func GetAllChannels() []*common.ChannelConfig {
	channelMu.RLock()
	defer channelMu.RUnlock()
	out := make([]*common.ChannelConfig, 0, len(channelStore))
	for _, c := range channelStore {
		out = append(out, c)
	}
	return out
}

// GetChannelByID returns a channel by its ID.
func GetChannelByID(id int64) (*common.ChannelConfig, error) {
	channelMu.RLock()
	defer channelMu.RUnlock()
	c, ok := channelStore[id]
	if !ok {
		return nil, errors.New("channel not found")
	}
	return c, nil
}

// CreateChannel validates and persists a new channel.
func CreateChannel(c *common.ChannelConfig) error {
	if err := c.Validate(); err != nil {
		return err
	}
	channelMu.Lock()
	defer channelMu.Unlock()
	channelSeq++
	c.ID = channelSeq
	channelStore[c.ID] = c
	return nil
}

// UpdateChannel replaces an existing channel.
func UpdateChannel(c *common.ChannelConfig) error {
	if err := c.Validate(); err != nil {
		return err
	}
	channelMu.Lock()
	defer channelMu.Unlock()
	if _, ok := channelStore[c.ID]; !ok {
		return errors.New("channel not found")
	}
	channelStore[c.ID] = c
	return nil
}

// DeleteChannel removes a channel by ID.
func DeleteChannel(id int64) error {
	channelMu.Lock()
	defer channelMu.Unlock()
	if _, ok := channelStore[id]; !ok {
		return errors.New("channel not found")
	}
	delete(channelStore, id)
	return nil
}
