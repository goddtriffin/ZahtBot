package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"

	"github.com/andersfylling/disgord"
	"github.com/sirupsen/logrus"
)

var hardcodedChannelID = disgord.ParseSnowflakeString("559468274031656963")

// ZahtBot is the Discord ZahtBot.
type ZahtBot struct {
	*disgord.Client

	dca []byte
}

// NewZahtBot creates a new ZahtBot.
func NewZahtBot(botToken string) (*ZahtBot, error) {
	// Zaht audio file
	dca, err := loadDCA()
	if err != nil {
		log.Printf("Load DCA error: %+v\n", err)
		return nil, err
	}

	zb := &ZahtBot{
		Client: disgord.New(disgord.Config{
			ProjectName: "ZahtBot",
			BotToken:    botToken,
			Logger: &logrus.Logger{
				Out:       os.Stderr,
				Formatter: new(logrus.JSONFormatter),
				Hooks:     make(logrus.LevelHooks),
				Level:     logrus.InfoLevel,
			},
		}),
		dca: dca,
	}

	zb.Ready(func() {
		log.Println("ZahtBot is online!")
	})

	zb.On(disgord.EvtMessageCreate, zb.mux)

	return zb, nil
}

func (zb *ZahtBot) mux(session disgord.Session, m *disgord.MessageCreate) {
	if m.Message.Content == "!zaht" {
		zb.zaht(session, m)
	}
}

func (zb *ZahtBot) zaht(session disgord.Session, m *disgord.MessageCreate) {
	log.Println("Zahting...")

	voice, err := session.VoiceConnect(m.Message.GuildID, hardcodedChannelID)
	if err != nil {
		log.Printf("Voice Connect error: %+v\n", err)
		return
	}
	defer voice.Close()

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

	log.Println("Zahted!")
}

func loadDCA() ([]byte, error) {
	buf, err := ioutil.ReadFile("assets/zaht.dca")
	if err != nil {
		return nil, err
	}

	return buf, nil
}
