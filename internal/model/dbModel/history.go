package dbmodel

import "time"

type History struct {
	Id             int64     `db:"id"`
	Url            string    `db:"url"`
	ChatId         int64     `db:"chat_id"`
	Status         int       `db:"status"`
	ResponseTimeMs int       `db:"response_time_ms"`
	CreatedAt      time.Time `db:"created_at"`
}
