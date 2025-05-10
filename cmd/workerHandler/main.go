package main

import (
	"encoding/json"
	"log/slog"
	mqconsumer "service-healthz-checker/internal/MQ/consumer"
	"service-healthz-checker/internal/config"
	"service-healthz-checker/internal/errs"
	"service-healthz-checker/internal/logger"
	requestmodel "service-healthz-checker/internal/model/requestModel"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.MustLoad()
	slogger := logger.SetupLogger(cfg.Env)

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	errs.FailOnError(err, "failed connect to TG BOT")

	bot.Debug = true
	slogger.Debug("Authorized on account", slog.String("username", bot.Self.UserName))

	topic := "worker-handler"
	consumer := mqconsumer.New(topic, cfg.RabbitHost, slogger)

	msgs, err := consumer.Consume(topic + "_1")
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	for msg := range msgs {
		reqModel := &requestmodel.RequestCommand{}
		err := json.Unmarshal(msg.Body, reqModel)
		if err != nil {
			slog.Debug("Failed to unmarshal msg", slog.Any("data", string(msg.Body)))
			msg := tgbotapi.NewMessage(reqModel.ChatID, "Введены некорректные данные")
			bot.Send(msg)
			continue
		}
		msg := tgbotapi.NewMessage(reqModel.ChatID, "good")
		bot.Send(msg)
		slog.Info("MSG", slog.Any("data", reqModel))
	}

	slogger.Info("Waiting for messages...")
	<-forever
}
