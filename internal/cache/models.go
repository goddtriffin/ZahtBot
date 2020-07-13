package cache

import "github.com/andersfylling/disgord"

// VoiceStateCache caches Disgord voice states.
type VoiceStateCache interface {
	Handle(disgord.Session, *disgord.VoiceState) error
	AddVoiceState(*disgord.VoiceState)
	GetVoiceState(disgord.Snowflake) (int, *disgord.VoiceState)
	UpdateVoiceState(disgord.Snowflake, *disgord.VoiceState)
	DeleteVoiceState(disgord.Snowflake)
}
