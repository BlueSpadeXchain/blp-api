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
    status TEXT
);

-- Create Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid VARCHAR(16) UNIQUE NOT NULL DEFAULT encode(gen_random_bytes(8), 'hex'), -- 16-character hex value
    wallet_address VARCHAR(255) UNIQUE NOT NULL,
    wallet_type VARCHAR(50) NOT NULL,
    balance NUMERIC(30, 6) DEFAULT 0,       
    perp_balance NUMERIC(30, 6) DEFAULT 0,     
    escrow_balance NUMERIC(30, 6) DEFAULT 0,    
    stake_balance NUMERIC(30, 6) DEFAULT 0,     
    frozen_balance NUMERIC(30, 6) DEFAULT 0, 
    total_balance NUMERIC(30, 6) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);

-- add total_balance trigger to calc
CREATE OR REPLACE FUNCTION update_total_balance()
RETURNS TRIGGER AS $$
BEGIN
    NEW.total_balance := NEW.balance + NEW.perp_balance + NEW.escrow_balance + NEW.stake_balance + NEW.frozen_balance;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_total_balance_trigger
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_total_balance();


-- Create Orders Table
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid VARCHAR(16) REFERENCES users(userid) ON DELETE CASCADE,
    order_type VARCHAR(10) NOT NULL CHECK (order_type IN ('long', 'short')),
    leverage NUMERIC(5, 2) NOT NULL,
    pair VARCHAR(64) NOT NULL,
    collateral NUMERIC(20, 2) NOT NULL CHECK (collateral > 0),
    status VARCHAR(20) NOT NULL DEFAULT 'unsigned' CHECK (status IN ('unsigned', 'pending', 'filled', 'canceled', 'closed', 'liquidated')),
    entry_price NUMERIC(20, 2),
    liq_price NUMERIC(20, 2),
    created_at TIMESTAMP DEFAULT NOW(),
    signed_at TIMESTAMP DEFAULT NULL,
    ended_at TIMESTAMP DEFAULT NULL
);

-- Create Withdrawals Table
CREATE TABLE withdrawals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userid VARCHAR(16) REFERENCES users(userid) ON DELETE CASCADE,
    amount NUMERIC(20, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'completed', 'failed')),
    tx_hash VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create Deposit Table
CREATE TABLE deposits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Primary key
    userid VARCHAR(16) REFERENCES users(userid) ON DELETE CASCADE, -- References `users.userid`
    wallet_address VARCHAR(255) NOT NULL, -- Ethereum-style wallet address
    wallet_type VARCHAR(50) NOT NULL, -- signature type
    chain_id TEXT DEFAULT '0' NOT NULL, -- Chain ID as a string
    block TEXT DEFAULT '0' NOT NULL, -- Block number as a string
    block_hash VARCHAR(64) NOT NULL, -- 32-byte block hash (Ethereum format)
    tx_hash VARCHAR(64) UNIQUE NOT NULL, -- Unique transaction hash
    sender VARCHAR(64) NOT NULL, -- Sender wallet address
    deposit_nonce TEXT DEFAULT '0' NOT NULL, -- Unique deposit nonce as a string
    asset VARCHAR(64) NOT NULL, -- Asset/Token contract address
    amount TEXT NOT NULL DEFAULT '0', -- Large amount as string
    value NUMERIC(78, 9) NOT NULL DEFAULT 0, -- Value as whole number with 9 decimals
    created_at TIMESTAMP DEFAULT NOW() -- Timestamp for deposit creation
);