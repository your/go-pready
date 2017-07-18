package main

import (
	"time"

	slackbot "github.com/your/go-slackbot"
)

func runCron(bot *slackbot.Bot) {
	for {
		checkRepositories(bot)
		time.Sleep(checkInterval * time.Millisecond)
	}
}
