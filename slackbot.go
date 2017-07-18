package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	slackbot "github.com/your/go-slackbot"
	"golang.org/x/net/context"
)

const (
	withoutTyping = slackbot.WithoutTyping

	helpText = "I can do the following: :sparkles: \n" +
		"`probot hi` for a simple message.\n" +
		"`probot help` to see this again."
)

var greetingPrefixes = []string{"Hey", "Hi", "Hello"}
var repositories []string
var githubURLRegex, _ = regexp.Compile(`http(s)?:\/\/github.com\/[\w-]+\/[\w-]+(\/)?`)

func runBot() {
	bot := slackbot.New(os.Getenv("SLACK_TOKEN"))

	go runCron(bot, &repositories)

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()

	hi := "hi|hello|probot hi|probot hello"
	toMe.Hear(hi).MessageHandler(helloHandler)
	bot.Hear(hi).MessageHandler(helloHandler)
	bot.Hear("help|probot help").MessageHandler(helpHandler)
	bot.Hear(`^review watch (\S+)\s*(\w+)?`).MessageHandler(reviewWatchHandler)
	bot.Hear(`^review unwatch (\S+)\s*(\w+)?`).MessageHandler(reviewUnwatchHandler)
	bot.Hear("(probot ).*").MessageHandler(catchAllHandler)
	bot.Run()
}

func sendMessage(bot *slackbot.Bot, msg string, channelName string) {
	channels, err := bot.Client.GetChannels(true)
	if err != nil {
		log.Println("Slack API Error:", err)
		return
	}

	channelID := getChannelID(channelName, channels)
	if channelID == "" {
		log.Printf("Slack API Error: cannot fetch ID for channel #%s\n", channelName)
		return
	}

	bot.RTM.SendMessage(bot.RTM.NewOutgoingMessage(msg, channelID))
}

func getChannelID(channelName string, channels []slack.Channel) string {
	var channelID string
	for _, channel := range channels {
		if channel.Name == channelName {
			channelID = channel.ID
			break
		}
	}
	return channelID
}

func helloHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	msg := greetingPrefixes[rand.Intn(len(greetingPrefixes))] + " <@" + evt.User + ">!"
	bot.Reply(evt, msg, withoutTyping)

	if slackbot.IsDirectMessage(evt) {
		dmMsg := "Hey, what's up? Can I `help`? :confused:"
		bot.Reply(evt, dmMsg, withoutTyping)
	}
}

func reviewWatchHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	var inputMsg = evt.Msg.Text
	var outputMsg string

	if stringMatchesRegex(inputMsg, githubURLRegex) {
		repo := getRepoNameFromMsg(inputMsg, githubURLRegex)

		if isStringInSlice(repo, repositories) {
			outputMsg = fmt.Sprintf(":neutral_face: Repository `%s` is already in the watchlist.", repo)
		} else {
			repositories = append(repositories, repo)

			outputMsg = fmt.Sprintf(":wine_glass: Added repository `%s` to watchlist.", repo)
		}
	} else {
		outputMsg = ":warning: Sorry, the URL is invalid, please double check it."
	}

	bot.Reply(evt, outputMsg, withoutTyping)
}

func reviewUnwatchHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	var inputMsg = evt.Msg.Text
	var outputMsg string

	if stringMatchesRegex(inputMsg, githubURLRegex) {
		repo := getRepoNameFromMsg(inputMsg, githubURLRegex)

		if isStringInSlice(repo, repositories) {
			removeStringFromSlice(repo, &repositories)

			outputMsg = fmt.Sprintf(":wine_glass: Removed repository `%s` from watchlist.", repo)
		} else {
			outputMsg = fmt.Sprintf(":neutral_face: Repository `%s` is not in the watchlist.", repo)
		}
	} else {
		outputMsg = ":warning: Sorry, the URL is invalid, please double check it."
	}

	bot.Reply(evt, outputMsg, withoutTyping)
}

func catchAllHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	msg := fmt.Sprintf("Hey, I don't know how to help you with that. :confused:\n\n%s", helpText)
	bot.Reply(evt, msg, withoutTyping)
}

func helpHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, helpText, withoutTyping)
}

// TODO: find a better place.
func getRepoNameFromMsg(str string, r *regexp.Regexp) string {
	// take 4th from ["http:", "", github.com", "your", "_repo"]
	return strings.Split(r.FindAllString(str, 1)[0], "/")[4]
}
