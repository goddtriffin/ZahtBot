package main

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

func (zb *ZahtBot) commandHelp(session disgord.Session, evt *disgord.MessageCreate) {
	zb.reply(session, evt, &disgord.Embed{
		Description: "**ZahtBot Help**",
		Color:       15844367,
		Fields: []*disgord.EmbedField{
			{Name: "`!commands`", Value: "displays commands"},
		},
		Thumbnail: &disgord.EmbedThumbnail{URL: zb.thumbnailURL},
	})

	zb.Logger().Info(fmt.Sprintf("%s (%s) called !help\n", evt.Message.Author.Username, evt.Message.Author.ID))
}

func (zb *ZahtBot) commandCommands(session disgord.Session, evt *disgord.MessageCreate) {
	// convert list of commands to list of Disgord fields
	fields := []*disgord.EmbedField{}
	for _, command := range zb.commands {
		fields = append(fields, &disgord.EmbedField{
			Name:   fmt.Sprintf("`%s`", command.String()),
			Value:  command.Description,
			Inline: true,
		})
	}

	zb.reply(session, evt, &disgord.Embed{
		Description: "**ZahtBot Commands**",
		Color:       15844367,
		Fields:      fields,
		Thumbnail:   &disgord.EmbedThumbnail{URL: zb.thumbnailURL},
	})

	zb.Logger().Info(fmt.Sprintf("%s (%s) called !commands\n", evt.Message.Author.Username, evt.Message.Author.ID))
}

func (zb *ZahtBot) commandZaht(session disgord.Session, evt *disgord.MessageCreate) {
	channelID := zb.getVoiceChannelID(session, evt)
	if channelID == 0 {
		return
	}

	go zb.voice(session, evt.Message.GuildID, channelID, zb.dca, "Zaht")

	zb.Logger().Info(fmt.Sprintf("%s (%s) called %s\n", evt.Message.Author.Username, evt.Message.Author.ID, evt.Message.Content))
}
