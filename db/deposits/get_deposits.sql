CREATE OR REPLACE FUNCTION get_deposits_by_userid(user_id VARCHAR)
RETURNS TABLE(
    id UUID,
    userid VARCHAR,
    wallet_address VARCHAR,
    wallet_type VARCHAR,
    chain_id TEXT,
    block TEXT,
    block_hash VARCHAR(64),
    tx_hash VARCHAR(64),
    sender VARCHAR(64),
    deposit_nonce TEXT,
    asset VARCHAR(64),
    amount TEXT,
    value NUMERIC(78, 9),
    created_at TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM deposits WHERE deposits.userid = user_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_deposits_by_address(wallet_addr VARCHAR, wallet_t VARCHAR)
RETURNS TABLE(
    id UUID,
    userid VARCHAR,
    wallet_address VARCHAR,
    wallet_type VARCHAR,
    chain_id TEXT,
    block TEXT,
    block_hash VARCHAR(64),
    tx_hash VARCHAR(64),
    sender VARCHAR(64),
    deposit_nonce TEXT,
    asset VARCHAR(64),
    amount TEXT,
    value NUMERIC(78, 9),
    created_at TIMESTAMP
) AS $$
BEGIN
    -- First, check if the user exists
    RETURN QUERY 
    SELECT * 
    FROM deposits 
    WHERE deposits.wallet_address = wallet_addr AND deposits.wallet_type = wallet_t;
END;
$$ LANGUAGE plpgsql;