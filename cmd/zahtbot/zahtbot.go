package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/andersfylling/disgord"
)

// ZahtBot is the Discord ZahtBot.
type ZahtBot struct {
	*disgord.Client

	dca []byte
}

// NewZahtBot creates a new ZahtBot.
func NewZahtBot(botToken string) *ZahtBot {
	dca, err := loadDCA()
	if err != nil {
		log.Printf("Load DCA error: %+v\n", err)
		panic(err)
	}

	zb := &ZahtBot{
		Client: disgord.New(disgord.Config{
			ProjectName: "ZahtBot",
			BotToken:    botToken,
			Logger:      disgord.DefaultLogger(false),
		}),
		dca: dca,
	}

	zb.On(disgord.EvtMessageCreate, func(session disgord.Session, m *disgord.MessageCreate) {
		if m.Message.Content == "!zaht" {
			zb.zaht(session, m.Message.GuildID, m.Message.ChannelID)
		}
	})

	return zb
}

func (zb *ZahtBot) zaht(session disgord.Session, guildID, channelID disgord.Snowflake) {
	voice, err := session.VoiceConnect(guildID, disgord.ParseSnowflakeString("559468274031656963"))
	if err != nil {
		log.Printf("Voice Connect error: %+v\n", err)
		return
	}

	err = voice.StartSpeaking()
	if err != nil {
		log.Printf("Start Speaking error: %+v\n", err)
		return
	}

	err = voice.SendDCA(bytes.NewReader(zb.dca))
	if err != nil {
		log.Printf("Send DCA error: %+v\n", err)
		return
	}

	err = voice.StopSpeaking()
	if err != nil {
		log.Printf("Stop Speaking error: %+v\n", err)
		return
	}

	err = voice.Close()
	if err != nil {
		log.Printf("Voice Close error\n")
		return
	}
}

func loadDCA() ([]byte, error) {
	buf, err := ioutil.ReadFile("assets/zaht.dca")
	if err != nil {
		return nil, err
	}

	return buf, nil
}
