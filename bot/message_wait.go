package main

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type MessageWaiter struct {
	ChannelID string
	UserID string
	Channel chan *discordgo.Message
	Done bool
}

func (w MessageWaiter) Result(msg *discordgo.Message) {
	if w.Done {
		return
	}
	w.Done = true
	w.Channel <- msg
	close(w.Channel)
	for v := range waiters {
		waiters = append(waiters[:v], waiters[v+1:]...)
	}
}

var waiters []*MessageWaiter

func WaitForMessage(ChannelID string, UserID string, Timeout int) *discordgo.Message {
	Channel := make(chan *discordgo.Message)
	waiter := MessageWaiter{
		ChannelID: ChannelID,
		UserID: UserID,
		Channel: Channel,
	}
	waiters = append(waiters, &waiter)

	if Timeout != 0 {
		go func() {
			time.Sleep(time.Minute * time.Duration(Timeout))
			if !waiter.Done {
				waiter.Result(nil)
			}
		}()
	}

	j, m := <-Channel
	if m {
		waiter.Done = true
		return j
	}
	return nil
}

func MessageWaitHandler(msg *discordgo.Message) {
	for _, v := range waiters {
		if v.UserID == msg.Author.ID && v.ChannelID == msg.ChannelID {
			v.Result(msg)
			return
		}
	}
}
