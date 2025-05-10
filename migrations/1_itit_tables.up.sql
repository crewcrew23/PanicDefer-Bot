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