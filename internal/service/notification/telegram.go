package notification

import (
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TGNotifier struct {
	bot *tgbotapi.BotAPI
	log *slog.Logger
}

func NewTGNotifier(token string, log *slog.Logger) *TGNotifier {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Debug("failed connect to TG BOT", slog.String("ERR", err.Error()))
	}
	log.Debug("Authorized on account", slog.String("username", bot.Self.UserName))

	return &TGNotifier{bot: bot, log: log}
}

func (t *TGNotifier) Send(chat_id int64, service, txt string) {
	msg := tgbotapi.NewMessage(chat_id, prepareMessage(service, txt))
	t.bot.Send(msg)
}

func prepareMessage(service, txt string) string {
	msg := fmt.Sprintf("⚠️_______%s_______⚠️\nВнимание: %s", service, txt)
	return msg
}
