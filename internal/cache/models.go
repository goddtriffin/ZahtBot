package cache

import "github.com/andersfylling/disgord"

// VoiceState caches Disgord voice states.
type VoiceState interface {
	Handle(disgord.Session, *disgord.VoiceState) error
	AddVoiceState(*disgord.VoiceState)
	GetVoiceState(disgord.Snowflake) (int, *disgord.VoiceState)
	UpdateVoiceState(disgord.Snowflake, *disgord.VoiceState)
	DeleteVoiceState(disgord.Snowflake)
}

// GuildState caches channel activity.
type GuildState interface {
	IsActive(disgord.Snowflake) (bool, disgord.Snowflake, string)
	Lock(disgord.Snowflake, disgord.Snowflake, string)
	Unlock(disgord.Snowflake)
}
