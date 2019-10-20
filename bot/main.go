package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/nats-io/nats.go"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"os"
	"os/signal"
	"syscall"
)

var RethinkConnection *r.Session

func main() {
	RethinkHost := os.Getenv("RETHINK_HOST")
	if RethinkHost == "" {
		RethinkHost = "127.0.0.1:28015"
	}
	RethinkPass := os.Getenv("RETHINK_PASSWORD")
	RethinkUser := os.Getenv("RETHINK_USER")
	if RethinkUser == "" {
		RethinkUser = "admin"
	}
	s, err := r.Connect(r.ConnectOpts{
		Address: RethinkHost,
		Password: RethinkPass,
		Username: RethinkUser,
		Database: "slimwidgets",
	})
	if err != nil {
		panic(err)
	}
	RethinkConnection = s

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}

	DiscordInit(discord)

	nc, err := nats.Connect(os.Getenv("NATS_HOST"))
	if err != nil {
		panic(err)
	}
	NATSInit(discord, nc)

	err = discord.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot is running. Always listening.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	_ = discord.Close()
}
