CREATE TABLE IF NOT EXISTS request_logs (
    id UUID PRIMARY KEY,
    request_id UUID NOT NULL UNIQUE,
    model TEXT NOT NULL,
    provider TEXT,
    status_code INT NOT NULL,
    duration_ms INT NOT NULL,
    prompt_tokens INT NOT NULL DEFAULT 0,
    completion_tokens INT NOT NULL DEFAULT 0,
    total_tokens INT NOT NULL DEFAULT 0,
    input_cost_usd NUMERIC(12,8) NOT NULL DEFAULT 0,
    output_cost_usd NUMERIC(12,8) NOT NULL DEFAULT 0,
    total_cost_usd NUMERIC(12,8) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);