package bot

import (
	"log"
	"os"
	"os/signal"

	"tpc-discord-bot/handlers"
	"tpc-discord-bot/internal/bootstrap"
	"tpc-discord-bot/internal/config"
	"tpc-discord-bot/util"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	log.Print("Starting discord-bot-v3")
	session, cleanup, err := bootstrap.InitForBot(AddHandlers)
	if err != nil {
		log.Print(err.Error())
		return
	}
	defer cleanup()

	util.HandleApplicationCommandUpdates(session)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
}

// AddHandlers adds all the handlers to the session:.
func AddHandlers(s *discordgo.Session) {

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		go config.IntervalReloadConfigs()
		go handlers.HandleClientReady(s)
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
		go handlers.HandleGuildMemberUpdate(s, m)
	})
	s.AddHandler(func(s *discordgo.Session, m *discordgo.GuildBanAdd) {
		go handlers.HandleGuildBanAdd(s, m)
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.GuildBanRemove) {
		go handlers.HandleGuildBanRemove(s, m)
	})

	s.AddHandler(func(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
		go handlers.OnGuildMemberAdd(s, g)
	})
	s.AddHandler(func(s *discordgo.Session, g *discordgo.GuildMemberRemove) {
		go handlers.OnGuildMemberRemove(s, g)
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

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		go handlers.HandleQuizButton(s, i)
	})
}
