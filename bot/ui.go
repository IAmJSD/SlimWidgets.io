package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

func CreateNewMenu(MenuID string, msg discordgo.Message, client *discordgo.Session) *EmbedMenu {
	var Guild GuildDB
	cursor, err := r.Table("guilds").Get(msg.GuildID).Run(RethinkConnection)
	if err != nil {
		panic(err)
	}
	err = cursor.One(&Guild)
	if err != nil {
		panic(err)
	}

	ToggleInvites := func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
		Guild.Invites = !Guild.Invites
		_, err := r.Table("guilds").Get(Guild.Id).Update(map[string]interface{}{
			"invites": Guild.Invites,
		}).Run(RethinkConnection)
		if err != nil {
			panic(err)
		}
		menu.Display(ChannelID, MessageID, client)
	}

	SetDescription := func(ChannelID string, MessageID string, menu *EmbedMenu, client *discordgo.Session) {
		defer menu.Display(ChannelID, MessageID, client)

		embed := &discordgo.MessageEmbed{
			Title: "Waiting for your description...",
			Description: "Please enter your new description.",
		}
		_, err := client.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID: MessageID,
			Channel: ChannelID,
			Embed: embed,
		})
		if err != nil {
			return
		}

		UserMessage := WaitForMessage(ChannelID, menu.MenuInfo.Author, 5)
		if UserMessage == nil {
			return
		}

		_, err = r.Table("guilds").Get(Guild.Id).Update(map[string]interface{}{
			"description": UserMessage.Content,
		}).Run(RethinkConnection)
		if err != nil {
			panic(err)
		}
	}

	MainMenu := NewEmbedMenu(
		discordgo.MessageEmbed{
			Title: "SlimWidgets Manager",
			Description: "Using this bot, you can configure your SlimWidget instance.",
			Color: 255,
		}, &MenuInfo{
			MenuID: MenuID,
			Author: msg.Author.ID,
			Info: []string{},
		},
	)
	perms, err := client.UserChannelPermissions(msg.Author.ID, msg.ChannelID)
	if err != nil {
		panic(err)
	}
	if (perms & 0x00000008) != 0x00000008 {
		MainMenu.Embed.Description += " You require administrator to configure this bot."
		return &MainMenu
	}

	var EnableDisableInvites string
	if Guild.Invites {
		EnableDisableInvites = "Enable Invites"
	} else {
		EnableDisableInvites = "Disable Invites"
	}
	MainMenu.Reactions.Add(MenuReaction{
		button:   MenuButton{
			Emoji: "📬",
			Name: EnableDisableInvites,
			Description: "This will toggle invites. Note that to invite users, the bot needs permission to make invites.",
		},
		function: ToggleInvites,
	})

	MainMenu.Reactions.Add(MenuReaction{
		button:   MenuButton{
			Emoji: "💬",
			Name: "Set Description",
			Description: "This will set the description that will show on the widget.",
		},
		function: SetDescription,
	})

	c := MainMenu.NewChildMenu(discordgo.MessageEmbed{
		Title: "iFrame Code",
		Description: fmt.Sprintf("To embed this iFrame, simply add the following code: ```html\n" + `<iframe src="https:/slimwidgets.io/widget/%s" width="350" height="400" allowtransparency="true" frameborder="0"></iframe>` + "\n```", msg.GuildID),
		Color: 255,
	}, MenuButton{
		Emoji:       "🔗",
		Name:        "Get iFrame code",
		Description: "Gives you the iFrame code so you can embed the widget into your site.",
	})
	c.AddBackButton()

	return &MainMenu
}
