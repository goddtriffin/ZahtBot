package main

import (
	"log"
	"os"

	"github.com/andersfylling/disgord"
)

// ZahtBot is the Discord ZahtBot.
type ZahtBot struct {
	*disgord.Client
}

// NewZahtBot creates a new ZahtBot.
func NewZahtBot(botToken string) *ZahtBot {
	zb := &ZahtBot{
		Client: disgord.New(disgord.Config{
			ProjectName: "ZahtBot",
			BotToken:    botToken,
			Logger:      disgord.DefaultLogger(false),
		}),
	}

	zb.On(disgord.EvtMessageCreate, func(session disgord.Session, m *disgord.MessageCreate) {
		if m.Message.Content == "!zaht" {
			zb.zaht(session, m.Message.GuildID, m.Message.ChannelID)
		}
	})

	return zb
}

func (zb *ZahtBot) zaht(session disgord.Session, guildID, channelID disgord.Snowflake) {
	f, err := os.Open("assets/zaht.dca")
	if err != nil {
		log.Printf("File Open error: %+v\n", err)
		return
	}
	defer f.Close()

	voice, err := session.VoiceConnect(guildID, disgord.ParseSnowflakeString("559468274031656963"))
	if err != nil {
		log.Printf("Voice Connect error: %+v\n", err)
		return
	}
	defer zb.closeVoice(voice)

	// Sending a speaking signal is mandatory before sending voice data
	err = voice.StartSpeaking()
	if err != nil {
		log.Printf("Start Speaking error: %+v\n", err)
		return
	}

	err = voice.SendDCA(f)
	if err != nil {
		log.Printf("Send DCA error: %+v\n", err)
		return
	}

	err = voice.StopSpeaking()
	if err != nil {
		log.Printf("Stop Speaking error: %+v\n", err)
		return
	}
}

func (zb *ZahtBot) closeVoice(voice disgord.VoiceConnection) {
	err := voice.Close()
	if err != nil {
		log.Printf("Voice Close error\n")
		panic(err)
	}
}
