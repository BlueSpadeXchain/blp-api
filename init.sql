-- Enable pgcrypto extension for generating random hex values
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create the debug_logs table for logging errors
CREATE TABLE IF NOT EXISTS debug_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP DEFAULT NOW(),
    log_level VARCHAR(50),
    error JSONB,
    message TEXT,
    context JSONB,
    status TEXT,
);

-- Create Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid VARCHAR(16) UNIQUE NOT NULL DEFAULT encode(gen_random_bytes(8), 'hex'), -- 16-character hex value
    email VARCHAR(255) UNIQUE,
    username VARCHAR(255),
    social_accounts JSONB DEFAULT '[]',
    wallet_addresses JSONB DEFAULT '[]', -- Updated to store an array of [chainid, wallet public key]
    balance BIGINT DEFAULT 0,           -- Total balance in nano-USD (scaled by 1e9)
    escrow_balance BIGINT DEFAULT 0,    -- Escrow balance in nano-USD
    stake_balance BIGINT DEFAULT 0,     -- Staked balance in nano-USD
    frozen_balance BIGINT DEFAULT 0,    -- Frozen balance in nano-USD
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create Orders Table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid INT REFERENCES users(userid) ON DELETE CASCADE,
    order_type VARCHAR(50) NOT NULL CHECK (order_type IN ('long', 'short')),
    leverage NUMERIC(5, 2) NOT NULL CHECK (leverage <= 50),
    pair VARCHAR(50) NOT NULL,
    amount NUMERIC(20, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'filled', 'canceled')),
    entry_price NUMERIC(20, 2),
    mark_price NUMERIC(20, 2),
    liq_price NUMERIC(20, 2),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create Withdrawals Table
CREATE TABLE withdrawals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid INT REFERENCES users(userid) ON DELETE CASCADE,
    amount NUMERIC(20, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    tx_hash VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);


-- Example insert to test table
INSERT INTO debug_logs (log_level, error, message, context)
VALUES 
    ('ERROR', '{"code": 500, "message": "Server exited", "details": "A panic occurred and the server has exited.", "origin": "api.main.handler"}', 
    'An example error message for testing.', 
    '{"param": "value"}');
