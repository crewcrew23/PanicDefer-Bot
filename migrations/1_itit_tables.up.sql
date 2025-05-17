CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    chat_id INTEGER NOT NULL,
    last_ping TIMESTAMP,
    last_status INTEGER,
    response_time_ms INTEGER,
    is_active BOOLEAN DEFAULT TRUE,
    last_err_msg TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_services_chat_id ON services(chat_id);