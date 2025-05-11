package dbmodel

import "time"

type Service struct {
	Id             int64     `db:"id"`
	Url            string    `db:"url"`
	ChatID         int64     `db:"chat_id"`
	LastPing       time.Time `db:"last_ping"`
	LastStatus     int       `db:"last_status"`
	ResponseTimeMs int       `db:"response_time_ms"`
	IsActive       bool      `db:"is_active"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
