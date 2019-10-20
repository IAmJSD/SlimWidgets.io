package main

import "github.com/bwmarrin/discordgo"

func DiscordInit(session *discordgo.Session) {
	session.AddHandler(OnReady)
	session.AddHandler(OnMessage)
	session.AddHandler(OnReactionAdd)
	session.AddHandler(OnGuildCreate)
	session.AddHandler(OnGuildDelete)
}
