package main

import (
	"strings"

	"github.com/andersfylling/disgord"
)

func filterNonDM(evt interface{}) interface{} {
	m := evt.(*disgord.MessageCreate)
	if m.Message.IsDirectMessage() {
		return nil
	}
	return evt
}

func filterNonHelpCommands(evt interface{}) interface{} {
	m := evt.(*disgord.MessageCreate)
	if strings.ToLower(m.Message.Content) != "/help" {
		return nil
	}
	return evt
}

func filterNonZahtCommands(evt interface{}) interface{} {
	m := evt.(*disgord.MessageCreate)
	if !strings.HasPrefix(m.Message.Content, "/zaht") {
		return nil
	}
	return evt
}
