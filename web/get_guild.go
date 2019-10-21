package main

import (
	"encoding/json"
	"strconv"
	"time"
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

func GetGuild(GuildID string) *GuildInfo {
	_, err := strconv.ParseInt(GuildID, 10, 64)
	if err != nil {
		return nil
	}
	reply, err := NATSClient.Request("guild-get", []byte(GuildID), 5 * time.Second)
	if err != nil {
		return nil
	}
	var Guild *GuildInfo
	err = json.Unmarshal(reply.Data, &Guild)
	if err != nil {
		panic(err)
	}
	if Guild == nil {
		return nil
	}
	Guild.GuildName = Policy.Sanitize(Guild.GuildName)
	if Guild.Description != nil {
		d := Policy.Sanitize(*Guild.Description)
		Guild.Description = &d
	}
	return Guild
}
