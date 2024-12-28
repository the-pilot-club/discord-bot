package bot

import (
	"log"
	"os"
	"os/signal"
	"tpc-discord-bot/handlers"
	"tpc-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	Session *discordgo.Session
}

func NewClient(token string) *Client {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	return &Client{
		Session: session,
	}
}

func (c *Client) Start() {
	err := c.Session.Open()
	if err != nil {
		panic(err)
	}
	// Block forever
	select {}
}

func Run() {
	log.Print("Starting discord-bot-v3")
	client := NewClient(config.DiscordToken)

	client.Session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	client.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		go config.IntervalReloadConfigs()
		go handlers.HandleCLientReady(s)
	})

	client.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		go handlers.MessageCreateHandler(s, m)
		return
	})

	err := client.Session.Open()
	if err != nil {
		println(err.Error())
	}
	defer client.Session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}
