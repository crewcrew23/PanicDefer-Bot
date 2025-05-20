package store

import (
	dbmodel "service-healthz-checker/internal/model/dbModel"
	requestmodel "service-healthz-checker/internal/model/requestModel"
)

type Store interface {
	//handler worker
	Save(model *requestmodel.Service) error
	AllUserServices(chatId int64) ([]*dbmodel.Service, error)
	ServiceInfoById(id, chatId int64) (*dbmodel.Service, error)
	RemoveService(id, chatId int64) error
	ChangeActiveSet(id, chatId int64) error

	//pinger worker
	DataForPing() ([]*dbmodel.Service, error)
	UpdateData([]*dbmodel.Service)
	SaveHistory([]*dbmodel.Service)
}
