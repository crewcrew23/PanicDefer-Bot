CREATE TABLE IF NOT EXISTS history (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    chat_id INTEGER NOT NULL,
    status INTEGER,
    response_time_ms INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX IF NOT EXISTS idx_history_chat_id ON history(chat_id);
CREATE INDEX IF NOT EXISTS idx_history_created_at ON history(created_at);