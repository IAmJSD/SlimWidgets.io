package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var MenuCache map[string]*EmbedMenu

type MenuInfo struct {
	MenuID string
	Author string
	Info []string
}

type MenuButton struct {
	Emoji string
	Name string
	Description string
}

type MenuReaction struct {
	button MenuButton
	function func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session)
}

type MenuReactions struct {
	ReactionSlice []MenuReaction
}

func (mr *MenuReactions) Add(reaction MenuReaction) {
	Slice := append(mr.ReactionSlice, reaction)
	mr.ReactionSlice = Slice
}

type EmbedMenu struct {
	Reactions *MenuReactions
	parent *EmbedMenu
	Embed *discordgo.MessageEmbed
	MenuInfo *MenuInfo
}

func (emm EmbedMenu) Display(ChannelID string, MessageID string, client *discordgo.Session) *error {
	MenuCache[emm.MenuInfo.MenuID] = &emm

	EmbedCopy := emm.Embed
	EmbedCopy.Footer = &discordgo.MessageEmbedFooter{
		Text: emm.MenuInfo.MenuID,
	}
	Fields := make([]*discordgo.MessageEmbedField, 0)
	for _, k := range emm.Reactions.ReactionSlice {
		Fields = append(Fields, &discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%s %s", k.button.Emoji, k.button.Name),
			Value: k.button.Description,
			Inline: false,
		})
	}
	EmbedCopy.Fields = Fields

	_, err := client.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Embed: EmbedCopy,
		ID: MessageID,
		Channel: ChannelID,
	})
	if err != nil {
		return &err
	}
	for _, k := range emm.Reactions.ReactionSlice {
		err := client.MessageReactionAdd(ChannelID, MessageID, k.button.Emoji)
		if err != nil {
			return &err
		}
	}
	return nil
}

func (emm EmbedMenu) NewChildMenu(embed discordgo.MessageEmbed, item MenuButton) *EmbedMenu {
	NewEmbedMenu := NewEmbedMenu(embed, emm.MenuInfo)
	NewEmbedMenu.parent = &emm
	Reaction := MenuReaction{
		button:   item,
		function: func(ChannelID string, MessageID string, _ *EmbedMenu, client *discordgo.Session) {
			_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)
			NewEmbedMenu.Display(ChannelID, MessageID, client)
		},
	}
	emm.Reactions.Add(Reaction)
	return &NewEmbedMenu
}

func (emm EmbedMenu) AddBackButton() {
	Reaction := MenuReaction{
		button:   MenuButton{
			Description: "Goes back to the parent menu.",
			Name: "Back",
			Emoji: "⬆",
		},
		function: func(ChannelID string, MessageID string, _ *EmbedMenu, client *discordgo.Session) {
			_ = client.MessageReactionsRemoveAll(ChannelID, MessageID)
			emm.parent.Display(ChannelID, MessageID, client)
		},
	}
	emm.Reactions.Add(Reaction)
}

func NewEmbedMenu(embed discordgo.MessageEmbed, info *MenuInfo) EmbedMenu {
	var reactions []MenuReaction
	menu := EmbedMenu{
		Reactions: &MenuReactions{
			ReactionSlice: reactions,
		},
		Embed: &embed,
		MenuInfo: info,
	}
	return menu
}

func HandleMenuReactionEdit(client *discordgo.Session, reaction *discordgo.MessageReactionAdd, MenuID string) {
	_ = client.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
	menu := MenuCache[MenuID]
	if menu == nil {
		return
	}

	if menu.MenuInfo.Author != reaction.UserID {
		return
	}

	for _, v := range menu.Reactions.ReactionSlice {
		if v.button.Emoji == reaction.Emoji.Name {
			v.function(reaction.ChannelID, reaction.MessageID, menu, client)
			return
		}
	}
}
