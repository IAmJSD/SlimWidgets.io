package main

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func HTTPInit(router *fasthttprouter.Router) {
	router.GET("/widget/:guildId/json", WidgetJSON)
	router.GET("/widget/:guildId", WidgetHTML)
	router.GET("/invite/:guildId", InviteHTML)
	router.GET("/invite/:guildId/:captchaResult", InviteCaptchaHandler)
	router.GET("/", func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Content-Type", "text/html")
		ctx.Response.SetStatusCode(200)
		ctx.Response.SetBody(IndexHTML)
	})
}
