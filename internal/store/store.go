package store

import (
	dbmodel "service-healthz-checker/internal/model/dbModel"
	requestmodel "service-healthz-checker/internal/model/requestModel"
)

type Store interface {
	Save(model requestmodel.Service) error
	AllUserServices(chatId int64) ([]*dbmodel.Service, error)
	ServiceInfoById(id int64) (*dbmodel.Service, error)
	RemoveService(id int64, chatId int64) error
}
