package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	mqconsumer "service-healthz-checker/internal/MQ/consumer"
	"service-healthz-checker/internal/command"
	"service-healthz-checker/internal/config"
	"service-healthz-checker/internal/errs"
	"service-healthz-checker/internal/lib/grapth"
	"service-healthz-checker/internal/logger"
	requestmodel "service-healthz-checker/internal/model/requestModel"
	"service-healthz-checker/internal/service"
	"service-healthz-checker/internal/store/sqlstore"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.MustLoad()
	slogger := logger.SetupLogger(cfg.Env)

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	errs.FailOnError(err, "failed connect to TG BOT")

	bot.Debug = true
	slogger.Debug("Authorized on account", slog.String("username", bot.Self.UserName))

	consumer := mqconsumer.New(cfg.MqConfig.Topics.FromServerTopic, cfg.MqConfig.Host, slogger)

	db, err := sqlx.Open("postgres", cfg.DbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store := sqlstore.New(db, cfg.TimeToPing, slogger)
	service := service.NewDomainService(store, slogger, bot)

	msgs, err := consumer.Consume(cfg.MqConfig.Topics.FromServerTopic)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	for msg := range msgs {
		reqModel := &requestmodel.RequestCommand{}
		err := json.Unmarshal(msg.Body, reqModel)
		if err != nil {
			slog.Debug("Failed to unmarshal msg", slog.Any("data", string(msg.Body)))
			sendMessage(reqModel.ChatID, "Введены некорректные данные", bot)
			continue
		}

		switch reqModel.Command {
		case command.ADD:
			serviceModel := &requestmodel.Service{ChatID: reqModel.ChatID, Url: reqModel.Value}
			err := service.Save(serviceModel)
			if err != nil {
				continue
			}

			sendMessage(reqModel.ChatID, "Данные сохранены", bot)
			continue
		case command.REMOVE:
			convINT, _ := strconv.Atoi(reqModel.Value)
			serviceModel := &requestmodel.Service{ChatID: reqModel.ChatID, Id: int64(convINT)}
			err := service.RemoveService(serviceModel)
			if err == nil {
				sendMessage(reqModel.ChatID, "Данные удаленны", bot)
				continue
			}
			continue
		case command.LIST:
			serviceModel := &requestmodel.Service{ChatID: reqModel.ChatID, Url: reqModel.Value}
			res, err := service.AllUserServices(serviceModel)
			if err != nil || res == nil {
				continue
			}
			text := "📡 *Ваши сервисы для мониторинга* 📡"

			for _, v := range res {
				text += fmt.Sprintf("\n• [%d]: %s", v.Id, v.Url)
				text += "\n──────────────"
			}
			sendMessage(reqModel.ChatID, text, bot)
			continue
		case command.GET:
			convINT, _ := strconv.Atoi(reqModel.Value)
			serviceModel := &requestmodel.Service{ChatID: reqModel.ChatID, Id: int64(convINT)}
			res, err := service.ServiceInfoById(serviceModel)
			if res == nil || err != nil {
				continue
			}
			var text string
			text += fmt.Sprintf("📡 *%s* 📡\n", res.Url)
			text += fmt.Sprintf("ID: %d\n", res.Id)
			text += fmt.Sprintf("URL: %s\n", res.Url)
			text += fmt.Sprintf("Last Ping: %v\n", res.LastPing)
			text += fmt.Sprintf("Last Status: %d\n", res.LastStatus)
			text += fmt.Sprintf("Response Time (ms): %d\n", res.ResponseTimeMs)
			text += fmt.Sprintf("IS Active: %v\n", res.IsActive)
			text += fmt.Sprintf("Created At: %v\n", res.CreatedAt)
			text += fmt.Sprintf("Updated At: %v\n", res.UpdatedAt)

			sendMessage(reqModel.ChatID, text, bot)
			continue
		case command.CH:
			convINT, _ := strconv.Atoi(reqModel.Value)
			serviceModel := &requestmodel.Service{ChatID: reqModel.ChatID, Id: int64(convINT)}
			err := service.ChangeActiveSet(serviceModel)
			if err != nil {
				continue
			}

			text := "Статус изменён"
			sendMessage(reqModel.ChatID, text, bot)
			continue

		case command.HISTORY:
			convINT, _ := strconv.Atoi(reqModel.Value)
			serviceModel := &requestmodel.Service{ChatID: reqModel.ChatID, Id: int64(convINT)}
			res, err := service.History(serviceModel)
			if err != nil {
				continue
			}

			if res == nil {
				sendMessage(serviceModel.ChatID, "Записей не найденно", bot)
				continue
			}

			imgBytes, err := grapth.CreateGrapth(res)
			if err != nil {
				sendMessage(serviceModel.ChatID, "не удалось создать график", bot)
				continue
			}

			sendPhoto(serviceModel.ChatID, imgBytes, bot)
			continue

		case command.HELP:
			sendMessageMarkDown(reqModel.ChatID, command.HelpTXT(), bot)
			continue
		case command.START:
			sendMessageMarkDown(reqModel.ChatID, command.HelpTXT(), bot)
			continue

		}

	}

	slogger.Info("Waiting for messages...")
	<-forever
}

func sendMessage(chatID int64, message string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, message)
	bot.Send(msg)
}

func sendMessageMarkDown(chatID int64, message string, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

func sendPhoto(chatID int64, photoBytes []byte, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{
		Name:  "response_graph",
		Bytes: photoBytes,
	})
	bot.Send(msg)
}
