package bot

import (
	"log"
	"os"
	"os/signal"
	"tpc-discord-bot/handlers"
	"tpc-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

func Session() (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, err
	}
	return discord, nil
}

func Run() {
	log.Print("Starting discord-bot-v3")
	session, err := Session()
	if err != nil {
		println(err.Error())
	}
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// add all handlers to the session
	AddHandlers(session)

	err = session.Open()
	if err != nil {
		println(err.Error())
	}
	defer session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}

// AddHandlers adds all the handlers to the session
func AddHandlers(session *discordgo.Session) {

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		go config.IntervalReloadConfigs()
		go handlers.HandleCLientReady(s)
	})

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		go handlers.MessageCreateHandler(s, m)
	})

	session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildScheduledEventCreate) {
		go handlers.HandleScheduledEventCreate(s, m)
	})

	session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildScheduledEventUpdate) {
		go handlers.HandleScheduledEventUpdate(s, m)
	})

	session.AddHandler(func(s *discordgo.Session, m *discordgo.GuildScheduledEventDelete) {
		go handlers.HandleScheduledEventDelete(s, m)
	})

}
