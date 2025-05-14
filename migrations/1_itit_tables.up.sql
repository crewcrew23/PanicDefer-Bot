CREATE TABLE IF NOT EXISTS services (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    chat_id INTEGER NOT NULL,
    last_ping TIMESTAMP,
    last_status INTEGER,
    response_time_ms INTEGER,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_services_chat_id ON services(chat_id);
CREATE INDEX CONCURRENTLY idx_services_active_last_ping ON services(is_active, last_ping);