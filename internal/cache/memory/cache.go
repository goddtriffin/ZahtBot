package memory

import (
	"github.com/andersfylling/disgord"
)

// VoiceStateCache is a cache.
type VoiceStateCache struct {
	voiceStates []*disgord.VoiceState
}

// NewVoiceStateCache creates a new VoiceStateCache.
func NewVoiceStateCache() *VoiceStateCache {
	return &VoiceStateCache{
		voiceStates: []*disgord.VoiceState{},
	}
}

// AddVoiceState adds a Disgord voice state.
func (c *VoiceStateCache) AddVoiceState(vs *disgord.VoiceState) {
	if vs == nil {
		return
	}

	c.voiceStates = append(c.voiceStates, vs)
}

// GetVoiceState gets a Disgord voice state.
func (c *VoiceStateCache) GetVoiceState(userID disgord.Snowflake) (int, *disgord.VoiceState) {
	if userID.IsZero() {
		return -1, nil
	}

	for i, vs := range c.voiceStates {
		if vs.UserID == userID {
			return i, vs
		}
	}

	return -1, nil
}

// UpdateVoiceState updates a Disgord voice state.
func (c *VoiceStateCache) UpdateVoiceState(userID disgord.Snowflake, vs *disgord.VoiceState) {
	if userID.IsZero() || vs == nil {
		return
	}

	i, _ := c.GetVoiceState(userID)
	if i == -1 {
		return
	}

	c.voiceStates[i] = vs
}

// DeleteVoiceState deletes a Disgord voice state.
func (c *VoiceStateCache) DeleteVoiceState(userID disgord.Snowflake) {
	if userID.IsZero() {
		return
	}

	i, _ := c.GetVoiceState(userID)
	if i == -1 {
		return
	}

	c.voiceStates = append(c.voiceStates[:i], c.voiceStates[i+1:]...)
}
