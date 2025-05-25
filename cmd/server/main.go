package main

import (
	"log/slog"
	mqproducer "service-healthz-checker/internal/MQ/producer"
	"service-healthz-checker/internal/command"
	"service-healthz-checker/internal/config"
	"service-healthz-checker/internal/errs"
	"service-healthz-checker/internal/logger"
	requestmodel "service-healthz-checker/internal/model/requestModel"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.MustLoad()
	slogger := logger.SetupLogger(cfg.Env)
	bot, updates := setupBot(cfg.BotToken, slogger)

	producer := mqproducer.New(cfg.MqConfig.Topics.FromServerTopic, cfg.MqConfig.Host, slogger)

	for update := range updates {
		if update.Message != nil {
			slogger.Debug("REQUEST:",
				slog.String("User: ", update.Message.From.UserName),
				slog.String("TEXT: ", update.Message.Text))

			reqModel := parseCommand(&update, bot, slogger)
			if reqModel == nil {
				continue
			}
			producer.WriteToTopic("", reqModel)
		}
	}
}

func setupBot(token string, log *slog.Logger) (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(token)
	errs.FailOnError(err, "failed connect to TG BOT")

	bot.Debug = true
	log.Debug("Authorized on account", slog.String("username", bot.Self.UserName))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	return bot, updates
}

func parseCommand(upt *tgbotapi.Update, bot *tgbotapi.BotAPI, slogger *slog.Logger) *requestmodel.RequestCommand {
	if upt.Message.Command() == command.LIST || upt.Message.Command() == command.START || upt.Message.Command() == command.HELP {
		return &requestmodel.RequestCommand{
			Command: upt.Message.Command(),
			Value:   "",
			ChatID:  upt.Message.Chat.ID,
		}
	}

	if upt.Message.Command() != command.ADD &&
		upt.Message.Command() != command.REMOVE &&
		upt.Message.Command() != command.GET &&
		upt.Message.Command() != command.CH &&
		upt.Message.Command() != command.HISTORY &&
		upt.Message.Command() != command.START &&
		upt.Message.Command() != command.HELP {
		sendTo("Неизвестная команда", upt, bot, slogger)
		return nil
	}

	spliVal := strings.Split(upt.Message.Text, " ")
	if len(spliVal) != 2 {
		sendTo("Введите параметр для команды", upt, bot, slogger)
		return nil
	}

	return &requestmodel.RequestCommand{
		Command: upt.Message.Command(),
		Value:   spliVal[1],
		ChatID:  upt.Message.Chat.ID,
	}
}

func sendTo(text string, update *tgbotapi.Update, bot *tgbotapi.BotAPI, slogger *slog.Logger) {
	slogger.Debug(
		"RESPONCE",
		slog.Int("SEND TO ChatID: ",
			int(update.Message.Chat.ID)),
		slog.String("TEXT", text),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	bot.Send(msg)
}
