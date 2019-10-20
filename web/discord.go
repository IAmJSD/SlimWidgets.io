package main

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"os"
	"time"
)

var Token = os.Getenv("DISCORD_TOKEN")

func MakeInvite(Channel string) *string {
	URL := "https://discordapp.com/api/v6/channels/" + Channel + "/invites"
	Client := http.Client{
		Timeout:       10 * time.Second,
	}
	Req, err := http.NewRequest(
		"POST", URL,
		bytes.NewBuffer([]byte(`{"max_age": 3600, "max_uses": 1, "unique": true}`)))
	if err != nil {
		panic(err)
	}
	Req.Header.Set("Content-Type", "application/json")
	Req.Header.Set("Authorization", "Bot " + Token)
	Response, err := Client.Do(Req)
	if err != nil {
		panic(err)
	}

	if math.Floor(float64(Response.StatusCode) / 100) == 4 {
		return nil
	}

	var ResponseJSON map[string]interface{}
	defer Response.Body.Close()
	var b []byte
	_, err = Response.Body.Read(b)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, &ResponseJSON)
	if err != nil {
		panic(err)
	}

	Invite := "https://discord.gg/" + ResponseJSON["code"].(string)
	return &Invite
}
