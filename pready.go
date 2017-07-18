package main

import (
	"time"

	slackbot "github.com/your/go-slackbot"
)

func runCron(bot *slackbot.Bot, repositories *[]string) {
	for {
		checkRepositories(bot, repositories)
		time.Sleep(checkInterval * time.Millisecond)
	}
}
