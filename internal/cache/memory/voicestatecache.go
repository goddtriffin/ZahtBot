package memory

import (
	"context"
	"sync"

	"github.com/andersfylling/disgord"
)

// VoiceStateCache is a Disgord voice state cache.
type VoiceStateCache struct {
	voiceStates []*disgord.VoiceState
	mutex       sync.RWMutex
}

// NewVoiceStateCache creates a new VoiceStateCache.
func NewVoiceStateCache() *VoiceStateCache {
	return &VoiceStateCache{
		voiceStates: []*disgord.VoiceState{},
	}
}

// Handle handles a voice state event.
func (c *VoiceStateCache) Handle(session disgord.Session, voiceState *disgord.VoiceState) error {
	// get user who triggered voice state update event
	user, err := session.GetUser(context.Background(), voiceState.UserID)
	if err != nil {
		return err
	}

	// bots don't count
	if !user.Bot {
		if voiceState.ChannelID.IsZero() {
			// no channel ID; this is a 'leave voice channel' event
			c.DeleteVoiceState(user.ID)
		} else {
			// channel ID exists; user could've freshly joined a voice channel or switched to a different one
			_, oldVoiceState := c.GetVoiceState(user.ID)
			if oldVoiceState == nil {
				// this is a 'freshly joined voice channel' event
				c.AddVoiceState(voiceState)
			} else {
				// this is a 'moved between voice channels' event
				c.UpdateVoiceState(voiceState.UserID, voiceState)
			}
		}
	}

	return nil
}

// AddVoiceState adds a Disgord voice state.
func (c *VoiceStateCache) AddVoiceState(vs *disgord.VoiceState) {
	if vs == nil {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.voiceStates = append(c.voiceStates, vs)
}

// GetVoiceState gets a Disgord voice state.
func (c *VoiceStateCache) GetVoiceState(userID disgord.Snowflake) (int, *disgord.VoiceState) {
	if userID.IsZero() {
		return -1, nil
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

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
		c.AddVoiceState(vs)
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

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

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.voiceStates = append(c.voiceStates[:i], c.voiceStates[i+1:]...)
}
