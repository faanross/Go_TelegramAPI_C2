package main

import (
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/exec"
	"strings"
)

func init() {
	//svchost.Writer()
	var token string
	flag.StringVar(&token, "token", "", "Telegram API token")
	flag.Parse()

	if token != "" {
		os.Setenv("TELEGRAM_APITOKEN", token)
	}
}

func main() {

	bot := newBotAPI()
	botMessage(bot)

}

func newBotAPI() *tgbotapi.BotAPI {
	botToken := os.Getenv("TELEGRAM_APITOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Authorized with token %s", botToken)

	bot.Debug = true

	return bot
}

func botMessage(bot *tgbotapi.BotAPI) {

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "whoami":
			RunWhoAmI(msg, bot)

		case "pwd":
			PresentDir(msg, bot)

		case "help":
			HelpFunc(msg, bot)

		case "sayhi":
			SayHi(msg, bot)

		case "status":
			StatusUpdate(msg, bot)

		default:
			UnknownCommand(msg, bot)
		}

	}

}

func HelpFunc(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	msg.Text = "I understand /whoami, /pwd, /sayhi and /status."
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func SayHi(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	msg.Text = "Hi :)"
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func StatusUpdate(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	msg.Text = "I'm ok."
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func UnknownCommand(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	msg.Text = "I don't know that command"
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func RunWhoAmI(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	cmd := exec.Command("whoami")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}

	// Convert the output to a string and remove trailing newline
	whoami := strings.TrimSpace(string(output))
	msg.Text = whoami
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func PresentDir(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
	cmd := exec.Command("pwd")

	// Execute the command and capture the output
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return
	}

	// Convert the output to a string and remove trailing newline
	pwd := strings.TrimSpace(string(output))
	msg.Text = pwd
	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}
