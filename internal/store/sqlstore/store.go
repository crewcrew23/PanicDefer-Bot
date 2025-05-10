package sqlstore

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db  *sqlx.DB
	log *slog.Logger
}

func New(db *sqlx.DB, log *slog.Logger) *Store {
	return &Store{db: db, log: log}
}
