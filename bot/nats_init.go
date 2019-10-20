package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/nats-io/nats.go"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strconv"
	"strings"
)

type GuildCountInfo struct {
	Online int `json:"online"`
	Offline int `json:"offline"`
	InVoice int `json:"in_voice"`
}

type GuildInfo struct {
	Counts *GuildCountInfo `json:"counts"`
	InviteURL *string `json:"invite_url"`
	Description *string `json:"description"`
	GuildName string `json:"guild_name"`
	GuildIcon string `json:"guild_icon"`
	DefaultChannel string `json:"default_channel"`
}

type GuildDB struct {
	Id string `gorethink:"id,omitempty"`
	Description *string `gorethink:"description"`
	Invites bool `gorethink:"invites"`
}

func NATSInit(session *discordgo.Session, nc *nats.Conn) {
	_, err := nc.Subscribe("guild-get", func(m *nats.Msg) {
		GuildStr := string(m.Data)
		GuildID, err := strconv.ParseInt(GuildStr, 10, 64)
		if err != nil {
			panic(err)
		}

		if int((GuildID >> 22) % int64(session.ShardCount)) != session.ShardID {
			// This is not related to this shard. Return here.
			return
		}

		DiscordGuild, err := session.Guild(GuildStr)
		if err != nil || DiscordGuild == nil {
			err = m.Respond([]byte("null"))
			if err != nil {
				panic(err)
			}
			return
		}

		var Guild GuildDB
		cursor, err := r.Table("guilds").Get(GuildStr).Run(RethinkConnection)
		if err != nil {
			panic(err)
		}
		err = cursor.One(&Guild)
		if err != nil {
			panic(err)
		}
		Online := 0
		Offline := 0
		for _, v := range DiscordGuild.Presences {
			if v.Status == "offline" {
				Offline++
			} else {
				Online++
			}
		}
		DefaultChannel := DiscordGuild.Channels[0].ID
		for _, v := range DiscordGuild.Channels {
			if strings.Contains(strings.ToLower(v.Name), "rules") || strings.Contains(strings.ToLower(v.Name), "general") {
				DefaultChannel = v.ID
				break
			}
		}
		Info := GuildInfo{
			Counts: &GuildCountInfo{
				Online:  Online,
				Offline: Offline,
				InVoice: len(DiscordGuild.VoiceStates),
			},
			InviteURL:   nil,
			Description: Guild.Description,
			GuildName: DiscordGuild.Name,
			GuildIcon: DiscordGuild.IconURL(),
			DefaultChannel: DefaultChannel,
		}
		if Guild.Invites {
			URL := "https://slimwidgets.io/invite/" + GuildStr
			Info.InviteURL = &URL
		}

		b, err := json.Marshal(&Info)
		if err != nil {
			panic(err)
		}
		err = m.Respond(b)
		if err != nil {
			panic(err)
		}
	})
	if err != nil {
		panic(err)
	}
}
