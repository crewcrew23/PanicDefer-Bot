package dbmodel

import "time"

type Service struct {
	Id             int64
	Url            string
	ChatID         int64
	LastPing       time.Time
	LastStatus     int
	ResponseTimeMs int
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
