package bot

import (
	"log"
	"os"
	"os/signal"
	"tpc-discord-bot/commands"
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

	log.Println("Adding commands...")
	registeredGlobalCommands := make([]*discordgo.ApplicationCommand, len(commands.GlobalCommands))
	for i, v := range commands.GlobalCommands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredGlobalCommands[i] = cmd
	}
	registeredGuildCommands := make([]*discordgo.ApplicationCommand, len(commands.GuildCommands))
	for i, v := range commands.GuildCommands {
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, *commands.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredGuildCommands[i] = cmd
	}
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
