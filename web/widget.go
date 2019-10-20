package main

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"text/template"
)

func WidgetJSON(ctx *fasthttp.RequestCtx) {
	GuildID := ctx.UserValue("guildId").(string)
	Guild := GetGuild(GuildID)
	if Guild == nil {
		ctx.Response.SetStatusCode(404)
	} else {
		ctx.Response.SetStatusCode(200)
	}
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	b, err := json.Marshal(&Guild)
	if err != nil {
		panic(err)
	}
	ctx.Response.SetBody(b)
}

func WidgetHTML(ctx *fasthttp.RequestCtx) {
	GuildID := ctx.UserValue("guildId").(string)
	Guild := GetGuild(GuildID)
	if Guild == nil {
		ctx.Response.SetStatusCode(404)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.SetBody([]byte("Guild not found."))
		return
	}
	ctx.Response.SetStatusCode(200)
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.New("widget").Parse(WidgetTemplate)
	if err != nil {
		panic(err)
	}
	bw := bytes.Buffer{}
	err = t.Execute(&bw, *Guild)
	ctx.Response.SetBody(bw.Bytes())
}
