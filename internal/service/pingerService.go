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

func (p *PingService) SaveHistory(data []*dbmodel.Service) {
	p.store.SaveHistory(data)
}

func (p *PingService) AvgResTime(id int64) (float64, error) {
	res, err := p.store.AvgResTime(id)
	if err != nil {
		return -1, err
	}

	return res, nil
}
