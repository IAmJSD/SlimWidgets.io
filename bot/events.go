package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

func SetPresence(client *discordgo.Session) {
	err := client.UpdateStatus(0, "slimwidgets.io | Tag me to configure!")
	if err != nil {
		fmt.Println("Error setting game: ", err)
	}
}

func OnReady(client *discordgo.Session, _ *discordgo.Ready) {
	SetPresence(client)
}

func OnMessage(client *discordgo.Session, msg *discordgo.MessageCreate) {
	go func() {
		if msg.Author.Bot {
			return
		}

		channel, _ := client.Channel(msg.ChannelID)
		mentions := msg.Mentions
		if len(mentions) == 1 && mentions[0].ID == client.State.User.ID && channel.Type == discordgo.ChannelTypeGuildText {
			OpenMenu(client, msg)
		}
		MessageWaitHandler(msg.Message)
	}()
}

func OpenMenu(client *discordgo.Session, msg *discordgo.MessageCreate) {
	MenuID := uuid.Must(uuid.NewRandom()).String()
	m, err := client.ChannelMessageSendComplex(msg.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Loading...",
		},
	})
	if err != nil {
		return
	}
	Menu := CreateNewMenu(MenuID, *msg.Message, client)
	Menu.Display(msg.ChannelID, m.ID, client)
}

func OnReactionAdd(client *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	go func() {
		message, err := client.ChannelMessage(reaction.ChannelID, reaction.MessageID)
		if err != nil {
			return
		}
		user, err := client.User(reaction.UserID)
		if err != nil {
			return
		}
		if user.Bot {
			return
		}
		if message.Author.ID == client.State.User.ID && len(message.Embeds) == 1 && message.Embeds[0].Footer != nil {
			MenuID := message.Embeds[0].Footer.Text
			HandleMenuReactionEdit(client, reaction, MenuID)
		}
	}()
}
