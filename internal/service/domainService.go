package service

import (
	"errors"
	"fmt"
	"log/slog"
	dbmodel "service-healthz-checker/internal/model/dbModel"
	requestmodel "service-healthz-checker/internal/model/requestModel"
	"service-healthz-checker/internal/store"
	"service-healthz-checker/internal/store/dberrs"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DomainService struct {
	store store.Store
	log   *slog.Logger
	bot   *tgbotapi.BotAPI
}

func New(store store.Store, log *slog.Logger, bot *tgbotapi.BotAPI) *DomainService {
	return &DomainService{store: store, log: log, bot: bot}
}

func (d *DomainService) Save(model *requestmodel.Service) error {
	err := d.store.Save(model)
	if err != nil {
		if errors.Is(err, dberrs.ErrUniqueConstraint) {
			d.sendMessage(model.ChatID, "Сервис уже существует")
			return fmt.Errorf("%w", err)
		}

		if errors.Is(err, dberrs.ErrInvalidData) {
			d.sendMessage(model.ChatID, "Сервис уже существует")
			return fmt.Errorf("%w", err)
		}

		if errors.Is(err, dberrs.ErrIsNullField) {
			d.sendMessage(model.ChatID, "Переданно недостаточно параметров")
			return fmt.Errorf("%w", err)
		}

		d.sendMessage(model.ChatID, "Непредвиденная ошибка")
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (d *DomainService) AllUserServices(model *requestmodel.Service) ([]*dbmodel.Service, error) {
	res, err := d.store.AllUserServices(model.ChatID)
	if err != nil {
		if errors.Is(err, dberrs.ErrGetRows) {
			d.sendMessage(model.ChatID, "Произошла ошибка получения данных")
			return nil, fmt.Errorf("%w", err)
		}
	}

	if res == nil {
		d.sendMessage(model.ChatID, "Сервисов не найденно")
		return nil, nil
	}

	return res, nil
}

func (d *DomainService) ServiceInfoById(model *requestmodel.Service) (*dbmodel.Service, error) {
	res, err := d.store.ServiceInfoById(model.Id, model.ChatID)
	if err != nil {
		if errors.Is(err, dberrs.ErrGetRows) {
			d.sendMessage(model.ChatID, "Произошла ошибка получения данных")
			return nil, fmt.Errorf("%w", err)
		}
	}

	if res == nil {
		d.sendMessage(model.ChatID, "Сервис не найден")
		return nil, nil
	}

	return res, nil
}

func (d *DomainService) RemoveService(model *requestmodel.Service) error {
	err := d.store.RemoveService(model.Id, model.ChatID)
	if err != nil {
		if errors.Is(err, dberrs.ErrNoRows) {
			d.sendMessage(model.ChatID, "Сервис не найден")
			return err
		}

		if errors.Is(err, dberrs.ErrNotEnoughtArgument) {
			d.sendMessage(model.ChatID, "Переданно недостаточно параметров")
			return err
		}

		d.sendMessage(model.ChatID, "Непредвиденная ошибка")
		return err
	}

	return nil
}

func (d *DomainService) ChangeActiveSet(model *requestmodel.Service) error {
	err := d.store.ChangeActiveSet(model.Id, model.ChatID)
	if err != nil {

		if errors.Is(err, dberrs.ErrNotEnoughtArgument) {
			d.sendMessage(model.ChatID, "Переданно недостаточно параметров")
			return err
		}

		d.sendMessage(model.ChatID, "Непредвиденная ошибка")
		return err
	}

	return nil
}

func (d *DomainService) sendMessage(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	d.bot.Send(msg)
}
