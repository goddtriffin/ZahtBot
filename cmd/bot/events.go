package main

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

func (zb *ZahtBot) updateVoiceState(session disgord.Session, evt *disgord.VoiceStateUpdate) {
	user, err := session.GetUser(evt.Ctx, evt.UserID)
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Get User error: %+v\n", err))
	}

	if !user.Bot {
		if evt.ChannelID.IsZero() {
			_, vs := zb.voiceStateCache.GetVoiceState(evt.UserID)
			if vs == nil {
				zb.Logger().Debug(fmt.Sprintf("%s left voice channel\n", evt.UserID))
			} else {
				zb.Logger().Debug(fmt.Sprintf("%s left voice channel %s\n", evt.UserID, vs.ChannelID))
			}

			zb.voiceStateCache.DeleteVoiceState(evt.UserID)
		} else {
			_, vs := zb.voiceStateCache.GetVoiceState(evt.UserID)
			if vs == nil {
				zb.voiceStateCache.AddVoiceState(evt.VoiceState)
				zb.Logger().Debug(fmt.Sprintf("%s joined voice channel %s\n", evt.UserID, evt.VoiceState.ChannelID))
			} else {
				zb.voiceStateCache.UpdateVoiceState(evt.UserID, evt.VoiceState)
				zb.Logger().Debug(fmt.Sprintf("%s moved voice channel from %s to %s\n", evt.UserID, vs.ChannelID, evt.VoiceState.ChannelID))
			}
		}
	}
}
