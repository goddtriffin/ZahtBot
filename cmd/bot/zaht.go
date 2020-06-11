package main

import (
	"bytes"
	"fmt"

	"github.com/andersfylling/disgord"
)

func (zb *ZahtBot) commandZaht(s disgord.Session, evt *disgord.MessageCreate) {
	channelIDs := []disgord.Snowflake{}

	// TODO get voice channels of all mentioned users
	if len(evt.Message.Mentions) > 0 {
		for _, user := range evt.Message.Mentions {
			zb.Logger().Debug(fmt.Sprintf("mentions: %s\n", user.String()))
		}
	}

	// manually select channel IDs based on guild
	if evt.Message.GuildID == purdoobahGuildID {
		channelIDs = append(channelIDs, purdoobahChannelID)
	} else if evt.Message.GuildID == otherPurdoobahGuildID {
		channelIDs = append(channelIDs, otherPurdoobahChannelID)
	} else if evt.Message.GuildID == bangBrosGuildID {
		channelIDs = append(channelIDs, bangBrosChannelID)
	} else {
		zb.Logger().Debug(fmt.Sprintf("unknown Guild ID: %s\n", evt.Message.GuildID))
		return
	}

	if len(channelIDs) == 0 {
		// TODO if no mentions, default to channel ID of voice channel of original message author
	}

	for _, channelID := range channelIDs {
		go zb.voiceZaht(s, evt.Message.GuildID, channelID)
	}
}

func (zb *ZahtBot) voiceZaht(session disgord.Session, guildID, channelID disgord.Snowflake) {
	if zb.isVoiceChannelActive(channelID) {
		zb.Logger().Debug(fmt.Sprintf("already Zahting in channel %s, skipping\n", channelID))
		return
	}
	zb.setVoiceChannelActivity(channelID, true)

	zb.Logger().Info(fmt.Sprintf("Zahting...\tGuild: %s\tChannel: %v\n", guildID, channelID))

	voice, err := session.VoiceConnect(guildID, channelID)
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Voice Connect error: %+v\n", err))
		zb.setVoiceChannelActivity(channelID, false)
		return
	}

	err = voice.StartSpeaking()
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Start Speaking error: %+v\n", err))
		zb.setVoiceChannelActivity(channelID, false)
		return
	}

	err = voice.SendDCA(bytes.NewReader(zb.dca))
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Send DCA error: %+v\n", err))
		zb.setVoiceChannelActivity(channelID, false)
		return
	}

	err = voice.StopSpeaking()
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Stop Speaking error: %+v\n", err))
		zb.setVoiceChannelActivity(channelID, false)
		return
	}

	err = voice.Close()
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Voice Close error: %+v\n", err))
		zb.setVoiceChannelActivity(channelID, false)
		return
	}

	zb.setVoiceChannelActivity(channelID, false)
	zb.Logger().Info(fmt.Sprintf("Zahted!\tGuild: %s\tChannel: %v\n", guildID, channelID))
}
