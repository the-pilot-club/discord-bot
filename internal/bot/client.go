package bot

import (
	"github.com/getsentry/sentry-go"
	"log"
	"os"
	"os/signal"
	"tpc-discord-bot/handlers"
	"tpc-discord-bot/internal/config"
	"tpc-discord-bot/util"

	"github.com/bwmarrin/discordgo"
)

func Session() (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	return discord, nil
}

func Run() {

	log.Print("Starting discord-bot-v3")
	session, err := Session()
	if err != nil {
		sentry.CaptureException(err)
		println(err.Error())
	}
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// add all handlers to the session
	AddHandlers(session)

	session.AddHandler(handlers.OnGuildMemberAdd)
	session.AddHandler(handlers.OnGuildMemberRemove)

	err = session.Open()
	if err != nil {
		sentry.CaptureException(err)
		println(err.Error())
	}

	util.HandleApplicationCommandUpdates(session)

	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}

// AddHandlers adds all the handlers to the session
func AddHandlers(s *discordgo.Session) {

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		go config.IntervalReloadConfigs()
		go handlers.HandleCLientReady(s)
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		go handlers.MessageCreateHandler(s, m)
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.GuildScheduledEventCreate) {
		go handlers.HandleScheduledEventCreate(s, m)
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.GuildScheduledEventUpdate) {
		go handlers.HandleScheduledEventUpdate(s, m)
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.GuildScheduledEventDelete) {
		go handlers.HandleScheduledEventDelete(s, m)
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		go handlers.InteractionCreateHandler(s, i)
	})

}
