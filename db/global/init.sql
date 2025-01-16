CREATE TABLE global_state (
    key TEXT PRIMARY KEY,
    value NUMERIC(30, 6) NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);
