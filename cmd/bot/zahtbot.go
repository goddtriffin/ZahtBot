package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/purdoobahs/ZahtBot/internal/cache"
	"github.com/purdoobahs/ZahtBot/internal/cache/memory"
	"github.com/purdoobahs/ZahtBot/internal/command"
	"github.com/sirupsen/logrus"
)

// ZahtBot is the Discord ZahtBot.
type ZahtBot struct {
	*disgord.Client

	voiceStateCache cache.VoiceState
	guildStateCache cache.GuildState

	thumbnailURL string
	commands     []command.Command

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

		voiceStateCache: memory.NewVoiceStateCache(),
		guildStateCache: memory.NewGuildStateCache(),

		thumbnailURL: "https://www.cla.purdue.edu/facultyStaff/profiles/new/newfaculty-17/full/Sweet_Jonathan.jpg",
		commands: []command.Command{
			{Name: "help", Description: "displays help"},
			{Name: "zaht", Description: "ZAHTZAHTZAHTZAHTZAHTZAHT"},
		},

		dca:            dca,
		activeChannels: map[disgord.Snowflake]interface{}{},
	}

	zb.Ready(func() {
		zb.Logger().Info("ZahtBot is online!")
		zb.startHealthEndpointServer()
	})

	// filters
	filter, _ := std.NewMsgFilter(context.Background(), zb)
	filter.SetPrefix("/")

	// /help
	zb.On(
		disgord.EvtMessageCreate,

		filter.NotByBot,
		filter.HasPrefix,
		filterNonDM,

		filterNonHelpCommands,
		zb.commandHelp,
	)

	// /zaht
	zb.On(
		disgord.EvtMessageCreate,

		filter.NotByBot,
		filter.HasPrefix,
		filterNonDM,

		filterNonZahtCommands,
		zb.commandZaht,
	)

	// Guild Create (when bot joins guild)
	zb.On(disgord.EvtGuildCreate, zb.guildCreate)

	// Voice State Update
	zb.On(disgord.EvtVoiceStateUpdate, zb.voiceStateUpdate)

	return zb, nil
}

func (zb *ZahtBot) reply(session disgord.Session, evt *disgord.MessageCreate, reply interface{}) {
	_, err := evt.Message.Reply(context.Background(), session, reply)
	if err != nil {
		zb.Logger().Error(fmt.Sprintf("reply error: %+v\n", err))
	}
}

// getVoiceChannelID retrieves the voice channel ID of the message poster, if they're in one
func (zb *ZahtBot) getVoiceChannelID(session disgord.Session, evt *disgord.MessageCreate) disgord.Snowflake {
	_, vs := zb.voiceStateCache.GetVoiceState(evt.Message.Author.ID)
	if vs == nil {
		zb.Logger().Info(fmt.Sprintf("%s (%s) is not in a voice channel\n", evt.Message.Author.Username, evt.Message.Author.ID))
		return 0
	}

	return vs.ChannelID
}

func (zb *ZahtBot) startHealthEndpointServer() {
	// create the server
	addr := ":8080"
	srv := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			zb.Logger().Info(fmt.Sprintf("%s %s", req.Method, req.URL.Path))

			if req.URL.Path != "/api/v1/health" {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("OK"))
			if err != nil {
				zb.Logger().Error(err)
				return
			}
		}),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// run the server
	zb.Logger().Info(fmt.Sprintf("Health checkpoint is being served at: %s/api/v1/health", addr))
	err := srv.ListenAndServe()
	if err != nil {
		// print error on exit
		zb.Logger().Error(err)
		os.Exit(1)
	}
}
