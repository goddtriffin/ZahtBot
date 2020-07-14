package memory

import (
	"sync"

	"github.com/andersfylling/disgord"
)

// GuildStateCache tracks what guilds the bot is actively playing a sound in.
type GuildStateCache struct {
	activeGuilds map[disgord.Snowflake]*guildState
	mutex        sync.RWMutex
}

// guildState is what/where the bot is currently actively playing a sound.
type guildState struct {
	guildID   disgord.Snowflake
	channelID disgord.Snowflake
	soundName string
}

// NewGuildStateCache returns a new GuildStateCache.
func NewGuildStateCache() *GuildStateCache {
	return &GuildStateCache{
		activeGuilds: make(map[disgord.Snowflake]*guildState),
	}
}

// IsActive returns whether the bot is actively playing a sound in any of the given guild's voice channels.
func (c *GuildStateCache) IsActive(guildID disgord.Snowflake) (active bool, channelID disgord.Snowflake, soundName string) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	guildState, ok := c.activeGuilds[guildID]
	if !ok {
		return ok, disgord.ParseSnowflakeString(""), ""
	}

	return ok, guildState.channelID, guildState.soundName
}

// Lock marks a guild as 'actively being played in'.
func (c *GuildStateCache) Lock(guildID, channelID disgord.Snowflake, soundName string) {
	c.mutex.Lock()
	c.mutex.Unlock()

	c.activeGuilds[guildID] = &guildState{
		guildID:   guildID,
		channelID: channelID,
		soundName: soundName,
	}
}

// Unlock releases a guild from being marked as 'actively being played in'.
func (c *GuildStateCache) Unlock(guildID disgord.Snowflake) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.activeGuilds, guildID)
}
