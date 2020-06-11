package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/sirupsen/logrus"
)

var purdoobahGuildID = disgord.ParseSnowflakeString("559468273255841794")
var purdoobahChannelID = disgord.ParseSnowflakeString("559468274031656963")

var otherPurdoobahGuildID = disgord.ParseSnowflakeString("715262479831138345")
var otherPurdoobahChannelID = disgord.ParseSnowflakeString("715262480254894142")

var bangBrosGuildID = disgord.ParseSnowflakeString("720415671514562632")
var bangBrosChannelID = disgord.ParseSnowflakeString("720415671514562637")

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

	// !zaht [optional: <mention> <mention> ...]
	zb.On(
		disgord.EvtMessageCreate,

		filter.NotByBot,
		filter.HasPrefix,
		filterNonDM,

		filterNonZahtCommands,
		zb.commandZaht,
	)

	return zb, nil
}

func (zb *ZahtBot) isVoiceChannelActive(channelID disgord.Snowflake) bool {
	_, ok := zb.activeChannels[channelID]
	return ok
}

func (zb *ZahtBot) setVoiceChannelActivity(channelID disgord.Snowflake, active bool) {
	if active {
		zb.activeChannels[channelID] = nil
	} else {
		delete(zb.activeChannels, channelID)
	}
}
