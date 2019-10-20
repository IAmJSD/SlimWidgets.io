package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"text/template"
)

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
	GuildID := ctx.UserValue("guildId").(string)
	CaptchaResult := ctx.UserValue("captchaResult").(string)

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

	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
		"secret": {
			os.Getenv("RECAPTCHA_SECRET_KEY"),
		},
		"response": {
			CaptchaResult,
		},
	})
	if err != nil {
		panic(err)
	}
	var Response map[string]interface{}
	defer resp.Body.Close()
	RawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(RawBody, &Response)
	if err != nil {
		panic(err)
	}

	if !Response["success"].(bool) {
		ctx.Response.SetStatusCode(400)
		ctx.SetBody([]byte(fmt.Sprintf("%v", Response["error-codes"].([]string))))
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
