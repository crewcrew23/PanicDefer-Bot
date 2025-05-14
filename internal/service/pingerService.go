package service

import (
	"log/slog"
	dbmodel "service-healthz-checker/internal/model/dbModel"
	"service-healthz-checker/internal/store"
)

type PingService struct {
	store store.Store
	log   *slog.Logger
}

func NewPingService(store store.Store, log *slog.Logger) *PingService {
	return &PingService{store: store, log: log}
}

func (p *PingService) DataForPing() []*dbmodel.Service {
	data, err := p.store.DataForPing()
	if err != nil {
		p.log.Debug("ERROR GET DATA FROM STORE", slog.String("ERROR", err.Error()))
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	return data
}

func (p *PingService) UpdateData(data []*dbmodel.Service) {
	p.store.UpdateData(data)
}
