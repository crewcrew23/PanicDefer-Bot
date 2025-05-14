package main

import (
	"log"
	"service-healthz-checker/internal/config"
	"service-healthz-checker/internal/logger"
	"service-healthz-checker/internal/service"
	workerpool "service-healthz-checker/internal/service/workerPool"
	"service-healthz-checker/internal/store/sqlstore"
	"time"

	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.MustLoad()
	slogger := logger.SetupLogger(cfg.Env)

	db, err := sqlx.Open("postgres", cfg.DbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	store := sqlstore.New(db, slogger)
	service := service.NewPingService(store, slogger)

	workerpool.RunPool(service, slogger, 5, time.Duration(cfg.TimeToPing)*time.Millisecond)
}
