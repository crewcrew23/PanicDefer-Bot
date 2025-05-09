package main

import (
	"log"
	"service-healthz-checker/internal/command"
	"service-healthz-checker/internal/config"
	"service-healthz-checker/internal/errs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.MustLoad()

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	errs.FailOnError(err, "failed connect to TG BOT")

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Command() {
			case command.ADD:
				sendTo("add command", update, bot)
			case command.LIST:
				sendTo("list command", update, bot)
			case command.REMOVE:
				sendTo("remove command", update, bot)
			default:
				sendTo("unknown command", update, bot)
			}
		}
	}
}

func sendTo(text string, update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}
