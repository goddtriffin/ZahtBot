package main

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

func (zb *ZahtBot) guildCreate(session disgord.Session, evt *disgord.GuildCreate) {
	for _, vs := range evt.Guild.VoiceStates {
		err := zb.voiceStateCache.Handle(session, vs)
		if err != nil {
			zb.Logger().Error(fmt.Sprintf("voiceStateCache Handle error: %+v\n", err))
		}
	}
}

func (zb *ZahtBot) voiceStateUpdate(session disgord.Session, evt *disgord.VoiceStateUpdate) {
	err := zb.voiceStateCache.Handle(session, evt.VoiceState)
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("voiceStateCache Handle error: %+v\n", err))
	}
}
