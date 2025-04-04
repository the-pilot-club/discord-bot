package handlers

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"tpc-discord-bot/internal/fcp"
)

/*
s *discordgo.Session isn't being used currently,
but added for future expansion to access disc api
*/
func OnGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	go func() {
		err := fcp.AddUser(m.User.ID, m.GuildID)
		if err != nil {
			log.Printf(err.Error())
		}
	}()
	log.Printf("added user %s to FCP", m.User.ID)
}

func OnGuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	go func() {
		err := fcp.RemoveUser(m.User.ID, m.GuildID)
		if err != nil {
			log.Printf(err.Error())
		}
	}()
	log.Printf("removed user %s from FCP", m.User.ID)
}
