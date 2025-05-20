package sqlstore

import (
	"database/sql"
	"errors"
	"log/slog"
	dbmodel "service-healthz-checker/internal/model/dbModel"
	requestmodel "service-healthz-checker/internal/model/requestModel"
	"service-healthz-checker/internal/store/dberrs"
	"service-healthz-checker/internal/store/sqlstore/query"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db         *sqlx.DB
	timeToPing int
	log        *slog.Logger
}

func New(db *sqlx.DB, timeToPing int, log *slog.Logger) *Store {
	return &Store{db: db, timeToPing: timeToPing, log: log}
}

func (s *Store) Save(model *requestmodel.Service) error {
	service := &dbmodel.Service{
		Url:            model.Url,
		ChatID:         model.ChatID,
		LastPing:       time.Time{},
		LastStatus:     0,
		ResponseTimeMs: 0,
		IsActive:       true,
		LastErrMsg:     time.Now().UTC().Add(-time.Hour),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	_, err := s.db.NamedExec(query.CREATE_SERVICE_Q, service)
	if err != nil {
		s.log.Debug("Error from Save Method:")
		errIs := dberrs.IsUniqueConstraintError(err)
		if errIs {
			s.log.Debug("Error", slog.String("Value", err.Error()))
			return dberrs.ErrUniqueConstraint
		}

		errIs = dberrs.IsCheckConstraintError(err)
		if errIs {
			s.log.Debug("Error", slog.String("Value", err.Error()))
			return dberrs.ErrInvalidData
		}

		errIs = dberrs.IsNullFieldError(err)
		if errIs {
			s.log.Debug("Error", slog.String("Value", err.Error()))
			return dberrs.ErrIsNullField
		}
		s.log.Debug("Error", slog.String("Value", err.Error()))
		return err
	}

	return nil
}

func (s *Store) AllUserServices(chatId int64) ([]*dbmodel.Service, error) {
	var services []*dbmodel.Service
	err := s.db.Select(&services, query.SELECT_ALL_BY_CHAT_ID, chatId)
	if err != nil {
		s.log.Debug("Error from AllUserServices Method:")
		s.log.Debug("Error", slog.String("Value", err.Error()))
		return nil, dberrs.ErrGetRows
	}

	if len(services) == 0 {
		s.log.Debug("WARN", slog.String("MSG", "no rows is db"))
		return nil, nil
	}

	return services, nil
}

func (s *Store) ServiceInfoById(id, chatId int64) (*dbmodel.Service, error) {
	var service []dbmodel.Service
	err := s.db.Select(&service, query.SELECT_SERVICE_ID, id, chatId)
	if err != nil {
		s.log.Debug("Error from ServiceInfoById Method:")
		s.log.Debug("Error", slog.String("Value", err.Error()))
		return nil, dberrs.ErrGetRows
	}

	if len(service) == 0 {
		s.log.Debug("WARN", slog.String("MSG", "no rows is db"))
		return nil, nil
	}

	return &service[0], nil
}

func (s *Store) RemoveService(id, chatId int64) error {
	_, err := s.db.Query(query.DELETE_SERVICE_BY_ID, id, chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Debug("WARN", slog.String("MSG", "no rows is db"))
			return dberrs.ErrNoRows
		}

		if dberrs.IsNotEnoughtArgumentError(err) {
			s.log.Debug("WARN", slog.String("MSG", "not enought arg"))
			return dberrs.ErrNotEnoughtArgument
		}

		s.log.Debug("DB ERR", slog.String("ERR", err.Error()))
		return dberrs.ErrDbOperation
	}

	return nil
}

func (s *Store) ChangeActiveSet(id, chatId int64) error {
	_, err := s.db.Query(query.UPDATE_SERVICE_STATE_BY_ID, id, chatId)
	if err != nil {
		if dberrs.IsNotEnoughtArgumentError(err) {
			s.log.Debug("WARN", slog.String("MSG", "not enought arg"))
			return dberrs.ErrNotEnoughtArgument
		}

		s.log.Debug("DB ERR", slog.String("ERR", err.Error()))
		return dberrs.ErrDbOperation
	}

	return nil
}

func (s *Store) DataForPing() ([]*dbmodel.Service, error) {
	var services []*dbmodel.Service
	err := s.db.Select(&services, query.SELECT_DATA_FOR_PING, s.timeToPing/1000) //devision by 1000 because s.timeToPing has time in ms
	if err != nil {
		s.log.Debug("Error from DataForPing Method:")
		s.log.Debug("Error", slog.String("Value", err.Error()))
		return nil, dberrs.ErrGetRows
	}

	if len(services) == 0 {
		s.log.Debug("WARN", slog.String("MSG", "no rows is db"))
		return nil, nil
	}

	s.log.Info("GET DATA OK")
	return services, nil
}

func (s *Store) UpdateData(data []*dbmodel.Service) {
	if len(data) == 0 {
		return
	}

	for _, v := range data {
		v.LastPing = time.Now()
		v.UpdatedAt = time.Now()

		_, err := s.db.NamedExec(query.UPDATE_DATA, v)
		if err != nil {
			s.log.Debug("FAILED TO UPDATE DATA", slog.String("ERROR", err.Error()))
		}
	}
}

func (s *Store) SaveHistory(data []*dbmodel.Service) {
	for _, v := range data {
		h := &dbmodel.History{
			Url:            v.Url,
			ChatId:         v.ChatID,
			Status:         v.LastStatus,
			ResponseTimeMs: v.ResponseTimeMs,
			CreatedAt:      time.Now().UTC(),
		}

		_, err := s.db.NamedExec(query.SAVE_HISTORY_DATA, h)
		if err != nil {
			s.log.Debug("FAILED TO SAVE HISTORY DATA", slog.String("ERROR", err.Error()))
		}
	}
}
