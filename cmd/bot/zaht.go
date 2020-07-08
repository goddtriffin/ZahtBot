package main

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

func (zb *ZahtBot) commandZaht(session disgord.Session, evt *disgord.MessageCreate) {
	channelID := zb.getVoiceChannelID(session, evt)
	if channelID == 0 {
		return
	}

	go zb.voice(session, evt.Message.GuildID, channelID, zb.dca, "Zaht")

	zb.Logger().Info(fmt.Sprintf("%s (%s) called %s\n", evt.Message.Author.Username, evt.Message.Author.ID, evt.Message.Content))
}
