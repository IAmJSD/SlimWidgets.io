package main

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"gopkg.in/ezzarghili/recaptcha-go.v3"
	"os"
	"text/template"
	"time"
)

var CAPTCHA, _ = recaptcha.NewReCAPTCHA(os.Getenv("RECAPTCHA_SECRET_KEY"), recaptcha.V2, 10 * time.Second)

func InviteHTML(ctx *fasthttp.RequestCtx) {
	GuildID := ctx.UserValue("guildId").(string)
	Guild := GetGuild(GuildID)
	if Guild == nil {
		ctx.Response.SetStatusCode(404)
		ctx.SetBody([]byte("Guild not found."))
		return
	}
	if Guild.InviteURL == nil {
		ctx.Response.SetStatusCode(400)
		ctx.SetBody([]byte("This guild has invites off."))
		return
	}
	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.New("invite").Parse(InviteTemplate)
	if err != nil {
		panic(err)
	}
	bw := bytes.Buffer{}
	err = t.Execute(&bw, map[string]interface{}{
		"GuildName": Guild.GuildName,
		"SiteKey": os.Getenv("RECAPTCHA_SITE_KEY"),
	})
	ctx.Response.SetBody(bw.Bytes())
}

func InviteCaptchaHandler(ctx *fasthttp.RequestCtx) {
	CaptchaResult := ctx.UserValue("captchaResult").(string)
	GuildID := ctx.UserValue("guildId").(string)

	Guild := GetGuild(GuildID)
	if Guild == nil {
		ctx.Response.SetStatusCode(404)
		ctx.SetBody([]byte("Guild not found."))
		return
	}
	if Guild.InviteURL == nil {
		ctx.Response.SetStatusCode(400)
		ctx.SetBody([]byte("This guild has invites off."))
		return
	}

	err := CAPTCHA.Verify(CaptchaResult)
	if err != nil {
		ctx.Response.SetStatusCode(400)
		ctx.SetBody([]byte(err.Error()))
		return
	}

	invite := MakeInvite(Guild.DefaultChannel)
	if invite == nil {
		ctx.Response.SetStatusCode(400)
		ctx.SetBody([]byte("Failed to create the invite. Does the bot have permission?"))
		return
	}

	ctx.Redirect(*invite, 301)
}
