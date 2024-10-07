package bot

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"tpc-discord-bot/handlers"
	"tpc-discord-bot/internal/config"
)

func Session() (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, err
	}
	return discord, nil
}

func Run() {
	log.Print("Starting discord-bot-v2")
	session, err := Session()
	if err != nil {
		println(err.Error())
	}
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		handlers.HandleCLientReady(s)
	})

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handlers.MessageCreateHandler(s, m)
	})

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
