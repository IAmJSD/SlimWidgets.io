package main

import "github.com/buaazp/fasthttprouter"

func HTTPInit(router *fasthttprouter.Router) {
	router.GET("/widget/:guildId/json", WidgetJSON)
	router.GET("/widget/:guildId", WidgetHTML)
	router.GET("/invite/:guildId", InviteHTML)
	router.GET("/invite/:guildId/:captchaResult", InviteCaptchaHandler)
}
