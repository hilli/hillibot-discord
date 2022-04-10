package main

import (
	"github.com/hilli/hillibot-discord/bot"
	"github.com/hilli/hillibot-discord/config"
)

func main() {
	config.Load()
	bot.Start()
	<-make(chan struct{})
	return
}
