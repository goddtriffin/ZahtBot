package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/andersfylling/disgord"
)

func main() {
	token := flag.String("token", "", "Discord Bot Token")
	guildIDStr := flag.String("guildID", "", "Discord Guild ID")
	channelIDStr := flag.String("channelID", "", "Discord Channel ID")
	flag.Parse()

	if *token == "" || *guildIDStr == "" || *channelIDStr == "" {
		flag.Usage()
		os.Exit(1)
	}

	discord := disgord.New(disgord.Config{
		BotToken: *token,
		Logger:   disgord.DefaultLogger(true),
	})

	var voice disgord.VoiceConnection

	discord.Ready(func() {
		guildID := disgord.ParseSnowflakeString(*guildIDStr)
		channelID := disgord.ParseSnowflakeString(*channelIDStr)

		var err error
		voice, err = discord.VoiceConnect(guildID, channelID)
		if err != nil {
			panic(err)
		}
	})

	discord.On(disgord.EvtMessageCreate, func(_ disgord.Session, m *disgord.MessageCreate) {
		if m.Message.Content == "!zaht" {
			f, err := os.Open("zaht.wav")
			if err != nil {
				log.Println(err)
			}
			defer f.Close()

			// Sending a speaking signal is mandatory before sending voice data
			err = voice.StartSpeaking()
			if err != nil {
				log.Println(err)
			}

			err = voice.SendDCA(f)
			if err != nil {
				log.Println(err)
			}

			err = voice.StopSpeaking()
			if err != nil {
				log.Println(err)
			}
		}
	})

	err := discord.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	err = discord.DisconnectOnInterrupt()
	if err != nil {
		panic(err)
	}
}
