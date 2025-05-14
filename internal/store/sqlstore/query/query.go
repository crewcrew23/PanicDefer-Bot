package query

const (
	CREATE_SERVICE_Q = `
	INSERT INTO services (
		url, 
		chat_id, 
		last_status, 
		response_time_ms, 
		is_active,
		last_ping,
		updated_at,
		created_at
	)
	VALUES (
		:url, 
		:chat_id, 
		:last_status, 
		:response_time_ms, 
		:is_active,
		:last_ping,
		:updated_at,
		:created_at
	)`
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

	SELECT_DATA_FOR_PING = `
    SELECT id, url, chat_id, last_ping, last_status, 
           response_time_ms, is_active, created_at, updated_at 
    FROM services
    WHERE 
        is_active = TRUE 
        AND (
            last_ping IS NULL 
            OR last_ping < NOW() - INTERVAL '1 seconds'
        )
`

	//minute
	UPDATE_DATA = `
	UPDATE services
	SET 
		url = :url,
		chat_id = :chat_id,
		last_status = :last_status,
		response_time_ms = :response_time_ms,
		is_active = :is_active,
		last_ping = :last_ping,
		updated_at = :updated_at
	WHERE id = :id
	`
)
