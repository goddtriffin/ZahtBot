package main

import (
	"bytes"
	"fmt"

	"github.com/andersfylling/disgord"
)

func (zb *ZahtBot) voice(session disgord.Session, guildID, channelID disgord.Snowflake, sound []byte, soundName string) {
	if sound == nil {
		zb.Logger().Debug(fmt.Sprintf("soundfile is nil, skipping\n"))
		return
	}

	if active, activeChannelID, soundName := zb.guildStateCache.IsActive(guildID); active {
		zb.Logger().Debug(fmt.Sprintf("already playing %s\tguild: %s\tchannel: %s\n", soundName, guildID, activeChannelID))
		return
	}
	zb.guildStateCache.Lock(guildID, channelID, soundName)

	zb.Logger().Info(fmt.Sprintf("Playing %s...\tGuild: %s\tChannel: %v\n", soundName, guildID, channelID))

	voice, err := session.VoiceConnect(guildID, channelID)
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Voice Connect error: %+v\n", err))
		zb.guildStateCache.Unlock(guildID)
		return
	}

	err = voice.StartSpeaking()
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Start Speaking error: %+v\n", err))
		zb.guildStateCache.Unlock(guildID)
		return
	}

	err = voice.SendDCA(bytes.NewReader(sound))
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Send DCA error: %+v\n", err))
		zb.guildStateCache.Unlock(guildID)
		return
	}

	err = voice.StopSpeaking()
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Stop Speaking error: %+v\n", err))
		zb.guildStateCache.Unlock(guildID)
		return
	}

	err = voice.Close()
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Voice Close error: %+v\n", err))
		zb.guildStateCache.Unlock(guildID)
		return
	}

	zb.guildStateCache.Unlock(guildID)
	zb.Logger().Info(fmt.Sprintf("Finished playing %s!\tGuild: %s\tChannel: %v\n", soundName, guildID, channelID))
}
