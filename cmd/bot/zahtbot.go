package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/sirupsen/logrus"
)

var purdoobahGuildID = disgord.ParseSnowflakeString("559468273255841794")
var purdoobahChannelID = disgord.ParseSnowflakeString("707682640064413777")

var otherPurdoobahGuildID = disgord.ParseSnowflakeString("715262479831138345")
var otherPurdoobahChannelID = disgord.ParseSnowflakeString("715620876061507615")

// ZahtBot is the Discord ZahtBot.
type ZahtBot struct {
	*disgord.Client

	dca            []byte
	activeChannels map[disgord.Snowflake]interface{}
}

// NewZahtBot creates a new ZahtBot.
func NewZahtBot(botToken string) (*ZahtBot, error) {
	logger := &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	// Zaht audio file
	dca, err := loadDCA()
	if err != nil {
		logger.Debug(fmt.Sprintf("Load DCA error: %+v\n", err))
		return nil, err
	}

	zb := &ZahtBot{
		Client: disgord.New(disgord.Config{
			ProjectName: "ZahtBot",
			BotToken:    botToken,
			Logger:      logger,
		}),
		dca:            dca,
		activeChannels: map[disgord.Snowflake]interface{}{},
	}

	zb.Ready(func() {
		zb.Logger().Info("ZahtBot is online!")
	})

	// filters
	filter, _ := std.NewMsgFilter(context.Background(), zb)
	filter.SetPrefix("!")

	zb.On(
		disgord.EvtMessageCreate,

		filter.NotByBot,
		filter.HasPrefix,
		filterNonDM,

		filterNonZahtCommands,
		zb.zahtCommand,
	)

	return zb, nil
}

func (zb *ZahtBot) zahtCommand(session disgord.Session, m *disgord.MessageCreate) {
	channelIDs := []disgord.Snowflake{}

	// TODO get voice channels of all mentioned users
	if len(m.Message.Mentions) > 0 {
		for _, user := range m.Message.Mentions {
			zb.Logger().Debug(fmt.Sprintf("mentions: %s\n", user.String()))
		}
	}

	// manually select channel IDs based on guild
	if m.Message.GuildID == purdoobahGuildID {
		channelIDs = append(channelIDs, purdoobahChannelID)
	} else if m.Message.GuildID == otherPurdoobahGuildID {
		channelIDs = append(channelIDs, otherPurdoobahChannelID)
	} else {
		zb.Logger().Debug(fmt.Sprintf("unknown Guild ID: %s\n", m.Message.GuildID))
		return
	}

	if len(channelIDs) == 0 {
		// TODO if no mentions, default to channel ID of voice channel of original message author
	}

	for _, channelID := range channelIDs {
		go zb.zaht(session, m.Message.GuildID, channelID)
	}
}

func (zb *ZahtBot) zaht(session disgord.Session, guildID, channelID disgord.Snowflake) {
	if zb.isChannelActive(channelID) {
		zb.Logger().Debug("already Zahting in channel, skipping")
		return
	}
	zb.setChannelActivity(channelID, true)

	zb.Logger().Info(fmt.Sprintf("Zahting...\tGuild: %s\tChannel: %v\n", guildID, channelID))

	voice, err := session.VoiceConnect(guildID, channelID)
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Voice Connect error: %+v\n", err))
		zb.setChannelActivity(channelID, false)
		return
	}

	err = voice.StartSpeaking()
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Start Speaking error: %+v\n", err))
		zb.setChannelActivity(channelID, false)
		return
	}

	err = voice.SendDCA(bytes.NewReader(zb.dca))
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Send DCA error: %+v\n", err))
		zb.setChannelActivity(channelID, false)
		return
	}

	err = voice.StopSpeaking()
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Stop Speaking error: %+v\n", err))
		zb.setChannelActivity(channelID, false)
		return
	}

	err = voice.Close()
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("Voice Close error: %+v\n", err))
		zb.setChannelActivity(channelID, false)
		return
	}

	zb.setChannelActivity(channelID, false)
	zb.Logger().Info(fmt.Sprintf("Zahted!\tGuild: %s\tChannel: %v\n", guildID, channelID))
}

func (zb *ZahtBot) isChannelActive(channelID disgord.Snowflake) bool {
	_, ok := zb.activeChannels[channelID]
	return ok
}

func (zb *ZahtBot) setChannelActivity(channelID disgord.Snowflake, active bool) {
	if active {
		zb.activeChannels[channelID] = nil
	} else {
		delete(zb.activeChannels, channelID)
	}
}
