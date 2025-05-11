package query

const (
	CREATE_SERVICE_Q = `
		INSERT INTO services (url, chat_id, last_ping, last_status, response_time_ms, is_active)
		VALUES (:url, :chat_id, :last_ping, :last_status, :response_time_ms, :is_active)
	`

	SELECT_ALL_BY_CHAT_ID = `
		SELECT id, url, chat_id, last_ping, last_status, response_time_ms, is_active, created_at, updated_at 
		FROM services
		WHERE chat_id = $1
	`

	SELECT_SERVICE_ID = `
		SELECT id, url, chat_id, last_ping, last_status, response_time_ms, is_active, created_at, updated_at 
		FROM services
		WHERE id = $1 AND chat_id = $2
	`

	UPDATE_SERVICE_STATE_BY_ID = `
		UPDATE services
		SET is_active = NOT is_active
		WHERE id = $1 AND chat_id = $2
	`

	DELETE_SERVICE_BY_ID = `
		DELETE FROM services
		WHERE id = $1 AND chat_id = $2
	`
)
