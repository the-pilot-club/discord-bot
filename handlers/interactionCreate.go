package handlers

import (
	"github.com/bwmarrin/discordgo"
	eventresponses "tpc-discord-bot/event-responses"
)

func InteractionCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionModalSubmit {
		data := i.ModalSubmitData()
		if data.CustomID == "TrainingRequest" {
			eventresponses.TrainingRequestModal(s, i)
		}
	}
	if i.Type == discordgo.InteractionMessageComponent {
		eventresponses.HandleButtonSubmit(s, i)
	}
	if i.Type == discordgo.InteractionApplicationCommand {
		eventresponses.GuildCommands(s, i)
		eventresponses.GlobalCommands(s, i)
	}
}
