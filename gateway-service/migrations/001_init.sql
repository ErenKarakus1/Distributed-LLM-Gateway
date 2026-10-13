CREATE TABLE IF NOT EXISTS request_logs (
    id UUID PRIMARY KEY,
    request_id UUID NOT NULL UNIQUE,
    model TEXT NOT NULL,
    provider TEXT,
    status_code INT NOT NULL,
    duration_ms INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);