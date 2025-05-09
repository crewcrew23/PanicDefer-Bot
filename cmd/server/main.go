package main

import (
	"log/slog"
	"service-healthz-checker/internal/command"
	"service-healthz-checker/internal/config"
	"service-healthz-checker/internal/errs"
	"service-healthz-checker/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.MustLoad()
	slogger := logger.SetupLogger(cfg.Env)

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	errs.FailOnError(err, "failed connect to TG BOT")

	bot.Debug = true
	slogger.Debug("Authorized on account", slog.String("username", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			slogger.Debug("REQUEST:",
				slog.String("User: ", update.Message.From.UserName),
				slog.String("TEXT: ", update.Message.Text))

			switch update.Message.Command() {
			case command.ADD:
				sendTo("add command", update, bot, slogger)
			case command.LIST:
				sendTo("list command", update, bot, slogger)
			case command.REMOVE:
				sendTo("remove command", update, bot, slogger)
			default:
				sendTo("unknown command", update, bot, slogger)
			}
		}
	}
}

func sendTo(text string, update tgbotapi.Update, bot *tgbotapi.BotAPI, slogger *slog.Logger) {
	slogger.Info(
		"RESPONCE",
		slog.Int("SEND TO ChatID: ",
			int(update.Message.Chat.ID)),
		slog.String("TEXT", text),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}
